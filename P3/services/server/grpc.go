package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/ujjwal-shekhar/stripe-clone/services/server/handler/auth"
	pb "github.com/ujjwal-shekhar/stripe-clone/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/stripe-clone/services/common/utils"
	"github.com/ujjwal-shekhar/stripe-clone/services/server/handler/server"
)

func InformGateway(bankname string, address string, tlsCreds credentials.TransportCredentials) {
	// Get the connection to the gateway
	conn, err := grpc.NewClient(utils.PAYMENT_GATEWAY_URL, grpc.WithTransportCredentials(tlsCreds))
	if err != nil {
		log.Fatalf("Failed to connect to payment gateway: %v", err)
	}
	gatewayClient := pb.NewStripeServiceClient(conn)
	log.Printf("Connected to payment gateway for registration\n")
	
	// Inform the gateway about the bank
	_, err = gatewayClient.BankRegister(
		context.Background(), 
		&pb.BankRegistrationRequest{
			Bankname: bankname, 
			Address: address,
		},
	)

	if err != nil {
		log.Fatalf("Failed to register bank: %v", err)
	}
	log.Printf("Informed gateway about bank: %s", bankname)
}

func StartBankServer(bankname string, reqPath string, keyPath string) {
	// Listen on a random port
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Listening on port: %s", lis.Addr().String())

	// Load the TLS credentials
	tlsCreds, err := auth.LoadTLSCredentials(reqPath, keyPath)
	if err != nil {
		log.Fatalf("Failed to load TLS credentials: %v", err)
	}
	log.Printf("TLS credentials loaded successfully")

	// Register the bank with the gateway
	InformGateway(bankname, lis.Addr().String(), tlsCreds)

	// Create a new bank server
	grpcServer := grpc.NewServer(grpc.Creds(tlsCreds))
	bankServer, err := bank.NewBankServerTLS(bankname, tlsCreds)
	if err != nil {
		log.Fatalf("Failed to create bank server: %v", err)
	}
	pb.RegisterBankServiceServer(grpcServer, bankServer)
	log.Printf("Bank server is running on port %s", lis.Addr().String())


	// Serve on the random port
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	// Print the port we are listening on
	log.Printf("Bank server is serving on port %s", lis.Addr().String())
}
