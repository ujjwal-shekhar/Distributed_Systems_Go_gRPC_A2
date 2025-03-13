package main

import (
	"log"

	"github.com/ujjwal-shekhar/stripe-clone/services/client/handler/auth"
	"github.com/ujjwal-shekhar/stripe-clone/services/client/handler/client"
)

func StartClient(username string, bankname string, password string, 
				 reqPath string, keyPath string) (*client.Client) {

	// Load the TLS credentials
	tlsCreds, err := auth.LoadTLSCredentials(reqPath, keyPath)
	if err != nil {
		log.Fatalf("Failed to load TLS credentials: %v", err)
	}

	// Create a new client
	cli, err := client.NewClientWithPassword(username, bankname, password, tlsCreds)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	log.Printf("Client %s from %s connected successfully\n", username, bankname)

	return cli
}