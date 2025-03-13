package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/ujjwal-shekhar/stripe-clone/services/client/handler/client"
)

func main() {
	UNAME := flag.String("UNAME", "ujjwal", "The username of the client")
	BNAME := flag.String("BNAME", "Bank of America", "The name of the bank")
	PASS := flag.String("PASS", "password", "The password of the client")
	flag.Parse()

	cli, tm := StartClient(
		*UNAME, 
		*BNAME, 
		*PASS, 
		client.CLIENT_PREFIX + "cert/client-cert.pem",
		client.CLIENT_PREFIX + "cert/client-key.pem",
	)
	if cli == nil || tm == nil {
		log.Fatalf("Failed to start client")
	}

	// Keep taking user requests
	for {
		// Get the user input
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}

		// Parse the user input
		parts := strings.Split(input, " ")
		if len(parts) == 0 {
			log.Fatalf("Invalid input")
		}

		// Process the user input
		switch parts[0] {
		// case "pay":
		// 	if len(parts) != 3 {
		// 		log.Fatalf("Invalid input")
		// 	}
		// 	amount, err := strconv.ParseFloat(parts[1], 64)
		// 	if err != nil {
		// 		log.Fatalf("Invalid input")
		// 	}
		// 	cli.Pay(parts[2], amount)
		case "balance":
			cli.Balance()
		case "exit":
			return
		default:
			log.Fatalf("Invalid input")
		}
	}
}