package mysql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	cfg := Config{
		User:     "root",
		Password: "1111111",
		DBName:   "test-ss-db",
	}
	client, err := NewClient(&cfg)
	require.Nil(t, err)
	require.Nil(t, client.Ping())
}
