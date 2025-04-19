package conf

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"time"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/utils"
	v1 "github.com/fatedier/frp/pkg/config/v1"
)

func MasterDefaultSalt() string {
	ctx := context.Background()
	cfg := Get()
	if cfg.Master.CompatibleMode {
		logger.Logger(ctx).Warnf("master compatible mode enabled, use frp as default salt, which is not recommended")
		return "frp"
	}
	return utils.MD5(fmt.Sprintf("salt_%s:%d:%s",
		cfg.Master.InternalFRPServerHost,
		cfg.Master.InternalFRPServerPort,
		cfg.App.GlobalSecret))
}

func RPCListenAddr() string {
	cfg := Get()
	return fmt.Sprintf(":%d", cfg.Master.RPCPort)
}

func rpcCallAddr() string {
	cfg := Get()
	return fmt.Sprintf("%s:%d", cfg.Master.RPCHost, cfg.Master.RPCPort)
}

func InternalFRPServerToken() string {
	cfg := Get()
	return utils.MD5(fmt.Sprintf("%s:%d:%s",
		cfg.Master.InternalFRPServerHost,
		cfg.Master.InternalFRPServerPort,
		cfg.App.GlobalSecret))
}

func JWTSecret() string {
	cfg := Get()
	return utils.SHA1(fmt.Sprintf("%s:%d:%s", cfg.Master.APIHost, cfg.Master.APIPort, cfg.App.GlobalSecret))
}

func MasterAPIListenAddr() string {
	cfg := Get()
	return fmt.Sprintf(":%d", cfg.Master.APIPort)
}

func ServerAPIListenAddr() string {
	cfg := Get()
	return fmt.Sprintf(":%d", cfg.Server.APIPort)
}

func FRPsAuthOption(isDefault bool) v1.HTTPPluginOptions {
	cfg := Get()
	var port int
	if isDefault {
		port = cfg.Master.APIPort
	} else {
		port = cfg.Master.InternalFRPAuthServerPort
	}
	authUrl, err := url.Parse(fmt.Sprintf("http://%s:%d%s",
		cfg.Master.InternalFRPAuthServerHost,
		port,
		cfg.Master.InternalFRPAuthServerPath))
	if err != nil {
		logger.Logger(context.Background()).WithError(err).Fatalf("parse auth url error")
	}
	return v1.HTTPPluginOptions{
		Name: "multiuser",
		Ops:  []string{"Login"},
		Addr: authUrl.Host,
		Path: authUrl.Path,
	}
}

func GetCommonJWT(uid string) string {
	token, _ := utils.GetJwtTokenFromMap(JWTSecret(),
		time.Now().Unix(),
		int64(Get().App.CookieAge),
		map[string]string{common.UserIDKey: uid})
	return token
}

func GetCommonJWTWithExpireTime(uid string, expSec int) string {
	token, _ := utils.GetJwtTokenFromMap(JWTSecret(),
		time.Now().Unix(),
		int64(expSec),
		map[string]string{common.UserIDKey: uid})
	return token
}

func GetAPIURL() string {
	cfg := Get()

	if len(cfg.Client.APIUrl) != 0 {
		return cfg.Client.APIUrl
	}

	return fmt.Sprintf("%s://%s:%d", cfg.Master.APIScheme, cfg.Master.APIHost, cfg.Master.APIPort)
}

func GetCertTemplate() *x509.Certificate {
	cfg := Get()
	now := time.Now()
	return &x509.Certificate{
		SerialNumber: big.NewInt(now.Unix()),
		Subject: pkix.Name{
			CommonName:         cfg.Master.APIHost,
			Country:            []string{"CN"},
			Organization:       []string{"frp-panel"},
			OrganizationalUnit: []string{"frp-panel"},
		},
		SignatureAlgorithm:    x509.SHA512WithRSA,
		DNSNames:              []string{cfg.Master.APIHost},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
		NotBefore:             now,
		NotAfter:              now.AddDate(10, 0, 0),
		SubjectKeyId:          []byte{102, 114, 112, 45, 112, 97, 110, 101, 108},
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		KeyUsage: x509.KeyUsageKeyEncipherment |
			x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageKeyAgreement |
			x509.KeyUsageDataEncipherment,
	}
}

type LisOpt struct {
	MuxLis net.Listener
	ApiLis net.Listener
	RunAPI bool
}

func GetListener(c context.Context) LisOpt {
	runAPI := RPCListenAddr() != MasterAPIListenAddr()

	muxLis, err := net.Listen("tcp", RPCListenAddr())
	if err != nil {
		logger.Logger(c).WithError(err).Fatalf("failed to listen: %v", RPCListenAddr())
	}

	opt := LisOpt{
		MuxLis: muxLis,
		RunAPI: runAPI,
	}

	if runAPI {
		apiLis, err := net.Listen("tcp", MasterAPIListenAddr())
		if err != nil {
			logger.Logger(c).WithError(err).Warnf("failed to listen: %v, but mux server can handle http api", MasterAPIListenAddr())
		}
		opt.ApiLis = apiLis
	}

	if !runAPI {
		opt.ApiLis = nil
	}

	return opt
}

type ConnInfo struct {
	Host   string
	Scheme Scheme
}

type Scheme string

const (
	GRPC Scheme = "grpc"
	WS   Scheme = "ws"
	WSS  Scheme = "wss"
)

func GetRPCConnInfo() ConnInfo {
	cfg := Get()
	rpcUrl := cfg.Client.RPCUrl

	if len(rpcUrl) == 0 {
		return ConnInfo{
			Host:   rpcCallAddr(),
			Scheme: GRPC,
		}
	}

	parsedUrl, err := url.Parse(rpcUrl)
	if err != nil {
		logger.Logger(context.Background()).WithError(err).Fatalf("parse rpc url error")
	}

	return ConnInfo{
		Host:   parsedUrl.Host,
		Scheme: Scheme(parsedUrl.Scheme),
	}
}
