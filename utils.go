package main

import (
	"log"
	"net/url"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func checkHash(password string, hashed_password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(password))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// Check for valid URL
func isValidURL(str string) bool {
	u, err := url.Parse(str)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false // Require scheme and host
	}
	// Optionally restrict schemes
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return true
}
