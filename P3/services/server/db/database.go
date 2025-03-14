package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ujjwal-shekhar/stripe-clone/services/server/handler/auth"
)

var DB *sql.DB

func InitDB(dbPath string) {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		role INTEGER NOT NULL,
		balance INTEGER NOT NULL DEFAULT 0
	);`

	_, err = DB.Exec(createUserTable)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Println("Database initialized successfully")
}

func VerifyClientCredentials (username string, password string) (string, bool, error) {
	log.Printf("Verifying credentials for user: %s\n", username)

	var hashedPassword []byte
	var role string
	err := DB.QueryRow("SELECT password_hash, role FROM users WHERE username = ?", username).Scan(&hashedPassword, &role)
	if err != nil {
		return "", false, err
	}

	log.Printf("Role of user: %s is %s\n", username, role)

	return role, auth.ComparePasswords(hashedPassword, password), nil
}

func GetBalance(username string) (int, error) {
	log.Printf("Getting balance for user: %s\n", username)

	var balance int
	err := DB.QueryRow("SELECT balance FROM users WHERE username = ?", username).Scan(&balance)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func CapableOfDeducting(username string, amount int) (bool, error) {
	log.Printf("Checking if user: %s can deduct %d\n", username, amount)

	var balance int
	err := DB.QueryRow("SELECT balance FROM users WHERE username = ?", username).Scan(&balance)
	if err != nil {
		return false, err
	}

	return balance >= amount, nil
}