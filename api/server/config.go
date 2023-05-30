package server

import (
	"errors"
	"time"

	"github.com/SyntSugar/ss-infra-go/consts"
	"github.com/SyntSugar/ss-infra-go/tracing"
)

const (
	defaultAPIAddr   = "127.0.0.1:8080"
	defaultAdminAddr = "127.0.0.1:9999"
)

type AdminCfg struct {
	Addr     string `mapstructure:"addr" json:"addr"`
	BasePath string `mapstructure:"basepath" json:"base_path"`
}

type APICfg struct {
	Addr     string `mapstructure:"addr" json:"addr"`
	BasePath string `mapstructure:"basepath" json:"base_path"`
}

type AccessLogCfg struct {
	Enabled bool   `mapstructure:"enabled"`
	Pattern string `mapstructure:"pattern"`
}

type Config struct {
	API           *APICfg            `mapstructure:"api"`
	Admin         *AdminCfg          `mapstructure:"admin"`
	OpenTelemetry *tracing.OTLConfig `mapstructure:"opentelemetry"`
	AccessLog     AccessLogCfg       `mapstructure:"access_log"`
	Shutdown      time.Duration      `mapstructure:"shutdown"`
}

func DefaultConfig() *Config {
	return &Config{
		API:   &APICfg{Addr: defaultAPIAddr},
		Admin: &AdminCfg{Addr: defaultAdminAddr},
		AccessLog: AccessLogCfg{
			Pattern: consts.DefaultAccessLogPattern,
		},
	}
}

func (cfg *Config) Validate() error {
	if cfg.API == nil && cfg.Admin == nil {
		return errors.New("api/admin config SHOULD NOT be empty at the same time")
	}
	return nil
}

func (cfg *Config) init() {
	if cfg.API != nil && cfg.Admin == nil {
		cfg.Admin = &AdminCfg{Addr: defaultAdminAddr}
	}
}
