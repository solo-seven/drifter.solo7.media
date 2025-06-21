package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Create a new router
	router := mux.NewRouter()

	// Register endpoints
	router.HandleFunc("/health", healthCheck).Methods("GET")
	router.HandleFunc("/environments", saveEnvironment).Methods("POST")

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

func saveEnvironment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	// Prepare log entry
	record := map[string]json.RawMessage{}
	ts := time.Now().UTC().Format(time.RFC3339)
	record["timestamp"] = json.RawMessage([]byte("\"" + ts + "\""))
	record["environment"] = body
	line, err := json.Marshal(record)
	if err != nil {
		http.Error(w, "failed to marshal record", http.StatusInternalServerError)
		return
	}

	logPath := os.Getenv("ENV_LOG_FILE")
	if logPath == "" {
		logPath = "logs/environments.log"
	}
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		http.Error(w, "failed to prepare log directory", http.StatusInternalServerError)
		return
	}

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		http.Error(w, "failed to open log", http.StatusInternalServerError)
		return
	}
	defer f.Close()
	if _, err := f.Write(append(line, '\n')); err != nil {
		http.Error(w, "failed to write log", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "saved"})
}
