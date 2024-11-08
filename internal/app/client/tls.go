package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

func applyTLS() (credentials.TransportCredentials, error) {
	caCert, err := os.ReadFile("tls/server.crt")
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}

	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("failed to append CA certificate to CA pool")
	}

	tlsConfig := &tls.Config{
		RootCAs:    caPool,
		MinVersion: tls.VersionTLS13,
	}

	creds := credentials.NewTLS(tlsConfig)

	return creds, nil
}
