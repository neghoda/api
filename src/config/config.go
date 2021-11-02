package config

import (
	"time"
)

type Config struct {
	LogLevel          string   `mapstructure:"LOG_LEVEL" default:"DEBUG"`
	HTTPConfig        HTTP     `mapstructure:"HTTP_SERVER"`
	PostgresCfg       Postgres `mapstructure:"POSTGRES_MASTER"`
	PostgresTestCfg   Postgres `mapstructure:"POSTGRES_TEST"`
	RefreshTokenLen   int      `mapstructure:"REFRESH_TOKEN_LEN" default:"32"`
	AccessTokenTTL    int      `mapstructure:"ACCESS_TOKEN_TTL_SEC" default:"900"` // seconds
	AccessTokenSecret string   `mapstructure:"ACCESS_TOKEN_SECRET" default:""`
	UserSessionTTL    int      `mapstructure:"USER_SESSION_TTL_SEC" default:"86400"` // seconds
}

type HTTP struct {
	Port              int    `mapstructure:"PORT"  default:"8080"`
	URLPrefix         string `mapstructure:"URL_PREFIX"  default:"/api"`
	SwaggerEnable     bool   `mapstructure:"SWAGGER_ENABLE"  default:"true"`
	SwaggerServeDir   string `mapstructure:"SWAGGER_SERVE_DIR"  default:"./src/server/http/static/"`
	CSRFSecuredCookie bool   `mapstructure:"CSFR_SECURED_COOKIE"  default:"true"`
	CORSAllowedHost   string `mapstructure:"CORS_ALLOWED_HOST"  default:"*"`
}

type Postgres struct {
	Host         string        `mapstructure:"HOST"          default:"localhost"`
	Port         string        `mapstructure:"PORT"          default:"5432"`
	Name         string        `mapstructure:"NAME"          default:"erp-db"`
	User         string        `mapstructure:"USER"          default:"postgres"`
	Password     string        `mapstructure:"PASSWORD"      default:"12345"`
	PoolSize     int           `mapstructure:"POOL_SIZE"     default:"10"`
	MaxRetries   int           `mapstructure:"MAX_RETRIES"   default:"5"`
	ReadTimeout  time.Duration `mapstructure:"READ_TIMEOUT"  default:"10s"`
	WriteTimeout time.Duration `mapstructure:"WRITE_TIMEOUT" default:"10s"`
	EnableLogger bool          `mapstructure:"ENABLE_LOGGER" default:"true"`
}
