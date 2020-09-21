package config

import (
	"github.com/caarlos0/env/v6"
)

// PgConnection contains variables for postgres db connection
type PgConnection struct {
	Username  string `env:"DB_USERNAME,required"`
	Password  string `env:"DB_PASSWORD,required"`
	Host      string `env:"DB_HOST" envDefault:"localhost"`
	Port      int    `env:"DB_PORT" envDefault:"5432"`
	DbName    string `env:"DB_NAME,required"`
	SslEnable bool   `env:"DB_SSL_ENABLE" envDefault:"false"`
}

// VkAuth contains variables for vk authorization
type VkAuth struct {
	AppID        string `env:"VK_APP_ID,required"`
	ClientSecret string `env:"VK_CLIENT_SECRET,required"`
	RedirectURI  string `env:"VK_REDIRECT_URI,required"`
}

// Config all app variables are stored here
type Config struct {
	Db        PgConnection
	Vk        VkAuth
	DebugMode bool   `env:"DEBUG_MODE" envDefault:"false"`
	JWTSecret string `env:"JWT_SECRET" envDefault:"secret"`
}

// New returns a new Config struct
func New() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
