package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

const tlsKey = "custom"

func NewClient(cfg *Config) (*sql.DB, error) {
	if cfg == nil {
		return nil, errors.New("cfg was nil")
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	cfg.init()

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%s&readTimeout=%s&writeTimeout=%s&parseTime=%t&charset=%s",
		cfg.User,
		cfg.Password,
		cfg.Addr,
		cfg.DBName,
		cfg.DialTimeout.String(),
		cfg.ReadTimeout.String(),
		cfg.WriteTimeout.String(),
		cfg.EnableParseTime,
		cfg.Charset,
	)
	if cfg.TLS != nil {
		tlsConfig, err := cfg.TLS.Build()
		if err != nil {
			return nil, fmt.Errorf("build tls config err: %w", err)
		}
		if err := mysql.RegisterTLSConfig(tlsKey, tlsConfig); err != nil {
			return nil, fmt.Errorf("register tls config err: %w", err)
		}
		dsn += "&tls=" + tlsKey
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(cfg.ConnMaxIldeTime)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	return db, nil
}
