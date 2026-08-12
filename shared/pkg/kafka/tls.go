package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
)

func buildTLSConfig(caCertPEM string) (*tls.Config, error) {
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM([]byte(caCertPEM)) {
		return nil, fmt.Errorf("failed to parse CA certificate")
	}
	return &tls.Config{RootCAs: caCertPool}, nil
}