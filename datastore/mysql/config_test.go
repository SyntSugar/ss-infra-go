package mysql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_validate(t *testing.T) {
	cfg := Config{}
	cfg.init()
	err := cfg.validate()
	require.EqualError(t, err, "user can't be empty")
	cfg.User = "test"
	err = cfg.validate()
	require.EqualError(t, err, "password can't be empty")
	cfg.Password = "test"
	err = cfg.validate()
	require.EqualError(t, err, "dbname can't be empty")
	cfg.DBName = "test"
	err = cfg.validate()
	require.Nil(t, err)
	require.Equal(t, 16, cfg.MaxIdleConns)
}
