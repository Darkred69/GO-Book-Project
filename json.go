package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Function that specify response with error

func responseWithError(w http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Println("Responding with 5xx error: ", message)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	responseWithJSON(w, code, errResponse{Error: message})
}

// Function that specify the response with JSON
func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	log.Printf("Response: %v", string(response))
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Failed to marshal JSON. %v", err)
		return
	}
	w.Header().Add("Content-Type", "application/json") // Add header to JSON response
	w.WriteHeader(code)
	w.Write(response)
}
