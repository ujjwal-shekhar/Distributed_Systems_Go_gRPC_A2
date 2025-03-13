package client

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/ujjwal-shekhar/stripe-clone/services/client/handler/auth"
	pb "github.com/ujjwal-shekhar/stripe-clone/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/stripe-clone/services/common/utils"
)

type Client struct {
	stripeClient 		pb.StripeServiceClient
	username 			string
	bankname 			string
	jwt_token 			string
}

func NewClientWithPassword(username string, bankname string, password string,
						   tlsCreds credentials.TransportCredentials) (*Client, error) {
	// Get the connection woohoo
	conn, err := grpc.NewClient(utils.PAYMENT_GATEWAY_URL, grpc.WithTransportCredentials(tlsCreds))
	if err != nil {
		log.Fatalf("Failed to connect to payment gateway: %v", err)
	}
	stripeClient := pb.NewStripeServiceClient(conn)

	// At this point the connection is AUTHENTICATED
	// We start to AUTHORIZE the connection by using many RPCs

	resp, err := stripeClient.ClientLogin(context.Background(), 
										&pb.ClientLoginRequest{Username: username, 
															   Bankname: bankname,
															   Password: password,})
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
		return nil, err
	}
	
	if !resp.Success {
		log.Fatalf("Failed to login: Unsuccesful")
		return nil, err
	}

	// Now a new connection with token injection enabled
	conn, err = grpc.NewClient(utils.PAYMENT_GATEWAY_URL, grpc.WithTransportCredentials(tlsCreds),
								grpc.WithUnaryInterceptor(auth.TokenUnaryInterceptor(resp.Token)))
	if err != nil {
		log.Fatalf("Failed to connect to payment gateway: %v", err)
	}
	stripeClient = pb.NewStripeServiceClient(conn)
	
	return &Client{
		stripeClient: 	stripeClient,
		username: 		username,
		bankname: 		bankname,
		jwt_token: 		resp.Token,
	}, nil
}

// func (c *Client) Balance() (int64, error) {
// }