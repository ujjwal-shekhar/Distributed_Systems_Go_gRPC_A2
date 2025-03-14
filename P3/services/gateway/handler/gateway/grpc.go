package gateway

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/ujjwal-shekhar/stripe-clone/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/stripe-clone/services/gateway/handler/auth"
)

// BankClient represents a connection to a bank service.
type BankClient struct {
	Client  pb.BankServiceClient
	Address string
}

// Client represents a registered user in the payment gateway.
type Client struct {
	StripeClient pb.StripeServiceClient
	Username     string
	Passhash     string
	SessionToken string
	Role         pb.Role
}

// PaymentGatewayServer manages clients and banks.
type PaymentGatewayServer struct {
	pb.UnimplementedStripeServiceServer
	mu      sync.Mutex
	banks   map[string]BankClient
	clients map[string]Client
	Crm     *CachedResponseMap
}

// NewGatewayServerTLS initializes a new gateway server with TLS.
func NewGatewayServerTLS(tlsConfig credentials.TransportCredentials) *PaymentGatewayServer {
	return &PaymentGatewayServer{
		banks:   make(map[string]BankClient),
		clients: make(map[string]Client),
		Crm:     NewCachedResponseMap(),
	}
}

// UnaryMiddleware logs incoming unary RPC requests.
func UnaryMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Unary method: %s", info.FullMethod)
	return handler(ctx, req)
}

// ClientLogin authenticates a client and issues a JWT token.
func (s *PaymentGatewayServer) ClientLogin(ctx context.Context, req *pb.ClientLoginRequest) (*pb.ClientSessionResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Fetching the bank client to forward this request
	bankClient, exists := s.banks[req.Bankname]
	if !exists {
		return &pb.ClientSessionResponse{Success: false}, errors.New("bank not registered / offline")
	}
	log.Printf("ClientLogin: %s @ %s", req.Username, req.Bankname)

	// RPC to the bank service to verify client details
	clientSession, err := bankClient.Client.GetClientSession(ctx, req)
	if err != nil || !clientSession.Success {
		log.Printf("ClientLogin failed: WRONG CREDS : %s, Error: %v", req.Username, err)
		return &pb.ClientSessionResponse{Success: false}, err
	}
	log.Printf("ClientLogin successful: VERIFIED CREDS : %s, Role: %s", req.Username, clientSession.Role)

	token, err := auth.GenerateJWT(req.Username, req.Bankname, clientSession.Role)
	if err != nil {
		return &pb.ClientSessionResponse{Success: false}, err
	}
	log.Printf("ClientLogin successful: JWT GENERATED : %s", req.Username)

	return &pb.ClientSessionResponse{Success: true, Token: token, Role: clientSession.Role}, nil
}

// BankRegister registers a new bank with the payment gateway.
func (s *PaymentGatewayServer) BankRegister(ctx context.Context, req *pb.BankRegistrationRequest) (*pb.BankRegistrationResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf("Registering bank: %s", req.Bankname)

	tlsCreds, err := auth.LoadTLSCredentials(
		PAYMENT_GATEWAY_PREFIX+"cert/gateway-cert.pem",
		PAYMENT_GATEWAY_PREFIX+"cert/gateway-key.pem",
	)
	if err != nil {
		return &pb.BankRegistrationResponse{Success: false}, err
	}

	conn, err := grpc.NewClient(req.Address, grpc.WithTransportCredentials(tlsCreds))
	if err != nil {
		return &pb.BankRegistrationResponse{Success: false}, err
	}

	s.banks[req.Bankname] = BankClient{
		Client:  pb.NewBankServiceClient(conn),
		Address: req.Address,
	}

	log.Printf("Bank registered: %s", req.Bankname)
	return &pb.BankRegistrationResponse{Success: true}, nil
}

func (s *PaymentGatewayServer) CheckBalance(ctx context.Context, req *pb.CheckBalanceRequest) (*pb.CheckBalanceResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Forward the request to the bank
	bank, exists := s.banks[req.Bankname]
	if !exists {
		return &pb.CheckBalanceResponse{Balance: 0}, errors.New("bank not registered")
	}
	log.Printf("CheckBalance: %s @ %s", req.Username, req.Bankname)

	resp, err := bank.Client.CheckBalance(ctx, req)
	if err != nil {
		return &pb.CheckBalanceResponse{Balance: 0}, err
	}

	return resp, nil
}

