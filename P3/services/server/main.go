package main

import (
	"flag"
	"fmt"

	"github.com/ujjwal-shekhar/stripe-clone/services/server/db"
	"github.com/ujjwal-shekhar/stripe-clone/services/server/handler/server"
)


func main() {
	// Fetch the Bank name from the flags and start the server
	BankName := flag.String("NAME", "Bank of America", "The name of the bank")
	flag.Parse()
	
	// Initialize the db
	db.InitDB(bank.BANK_DB_PREFIX + fmt.Sprintf("%s-db.db", *BankName))
	db.SeedUsers()

	// Start the bank server
	StartBankServer(*BankName, bank.BANK_PREFIX + "cert/bank-cert.pem", 
							   bank.BANK_PREFIX + "cert/bank-key.pem")

	// Send periodic heartbeats to the gateway
	// go SendHeartbeats(bankServer)
}