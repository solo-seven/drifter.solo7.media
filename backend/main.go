package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Create a new router
	router := mux.NewRouter()

	// Register the health check endpoint
	router.HandleFunc("/health", healthCheck).Methods("GET")

	// Set up CORS
	callowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Start the HTTP server with CORS middleware
	log.Println("Drifter backend starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(callowedOrigins, allowedMethods)(router)))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
