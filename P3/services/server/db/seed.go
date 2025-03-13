package db

import (
	"log"

	"github.com/ujjwal-shekhar/stripe-clone/services/server/handler/auth"
)

func SeedUsers() {
	users := []struct {
		Username string
		Password string
		Role     string
		Balance  int
	}{
		{"admin", "adminpass", "ADMIN", 10000},
		{"user1", "user1pass", "CUSTOMER", 5000},
		{"user2", "user2pass", "CUSTOMER", 3000},
	}

	for _, u := range users {
		hashedPassword, err := auth.EncryptPassword(u.Password)
		if err != nil {
			log.Fatalf("Error hashing password: %v", err)
		}

		_, err = DB.Exec(`
				INSERT INTO users (username, password_hash, role, balance)
				SELECT ?, ?, ?, ?
				WHERE NOT EXISTS (
					SELECT 1 FROM users WHERE username = ?
				)
			`, u.Username, hashedPassword, u.Role, u.Balance, u.Username)
		if err != nil {
			log.Printf("Skipping user %s: %v", u.Username, err)
		}
	}
	log.Println("Seed users inserted")
}
