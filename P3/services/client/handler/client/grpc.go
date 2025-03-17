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
						   tlsCreds credentials.TransportCredentials,
						   tm *TransactionManager) (*Client, error) {
	// Get the connection woohoo
	conn, err := grpc.NewClient(utils.PAYMENT_GATEWAY_URL, grpc.WithTransportCredentials(tlsCreds))
	if err != nil {
		log.Fatalf("Failed to connect to payment gateway: %v", err)
	}
	stripeClient := pb.NewStripeServiceClient(conn)
	log.Printf("Connected to payment gateway without RBAC\n")

	// At this point the connection is AUTHENTICATED
	// We start to AUTHORIZE the connection by using many RPCs
	// We forward a login request, this can be made without a token 
	// and has no RBAC restrictions
	resp, err := stripeClient.ClientLogin(
		context.Background(), 
		&pb.ClientLoginRequest{
			Username: username, 
			Bankname: bankname,
			Password: password,
		},
	)
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
		return nil, err
	}
	log.Printf("Client %s @ %s logged in successfully\n", username, bankname)

	// Close the old connection
	conn.Close()

	// Now a new connection with token injection enabled
	conn, err = grpc.NewClient(
		utils.PAYMENT_GATEWAY_URL, 
		grpc.WithTransportCredentials(tlsCreds),
		grpc.WithChainUnaryInterceptor(
			auth.TokenUnaryInterceptor(resp.Token),
			tm.IdempotencyKeyUnaryInterceptor(),
			auth.LoggerUnaryInterceptor(),
		),
	)

	if err != nil {
		log.Fatalf("Failed to connect to payment gateway: %v", err)
	}
	stripeClient = pb.NewStripeServiceClient(conn)
	log.Printf("Connected to payment gateway with RBAC\n")
	
	return &Client{
		stripeClient: 	stripeClient,
		username: 		username,
		bankname: 		bankname,
		jwt_token: 		resp.Token,
	}, nil
}

func (c *Client) Balance() (int32, error) {
	// Only the JWT token matters since everything else gets 
	// overwritten by the claims in the token at gateway 
	// before reaching the service
	resp, err := c.stripeClient.CheckBalance(
		context.Background(), 
		&pb.CheckBalanceRequest{
			Username: c.username,
			Bankname: c.bankname,
		},
	)
	if err != nil {
		log.Printf("Failed to get balance: %v", err)
		return 0, err
	}
	log.Printf("Balance of %s: %d\n", c.username, resp.Balance)

	return resp.Balance, nil
}

func (c *Client) MakePayment(amount int32, recipient_username string, recipient_bankname string) error {
	resp, err := c.stripeClient.MakePayment(
		context.Background(), 
		&pb.MakePaymentRequest{
			SenderUsername: c.username,
			SenderBankname: c.bankname,
			ReceiverUsername: recipient_username,
			ReceiverBankname: recipient_bankname,
			Amount: amount,
		},
	)
	if err != nil || !resp.Success {
		log.Printf("Failed to make payment: %v", err)
		return err
	}
	log.Printf("Payment of %d to %s @ %s successful\n", 
				amount, recipient_username, recipient_bankname)

	return nil
}