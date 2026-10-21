package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/VaalaCat/frp-panel/biz/master/auth"
	masterserver "github.com/VaalaCat/frp-panel/biz/master/server"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/cache"
	"github.com/VaalaCat/frp-panel/services/dao"
	plugin "github.com/fatedier/frp/pkg/plugin/server"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type testDB struct {
	app.DBManager
	db *gorm.DB
}

func (d testDB) GetDefaultDB() *gorm.DB { return d.db }

func securityApp(t *testing.T) app.Application {
	t.Helper()
	config := conf.Config{}
	config.Master.CacheSize = 1
	cache.InitCache(config)
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "auth.db")), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	for _, sql := range []string{
		`CREATE TABLE users (user_id INTEGER PRIMARY KEY,user_name TEXT,email TEXT,password TEXT,token TEXT,deleted_at DATETIME)`,
		`CREATE TABLE clients (client_id TEXT PRIMARY KEY,connect_secret TEXT,deleted_at DATETIME)`,
		`CREATE TABLE servers (server_id TEXT PRIMARY KEY,connect_secret TEXT,deleted_at DATETIME)`,
		`CREATE TABLE certs (id INTEGER PRIMARY KEY,name TEXT,cert_file BLOB,key_file BLOB,deleted_at DATETIME)`,
	} {
		if err := db.Exec(sql).Error; err != nil {
			t.Fatal(err)
		}
	}
	hashBytes, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	hash := string(hashBytes)
	for _, stmt := range []struct {
		sql  string
		args []any
	}{
		{`INSERT INTO users(user_id,user_name,email,password,token) VALUES(1,?,?,?,?)`, []any{"security-existing", "security@example.test", hash, "valid-user-token"}},
		{`INSERT INTO clients(client_id,connect_secret) VALUES(?,?)`, []any{"test-client", "client-secret"}},
		{`INSERT INTO servers(server_id,connect_secret) VALUES(?,?)`, []any{"test-server", "server-secret"}},
		{`INSERT INTO certs(id,name,cert_file,key_file) VALUES(1,'default',?,?)`, []any{[]byte("public-certificate"), []byte("private-key-must-not-leak")}},
	} {
		if err := db.Exec(stmt.sql, stmt.args...).Error; err != nil {
			t.Fatal(err)
		}
	}
	a := app.NewApp()
	a.SetDBManager(testDB{db: db})
	return a
}

func TestCertificateRequiresKnownTypeAndCorrectSecret(t *testing.T) {
	a := securityApp(t)
	for _, tc := range []struct {
		kind       pb.ClientType
		id, secret string
		valid      bool
	}{
		{pb.ClientType_CLIENT_TYPE_UNSPECIFIED, "", "", false},
		{pb.ClientType(999), "test-client", "client-secret", false},
		{pb.ClientType(-1), "test-server", "server-secret", false},
		{pb.ClientType_CLIENT_TYPE_FRPC, "test-client", "wrong", false},
		{pb.ClientType_CLIENT_TYPE_FRPS, "missing", "wrong", false},
		{pb.ClientType_CLIENT_TYPE_FRPC, "test-client", "client-secret", true},
		{pb.ClientType_CLIENT_TYPE_FRPS, "test-server", "server-secret", true},
	} {
		res, err := auth.GetClientCert(app.NewContext(context.Background(), a), &pb.GetClientCertRequest{ClientType: tc.kind, ClientId: tc.id, ClientSecret: tc.secret})
		if tc.valid {
			if err != nil || string(res.GetCert()) != "public-certificate" {
				t.Fatalf("valid bootstrap %v failed: %v %v", tc.kind, res, err)
			}
		} else if err == nil || len(res.GetCert()) != 0 {
			t.Fatalf("unauthorized type %v returned certificate: %v %v", tc.kind, res, err)
		}
	}
}

func TestLoginFailuresAreIdenticalForJSONAndProtobuf(t *testing.T) {
	a := securityApp(t)
	for _, contentType := range []string{"application/json", "application/x-protobuf"} {
		router := gin.New()
		router.POST("/login", app.Wrapper(a, auth.LoginHandler))
		var previous *httptest.ResponseRecorder
		for _, username := range []string{"security-existing", "security-missing"} {
			password := "wrong-password"
			request := &pb.LoginRequest{Username: &username, Password: &password}
			var body []byte
			var err error
			if contentType == "application/json" {
				body, err = json.Marshal(request)
			} else {
				body, err = proto.Marshal(request)
			}
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", contentType)
			out := httptest.NewRecorder()
			router.ServeHTTP(out, req)
			if previous != nil && (previous.Code != out.Code || !bytes.Equal(previous.Body.Bytes(), out.Body.Bytes())) {
				t.Fatalf("%s enumerates users: %d %s / %d %s", contentType, previous.Code, previous.Body, out.Code, out.Body)
			}
			previous = out
		}
	}
	ok, user, err := dao.NewQuery(app.NewContext(context.Background(), a)).CheckUserPassword("security-existing", "correct-password")
	if err != nil || !ok || user.GetUserName() != "security-existing" {
		t.Fatalf("correct password rejected: %v %v", ok, err)
	}
}

func TestFRPPluginAndRPCUseUniformFailuresAndAcceptValidToken(t *testing.T) {
	a := securityApp(t)
	router := gin.New()
	router.POST("/auth", auth.MakeGinHandlerFunc(a, auth.HandleLogin))
	var rejected string
	for _, tc := range []struct {
		user, token string
		valid       bool
	}{
		{"security-existing", "wrong", false}, {"security-missing", "wrong", false}, {"", "", false}, {"security-existing", "valid-user-token", true},
	} {
		content := &plugin.LoginContent{}
		content.User = tc.user
		content.Metas = map[string]string{"token": tc.token}
		body, err := json.Marshal(&plugin.Request{Content: content})
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		out := httptest.NewRecorder()
		router.ServeHTTP(out, req)
		var response plugin.Response
		if err := json.Unmarshal(out.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		if tc.valid {
			if response.Reject || !response.Unchange {
				t.Fatalf("valid plugin credentials rejected: %s", out.Body)
			}
		} else {
			if !response.Reject {
				t.Fatal("invalid plugin accepted")
			}
			if rejected != "" && rejected != response.RejectReason {
				t.Fatal("plugin enumerates users")
			}
			rejected = response.RejectReason
		}
		responseRPC, err := masterserver.FRPAuth(app.NewContext(context.Background(), a), &pb.FRPAuthRequest{User: tc.user, Token: tc.token})
		if tc.valid {
			if err != nil || !responseRPC.GetOk() {
				t.Fatal("valid RPC credentials rejected")
			}
		} else if err == nil || responseRPC.GetOk() || responseRPC.GetStatus().GetMessage() != "invalid credentials" {
			t.Fatalf("invalid RPC differs: %v %v", responseRPC, err)
		}
	}
}
