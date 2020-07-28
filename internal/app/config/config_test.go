package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	os.Setenv("DB_USERNAME", "Trevor")
	os.Setenv("DB_PASSWORD", "Belmont")
	os.Setenv("DB_HOST", "castlevania.com")
	os.Setenv("DB_PORT", "7432")
	os.Setenv("DB_NAME", "draculaCastle")
	os.Setenv("DB_SSL_ENABLE", "true")

	c := New()
	require.Equal(t, "Trevor", c.Db.Username)
	require.Equal(t, "Belmont", c.Db.Password)
	require.Equal(t, "castlevania.com", c.Db.Host)
	require.Equal(t, 7432, c.Db.Port)
	require.Equal(t, "draculaCastle", c.Db.DbName)
	require.Equal(t, true, c.Db.SslEnable)
}

func TestGetEnv(t *testing.T) {
	os.Setenv("HOST", "db.com")
	require.Equal(t, "db.com", getEnv("HOST", "NOT_EXIST"))
}

func TestGetEnvDefault(t *testing.T) {
	os.Unsetenv("SOME_TEST_VAR")
	require.Equal(t, "NOT_EXIST", getEnv("SOME_TEST_VAR", "NOT_EXIST"))
}

func TestGetEnvInt(t *testing.T) {
	tests := []struct {
		testName   string
		varName    string
		envVal     string
		defaultVal int
		expected   int
	}{
		{"get from env", "PORT", "6432", 5432, 6432},
		{"take default", "PORT", "this_is_not_int", 5432, 5432},
		{"take default empty", "PORT", "", 5432, 5432},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			os.Setenv(tc.varName, tc.envVal)
			require.Equal(t, tc.expected, getEnvAsInt(tc.varName, tc.defaultVal))
		})
	}
}

func TestGetEnvBool(t *testing.T) {
	tests := []struct {
		testName   string
		varName    string
		envVal     string
		defaultVal bool
		expected   bool
	}{
		{"get from env", "DEBUG", "true", false, true},
		{"take default", "DEBUG", "maybe", false, false},
		{"take default empty", "DEBUG", "", false, false},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			os.Setenv(tc.varName, tc.envVal)
			require.Equal(t, tc.expected, getEnvAsBool(tc.varName, tc.defaultVal))
		})
	}
}

func TestGetEnvSlice(t *testing.T) {
	tests := []struct {
		testName   string
		varName    string
		envVal     string
		defaultVal []string
		expected   []string
		sep        string
	}{
        {"get from env", "SOME_LIST", "a,b,why", []string{"0", "zero"}, []string{"a", "b", "why"}, ","},
        {"get default", "SOME_LIST", "abwhy", []string{"0", "zero"}, []string{"abwhy"}, ","},
        {"get default empty", "SOME_LIST", "", []string{"0", "zero"}, []string{"0", "zero"}, ","},
        {"wrong sep", "SOME_LIST", "a,b,why", []string{"0", "zero"}, []string{"a,b,why"}, "."},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			os.Setenv(tc.varName, tc.envVal)
			require.Equal(t, tc.expected, getEnvAsSlice(tc.varName, tc.defaultVal, tc.sep))
		})
	}
}
