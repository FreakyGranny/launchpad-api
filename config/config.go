package config

import (
    "os"
    "strconv"
    "strings"
)

// PgConnection contains variables for postgres db connection
type PgConnection struct {
	username string
	password string
	host     string
	port     int
	dbName   string
}

// Config all app variables are stored here
type Config struct {
    db        PgConnection
    DebugMode bool
}

// New returns a new Config struct
func New() *Config {
    return &Config{
	db: PgConnection{
	    username: getEnv("DB_USERNAME", ""),
		password: getEnv("DB_PASSWORD", ""),
		host:     getEnv("DB_HOST", "localhost"),
		port:     getEnvAsInt("DB_PORT", 5432),
		dbName:   getEnv("DB_NAME", ""),
	},
	DebugMode: getEnvAsBool("DEBUG_MODE", true),
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
