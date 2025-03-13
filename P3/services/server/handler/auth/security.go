package auth

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"

	"github.com/ujjwal-shekhar/stripe-clone/services/common/utils"
)

func LoadTLSCredentials(reqPath string, keyPath string) (credentials.TransportCredentials, error) {
	// Root CA handling
	rootCA, err := os.ReadFile(utils.ROOT_CA_CERT_PATH)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(rootCA) {
		return nil, fmt.Errorf("failed to append client certs")
	}

	// Client certificate handling
	bankCert, err := tls.LoadX509KeyPair(reqPath, keyPath)
	if err != nil {
		return nil, err
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{bankCert},
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: certPool,
		RootCAs: certPool,
	}), nil
}