package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	cfg := &Config{}
	cfg.init()

	cfg.API = &APICfg{
		Addr: defaultAPIAddr,
	}
	cfg.init()
	assert.Equal(t, cfg.Admin.Addr, defaultAdminAddr)
}

func TestValidateConfig(t *testing.T) {
	cfg := &Config{}
	assert.NotNil(t, cfg.Validate())
	cfg.API = &APICfg{
		Addr: defaultAPIAddr,
	}
	assert.Nil(t, cfg.Validate())
}
