package conf

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials"
)

type Config struct {
	App struct {
		Secret         string `env:"SECRET" env-description:"app secret, for client and server frp salt"`
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
		CacheSize                 int    `env:"CACHE_SIZE" env-default:"100" env-description:"cache size in MB"`
		RPCHost                   string `env:"RPC_HOST" env-default:"127.0.0.1" env-description:"master host, is a public ip or domain"`
		RPCPort                   int    `env:"RPC_PORT" env-default:"9001" env-description:"master rpc port"`
		InternalFRPServerHost     string `env:"INTERNAL_FRP_SERVER_HOST" env-description:"internal frp server host, used for client connection"`
		InternalFRPServerPort     int    `env:"INTERNAL_FRP_SERVER_PORT" env-default:"9002" env-description:"internal frp server port, used for client connection"`
		InternalFRPAuthServerHost string `env:"INTERNAL_FRP_AUTH_SERVER_HOST" env-default:"127.0.0.1" env-description:"internal frp auth server host"`
		InternalFRPAuthServerPort int    `env:"INTERNAL_FRP_AUTH_SERVER_PORT" env-default:"8999" env-description:"internal frp auth server port"`
		InternalFRPAuthServerPath string `env:"INTERNAL_FRP_AUTH_SERVER_PATH" env-default:"/auth" env-description:"internal frp auth server path"`
	} `env-prefix:"MASTER_"`
	Server struct {
		APIPort int `env:"API_PORT" env-default:"8999" env-description:"server api port"`
	} `env-prefix:"SERVER_"`
	DB struct {
		Type string `env:"TYPE" env-default:"sqlite3" env-description:"db type, mysql or sqlite3 and so on"`
		DSN  string `env:"DSN" env-default:"data.db" env-description:"db dsn, for sqlite is path, other is dsn, look at https://github.com/go-sql-driver/mysql#dsn-data-source-name"`
	} `env-prefix:"DB_"`
	Client struct {
		ID     string `env:"ID" env-description:"client id"`
		Secret string `env:"SECRET" env-description:"client secret"`
	} `env-prefix:"CLIENT_"`
}

var (
	config     *Config
	ClientCred credentials.TransportCredentials
)

func Get() *Config {
	return config
}

func InitConfig() {
	var (
		err        error
		useEnvFile bool
	)

	envFiles := []string{
		".env",
		"/etc/frpp/.env",
	}

	for i, envFile := range envFiles {
		if err = godotenv.Load(envFile); err == nil {
			logrus.Infof("load env file success: %s", envFile)
			useEnvFile = true
			break
		}
		if i == len(envFiles)-1 {
			logrus.Errorf("cannot load env file: %s, error: %v", envFile, err)
			useEnvFile = false
			break
		}
		logrus.Infof("cannot load env file: %s, error: %v, try next", envFile, err)
	}

	if !useEnvFile {
		logrus.Info("use runtime env variables")
	}

	cfg := Config{}
	if err = cleanenv.ReadEnv(&cfg); err != nil {
		logrus.Panic(err)
	}
	cfg.Complete()

	config = &cfg
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
