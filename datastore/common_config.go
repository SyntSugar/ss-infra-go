package datastore

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

// ClientTLSConfig gives the user a easy way to configure the tls
// Usage example
//	clientConfig := &ClientTLSConfig{
//		CAFile:   "path/to/your/cafile",
//		CertFile: "path/to/your/certfile",
//		KeyFile:  "path/to/your/keyfile",
//	}
// Use local file system to store these files.

type ClientTLSConfig struct {
	CAFile   string `mapstructure:"ca_file"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`

	CAFileBytes   []byte
	CertFileBytes []byte
	KeyFileBytes  []byte

	InsecureSkipVerify bool `mapstructure:"inescure_skip_verify"`
}

func (config *ClientTLSConfig) readCertFiles(ctx context.Context) error {
	if config.CAFile != "" {
		bytes, err := os.ReadFile(config.CAFile)
		if err != nil {
			return fmt.Errorf("read ca file(%s) err: %w", config.CAFile, err)
		}
		config.CAFileBytes = bytes
	}
	if config.CertFile != "" {
		bytes, err := os.ReadFile(config.CertFile)
		if err != nil {
			return fmt.Errorf("read cert file(%s) err: %w", config.CertFile, err)
		}
		config.CertFileBytes = bytes
	}
	if config.KeyFile != "" {
		bytes, err := os.ReadFile(config.KeyFile)
		if err != nil {
			return fmt.Errorf("read key file(%s) err: %w", config.KeyFile, err)
		}
		config.KeyFileBytes = bytes
	}
	return nil
}

func (config *ClientTLSConfig) Build() (*tls.Config, error) {
	if err := config.readCertFiles(context.Background()); err != nil {
		return nil, err
	}

	tlsConfig := new(tls.Config)
	if config.CAFileBytes != nil {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(config.CAFileBytes)
		tlsConfig.RootCAs = caCertPool
	}
	if config.CertFileBytes != nil && config.KeyFileBytes != nil {
		cert, err := tls.X509KeyPair(config.CertFileBytes, config.KeyFileBytes)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	tlsConfig.InsecureSkipVerify = config.InsecureSkipVerify
	return tlsConfig, nil
}