func (s *PaymentGatewayServer) MakePayment(ctx context.Context, req *pb.MakePaymentRequest) (*pb.MakePaymentResponse, error) {
	s.mu.Lock()
	senderBank, senderBankExists := s.banks[req.SenderBankname]
	receiverBank, receiverBankExists := s.banks[req.ReceiverBankname]
	if !senderBankExists || !receiverBankExists {
		return &pb.MakePaymentResponse{Success: false}, errors.New("bank not registered")
	}
	s.mu.Unlock()

	// I will have to do a 2PC + Timeout here
	finalVote := true

	// Query phase will be in two go routines which we will wait on
	// unless they timeout, and then the vote will be taken to be "no"
	timeout_sender := make(chan bool, 1)
	vote_sender := make(chan bool, 1)
	
	// Start the query phase
	go func() {
		resp, err :=senderBank.Client.QueryPayment(
			ctx, &pb.QueryPaymentRequest{
				Username: req.SenderUsername,
				IsSender: true,
				Amount: req.Amount,
			},
		)
		vote_sender <- err == nil && resp.Vote
	}()
	go func() {
		time.Sleep(TIMEOUT_2PC)
		timeout_sender <- false
	}()
	log.Printf("MakePayment: Query phase started for sender: %s @ %s", req.SenderUsername, req.SenderBankname)

	// Select on the vote channels and the timeout channels
	select {
	case <-timeout_sender:
		log.Printf("MakePayment: Query phase timed out for sender: %s @ %s", req.SenderUsername, req.SenderBankname)
		finalVote = false
	case vote := <-vote_sender:
		log.Printf("MakePayment: Query phase successful for sender: %s @ %s", req.SenderUsername, req.SenderBankname)
		finalVote = finalVote && vote
	}
	
	// Receiver query phase
	timeout_receiver := make(chan bool, 1)
	vote_receiver := make(chan bool, 1)

	go func() {
		log.Printf("Query phase started for receiver: %s @ %s", req.ReceiverUsername, req.ReceiverBankname)
		resp, err := receiverBank.Client.QueryPayment(
			ctx, &pb.QueryPaymentRequest{
				Username: req.ReceiverUsername,
				IsSender: false,
				Amount: req.Amount,
			},
		)
		vote_receiver <- err == nil && resp.Vote
	}()
	go func() {
		time.Sleep(TIMEOUT_2PC)
		timeout_receiver <- false
	}()
	log.Printf("MakePayment: Query phase started for receiver: %s @ %s", req.ReceiverUsername, req.ReceiverBankname)

	select {
	case <-timeout_receiver:
		log.Printf("MakePayment: Query phase timed out for receiver: %s @ %s", req.ReceiverUsername, req.ReceiverBankname)
		finalVote = false
	case vote := <-vote_receiver:
		log.Printf("MakePayment: Query phase successful for receiver: %s @ %s", req.ReceiverUsername, req.ReceiverBankname)
		finalVote = finalVote && vote
	}

	log.Printf("MakePayment: Final vote: %v | Starting Commit/Rollback", finalVote)

	// Based on the votes gathered we will send the commit/rollback
	senderBank.Client.PersistPayment(
		ctx, &pb.PersistPaymentRequest{
			Username: req.SenderUsername,
			Amount: req.Amount,
			ToCommit: true,
			IsSender: true,
		},
	)
	receiverBank.Client.PersistPayment(
		ctx, &pb.PersistPaymentRequest{
			Username: req.ReceiverUsername,
			Amount: req.Amount,
			ToCommit: true,
			IsSender: false,
		},
	)

	if finalVote {
		log.Printf("MakePayment: Commit successful for both banks")
		return &pb.MakePaymentResponse{Success: true}, nil
	} else {
		log.Printf("MakePayment: Rollback successful for both banks")
		return &pb.MakePaymentResponse{Success: false}, nil
	}
}