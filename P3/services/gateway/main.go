package main

import (
	"github.com/ujjwal-shekhar/stripe-clone/services/gateway/handler/gateway"
)

func main() {
	StartGatewayServer(
		gateway.PAYMENT_GATEWAY_PREFIX + "cert/gateway-cert.pem", 
		gateway.PAYMENT_GATEWAY_PREFIX + "cert/gateway-key.pem",
	)
}