package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Setenv("POSTGRES_HOST", "localhost")
	t.Setenv("POSTGRES_USER", "testuser")
	t.Setenv("POSTGRES_PASSWORD", "testpassword")
	t.Setenv("POSTGRES_PORT", "5432")
	t.Setenv("POSTGRES_DB", "testdb")
	t.Setenv("POSTGRES_SLLMODE", "disable")
	t.Setenv("POSTGRES_TIMEZONE", "UTC")

	config := Load()

	assert.Equal(t, "localhost", config.Database.Localhost)
	assert.Equal(t, "testuser", config.Database.User)
	assert.Equal(t, "testpassword", config.Database.Password)
	assert.Equal(t, "5432", config.Database.Port)
	assert.Equal(t, "testdb", config.Database.DbName)
	assert.Equal(t, "disable", config.Database.SSLMode)
	assert.Equal(t, "UTC", config.Database.TimeZone)
}
