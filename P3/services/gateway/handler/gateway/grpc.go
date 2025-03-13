package gateway

import (
	"context"
	"errors"
	"log"
	"sync"

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