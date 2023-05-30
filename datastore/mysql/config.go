package mysql

import (
	"errors"
	"time"

	"github.com/SyntSugar/ss-infra-go/datastore"
)

type Config struct {
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	Addr            string
	DBName          string
	User            string
	Password        string
	ConnMaxLifetime time.Duration
	ConnMaxIldeTime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
	EnableParseTime bool
	Charset         string

	TLS *datastore.ClientTLSConfig
}

func (cfg *Config) init() {
	cfg.Addr = "127.0.0.1:3306"
	cfg.DialTimeout = 3100 * time.Millisecond
	cfg.ReadTimeout = 5 * time.Second
	cfg.WriteTimeout = 5 * time.Second
	cfg.MaxOpenConns = 0
	cfg.MaxIdleConns = 16
	cfg.Charset = "utf8"
}

func (cfg *Config) validate() error {
	if cfg.User == "" {
		return errors.New("user can't be empty")
	}
	if cfg.Password == "" {
		return errors.New("password can't be empty")
	}
	if cfg.DBName == "" {
		return errors.New("dbname can't be empty")
	}
	return nil
}
