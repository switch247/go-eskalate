package config

import (
	"log"
	"os"
)

// Global variable to store the JWT secret
var JwtSecret = []byte("your_jwt_secret")

func init() {
	JwtSecretKey := os.Getenv("JWT_SECRETE_KEY")
	if JwtSecretKey != "" {
		JwtSecret = []byte(JwtSecretKey)
	} else {
		log.Fatal("JWT secret key not configured")
	}
}
