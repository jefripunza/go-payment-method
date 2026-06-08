package main

import (
	"github.com/joho/godotenv"
)

// parses the .env file using the godotenv library from local or parent directory.
func init() {
	// Try loading from .env in current directory, fallback to parent directory if it fails
	if err := godotenv.Load(".env"); err != nil {
		_ = godotenv.Load("../.env")
	}
}
