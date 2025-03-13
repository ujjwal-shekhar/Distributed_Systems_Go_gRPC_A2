package bank

import (
	"fmt"
	"log"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/ujjwal-shekhar/stripe-clone/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/stripe-clone/services/common/utils"
	"github.com/ujjwal-shekhar/stripe-clone/services/server/db"
)

type BankServer struct {
	pb.UnimplementedBankServiceServer

	gatewayClient pb.StripeServiceClient
	Bankname 	  string
	Address 	  string
}

func NewBankServerTLS(bankname string, tlsCreds credentials.TransportCredentials) (*BankServer, error) {
	// First connect with the gateway as a client
	// and then we will register ourselves with the gateway
	conn, err := grpc.NewClient(utils.PAYMENT_GATEWAY_URL, grpc.WithTransportCredentials(tlsCreds))
	if err != nil {
		log.Fatalf("Failed to connect to payment gateway: %v", err)
		return nil, err
	}
	gatewayClient := pb.NewStripeServiceClient(conn)

	return &BankServer{
		gatewayClient: gatewayClient,
		Bankname: bankname,
		Address: "",
	}, nil
}

func (s *BankServer) GetClientSession(ctx context.Context, req *pb.ClientLoginRequest) (*pb.ClientSessionResponse, error) {
	// Check in the database if the client exists
	// If not, then return unsuccessful
	role, credsValid, err := db.VerifyClientCredentials(req.Username, req.Password)
	if err != nil {
		return &pb.ClientSessionResponse{Success: false, Token: ""}, err
	}
	
	if !credsValid {
		return &pb.ClientSessionResponse{Success: false, Token: ""}, fmt.Errorf("invalid credentials")
	}

	// Else, return it with success
	return &pb.ClientSessionResponse{Success: true, Token: "", Role: role}, nil
}