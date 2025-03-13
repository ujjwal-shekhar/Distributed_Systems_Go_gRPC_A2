package auth

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
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
		return nil, fmt.Errorf("failed to append root certs")
	}
	log.Printf("Root CA loaded successfully\n")

	// Client certificate handling
	clientCert, err := tls.LoadX509KeyPair(reqPath, keyPath)
	if err != nil {
		return nil, err
	}
	log.Printf("Client certificate loaded successfully\n")

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{clientCert},
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: certPool,
		RootCAs: certPool,
	}), nil
}