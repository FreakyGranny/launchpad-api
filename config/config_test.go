package config

import (
    "os"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestNewConfig(t *testing.T) {
    t.Parallel()
    
    os.Setenv("DB_USERNAME", "Trevor")
    os.Setenv("DB_PASSWORD", "Belmont")
    os.Setenv("DB_HOST", "castlevania.com")
    os.Setenv("DB_PORT", "7432")
    os.Setenv("DB_NAME", "draculaCastle")

    c := New()
    assert.Equal(t, "Trevor", c.db.username)
    assert.Equal(t, "Belmont", c.db.password)
    assert.Equal(t, "castlevania.com", c.db.host)
    assert.Equal(t, 7432, c.db.port)
    assert.Equal(t, "draculaCastle", c.db.dbName)
}

func TestGetEnv(t *testing.T) {
	// t.Skip()
    t.Parallel()
    os.Setenv("HOST", "db.com")
	assert.Equal(t, "db.com", getEnv("HOST", "NOT_EXIST"))
}

func TestGetEnvDefault(t *testing.T) {
    t.Parallel()
    os.Unsetenv("SOME_TEST_VAR")
	assert.Equal(t, "NOT_EXIST", getEnv("SOME_TEST_VAR", "NOT_EXIST"))
}

func TestGetEnvInt(t *testing.T) {
    t.Parallel()
    os.Setenv("PORT", "6432")
	assert.Equal(t, 6432, getEnvAsInt("PORT", 5432))
}

func TestGetEnvIntWrong(t *testing.T) {
    t.Parallel()
    os.Setenv("PORT", "this_is_not_int")
	assert.Equal(t, 5432, getEnvAsInt("PORT", 5432))
}

func TestGetEnvBool(t *testing.T) {
    t.Parallel()
    os.Setenv("DEBUG", "true")
	assert.True(t, getEnvAsBool("DEBUG", false))
}

func TestGetEnvBoolWrong(t *testing.T) {
    t.Parallel()
    os.Setenv("DEBUG", "maybe")
	assert.False(t, getEnvAsBool("DEBUG", false))
}

func TestGetEnvSlice(t *testing.T) {
    t.Parallel()
    os.Setenv("SOME_LIST", "a,b,why")
	assert.Equal(t, []string{"a", "b", "why"},getEnvAsSlice("SOME_LIST", []string{"0","zero"}, ","))
}

func TestGetEnvSliceDefault(t *testing.T) {
    t.Parallel()
    os.Unsetenv("SOME_LIST")
	assert.Equal(t, []string{"0", "zero"},getEnvAsSlice("SOME_LIST", []string{"0","zero"}, ","))
}
