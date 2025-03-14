package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

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
		io.WriteString(os.Stdout, "Enter your command: ")
		fmt.Scanln(&input)

		// Process the user input
		switch input {
		case "pay":
			// Read the next three lines
			var recipient, bankname string
			var amount int32
			io.WriteString(os.Stdout, "Enter the recipient: ")
			fmt.Scanln(&recipient)
			io.WriteString(os.Stdout, "Enter the bankname: ")
			fmt.Scanln(&bankname)
			io.WriteString(os.Stdout, "Enter the amount: ")
			fmt.Scanln(&amount)
			cli.MakePayment(amount, recipient, bankname)
		case "balance":
			cli.Balance()
		case "exit":
			return
		default:
			log.Printf("Invalid input")
		}
	}
}