package config

import (
	"os"
	"strconv"
	"strings"
)

// PgConnection contains variables for postgres db connection
type PgConnection struct {
	Username  string
	Password  string
	Host      string
	Port      int
	DbName    string
	SslEnable bool
}

// VkAuth contains variables for vk authorization
type VkAuth struct {
	AppID        string
	ClientSecret string
	RedirectURI  string
}

// Config all app variables are stored here
type Config struct {
	Db        PgConnection
	Vk        VkAuth
	DebugMode bool
	JWTSecret string
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		Db: PgConnection{
			Username:  getEnv("DB_USERNAME", ""),
			Password:  getEnv("DB_PASSWORD", ""),
			Host:      getEnv("DB_HOST", "localhost"),
			Port:      getEnvAsInt("DB_PORT", 5432),
			DbName:    getEnv("DB_NAME", ""),
			SslEnable: getEnvAsBool("DB_SSL_ENABLE", false),
		},
		Vk: VkAuth{
			AppID:        getEnv("VK_APP_ID", ""),
			ClientSecret: getEnv("VK_CLIENT_SECRET", ""),
			RedirectURI:  getEnv("VK_REDIRECT_URI", ""),
		},
		DebugMode: getEnvAsBool("DEBUG_MODE", false),
		JWTSecret: getEnv("JWT_SECRET", "secret"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	return strings.Split(valStr, sep)
}
