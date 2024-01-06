package conf

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	App struct {
		GlobalSecret   string `env:"GLOBAL_SECRET" env-required:"true" env-description:"global secret, used in manager gen secret, keep it safe"`
		CookieAge      int    `env:"COOKIE_AGE" env-default:"86400" env-description:"cookie age in second, default is 1 day"`
		CookieName     string `env:"COOKIE_NAME" env-default:"frp-panel-cookie" env-description:"cookie name"`
		CookiePath     string `env:"COOKIE_PATH" env-default:"/" env-description:"cookie path"`
		CookieDomain   string `env:"COOKIE_DOMAIN" env-default:"" env-description:"cookie domain"`
		CookieSecure   bool   `env:"COOKIE_SECURE" env-default:"false" env-description:"cookie secure"`
		CookieHTTPOnly bool   `env:"COOKIE_HTTP_ONLY" env-default:"true" env-description:"cookie http only"`
	} `env-prefix:"APP_"`
	Master struct {
		APIPort                   int    `env:"API_PORT" env-default:"9000" env-description:"master api port"`
		APIHost                   string `env:"API_HOST" env-description:"master host, can behind proxy like cdn"`
		CacheSize                 int    `env:"CACHE_SIZE" env-default:"100" env-description:"cache size in MB"`
		RPCHost                   string `env:"RPC_HOST" env-required:"true" env-description:"master host, is a public ip or domain"`
		RPCPort                   int    `env:"RPC_PORT" env-default:"9001" env-description:"master rpc port"`
		InternalFRPServerHost     string `env:"INTERNAL_FRP_SERVER_HOST" env-description:"internal frp server host, used for client connection"`
		InternalFRPServerPort     int    `env:"INTERNAL_FRP_SERVER_PORT" env-default:"9002" env-description:"internal frp server port, used for client connection"`
		InternalFRPAuthServerHost string `env:"INTERNAL_FRP_AUTH_SERVER_HOST" env-default:"127.0.0.1" env-description:"internal frp auth server host"`
		InternalFRPAuthServerPort int    `env:"INTERNAL_FRP_AUTH_SERVER_PORT" env-default:"9000" env-description:"internal frp auth server port"`
		InternalFRPAuthServerPath string `env:"INTERNAL_FRP_AUTH_SERVER_PATH" env-default:"/auth" env-description:"internal frp auth server path"`
	} `env-prefix:"MASTER_"`
	Server struct {
		APIPort int `env:"API_PORT" env-default:"8999" env-description:"server api port"`
	} `env-prefix:"SERVER_"`
	DB struct {
		Type string `env:"TYPE" env-default:"sqlite3" env-description:"db type, mysql or sqlite3 and so on"`
		DSN  string `env:"DSN" env-default:"data.db" env-description:"db dsn, for sqlite is path, other is dsn, look at https://github.com/go-sql-driver/mysql#dsn-data-source-name"`
	} `env-prefix:"DB_"`
}

var (
	config *Config
)

func Get() *Config {
	return config
}

func InitConfig() {
	var (
		err error
	)

	if err = godotenv.Load(); err != nil {
		logrus.Infof("Error loading .env file, will use runtime env")
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
}
