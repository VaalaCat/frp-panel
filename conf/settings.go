package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/tidwall/pretty"
)

type Config struct {
	App struct {
		GlobalSecret   string `env:"GLOBAL_SECRET" env-default:"frp-panel" env-description:"global secret, used in manager gen secret, keep it safe"`
		CookieAge      int    `env:"COOKIE_AGE" env-default:"86400" env-description:"cookie age in second, default is 1 day"`
		CookieName     string `env:"COOKIE_NAME" env-default:"frp-panel-cookie" env-description:"cookie name"`
		CookiePath     string `env:"COOKIE_PATH" env-default:"/" env-description:"cookie path"`
		CookieDomain   string `env:"COOKIE_DOMAIN" env-default:"" env-description:"cookie domain"`
		CookieSecure   bool   `env:"COOKIE_SECURE" env-default:"false" env-description:"cookie secure"`
		CookieHTTPOnly bool   `env:"COOKIE_HTTP_ONLY" env-default:"true" env-description:"cookie http only"`
		EnableRegister bool   `env:"ENABLE_REGISTER" env-default:"false" env-description:"enable register, only allow the first admin to register"`
	} `env-prefix:"APP_"`
	Master struct {
		APIPort                   int    `env:"API_PORT" env-default:"9000" env-description:"master api port"`
		APIHost                   string `env:"API_HOST" env-description:"master host, can behind proxy like cdn"`
		APIScheme                 string `env:"API_SCHEME" env-default:"http" env-description:"master api scheme"`
		CacheSize                 int    `env:"CACHE_SIZE" env-default:"10" env-description:"cache size in MB"`
		RPCHost                   string `env:"RPC_HOST" env-default:"127.0.0.1" env-description:"master host, is a public ip or domain"`
		RPCPort                   int    `env:"RPC_PORT" env-default:"9001" env-description:"master rpc port"`
		InternalFRPServerHost     string `env:"INTERNAL_FRP_SERVER_HOST" env-description:"internal frp server host, used for client connection"`
		InternalFRPAuthServerHost string `env:"INTERNAL_FRP_AUTH_SERVER_HOST" env-default:"127.0.0.1" env-description:"internal frp auth server host"`
		InternalFRPAuthServerPort int    `env:"INTERNAL_FRP_AUTH_SERVER_PORT" env-default:"8999" env-description:"internal frp auth server port"`
		InternalFRPAuthServerPath string `env:"INTERNAL_FRP_AUTH_SERVER_PATH" env-default:"/auth" env-description:"internal frp auth server path"`
	} `env-prefix:"MASTER_"`
	Server struct {
		APIPort int `env:"API_PORT" env-default:"8999" env-description:"server api port"`
	} `env-prefix:"SERVER_"`
	DB struct {
		Type string `env:"TYPE" env-default:"sqlite3" env-description:"db type, mysql or sqlite3 and so on"`
		DSN  string `env:"DSN" env-default:"/data/data.db" env-description:"db dsn, for sqlite is path, other is dsn, look at https://github.com/go-sql-driver/mysql#dsn-data-source-name"`
	} `env-prefix:"DB_"`
	Client struct {
		ID                    string `env:"ID" env-description:"client id"`
		Secret                string `env:"SECRET" env-description:"client secret"`
		TLSRpc                bool   `env:"TLS_RPC" env-default:"true" env-description:"use tls for rpc connection"`
		RPCUrl                string `env:"RPC_URL" env-description:"rpc url, support ws or wss or grpc scheme, eg: ws://127.0.0.1:9000"`
		APIUrl                string `env:"API_URL" env-description:"api url, support http or https scheme, eg: http://127.0.0.1:9000"`
		TLSInsecureSkipVerify bool   `env:"TLS_INSECURE_SKIP_VERIFY" env-default:"true" env-description:"skip tls verify"`
	} `env-prefix:"CLIENT_"`
	IsDebug bool `env:"IS_DEBUG" env-default:"false" env-description:"is debug mode"`
	Logger  struct {
		DefaultLoggerLevel string `env:"DEFAULT_LOGGER_LEVEL" env-default:"info" env-description:"frp-panel internal default logger level"`
		FRPLoggerLevel     string `env:"FRP_LOGGER_LEVEL" env-default:"info" env-description:"frp logger level"`
	} `env-prefix:"LOGGER_"`
}

func NewConfig() Config {
	var (
		err        error
		useEnvFile bool
		ctx        = context.Background()
	)

	// 越前面优先级越高，后面的不会覆盖前面的
	envFiles := []string{
		defs.CurEnvPath,
		defs.SysEnvPath,
	}

	for _, envFile := range envFiles {
		if err = godotenv.Load(envFile); err == nil {
			logger.Logger(ctx).Infof("load env file success: %s", envFile)
			useEnvFile = true
		}
	}

	if !useEnvFile {
		logger.Logger(ctx).Info("use runtime env variables")
	}

	cfg := Config{}
	if err = cleanenv.ReadEnv(&cfg); err != nil {
		logger.Logger(ctx).Panic(err)
	}
	cfg.Complete()

	if !cfg.IsDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	return cfg
}

func (cfg *Config) Complete() {
	if len(cfg.Master.InternalFRPServerHost) == 0 {
		cfg.Master.InternalFRPServerHost = cfg.Master.RPCHost
	}

	if len(cfg.Master.APIHost) == 0 {
		cfg.Master.APIHost = cfg.Master.RPCHost
	}

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(cfg.Client.ID) == 0 {
		cfg.Client.ID = hostname
	}
}

func (cfg Config) PrintStr() string {
	raw, _ := json.Marshal(cfg)
	return string(pretty.Pretty(raw))
}
