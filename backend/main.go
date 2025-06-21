package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// corsMiddleware adds CORS headers to responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				// Get the method the client wants to use
				requestMethod := r.Header.Get("Access-Control-Request-Method")
				if requestMethod != "" {
					w.Header().Set("Access-Control-Allow-Methods", strings.ToUpper(requestMethod))
				}

				// Get the headers the client wants to use
				requestHeaders := r.Header.Get("Access-Control-Request-Headers")
				if requestHeaders != "" {
					w.Header().Set("Access-Control-Allow-Headers", requestHeaders)
				}

				// Cache preflight response for 24 hours
				w.Header().Set("Access-Control-Max-Age", "86400")


				// End the request for preflight
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func createHandler() http.Handler {
	router := mux.NewRouter()

	// Register routes with explicit methods
	healthHandler := http.HandlerFunc(healthCheck)
	environmentHandler := http.HandlerFunc(saveEnvironment)

	// Apply CORS middleware to each route
	router.Handle("/health", healthHandler).Methods("GET")
	router.Handle("/environments", environmentHandler).Methods("POST")

	// Apply CORS middleware to the router
	handler := corsMiddleware(router)

	return handler
}

// setupServer creates and configures an HTTP server with the provided handler
func setupServer(handler http.Handler, port string) *http.Server {
	if port == "" {
		port = "8080"
	}

	return &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}
}

// runServer starts the HTTP server and listens for requests
func runServer(server *http.Server) error {
	log.Printf("Server starting on %s\n", server.Addr)
	return server.ListenAndServe()
}

func main() {
	port := os.Getenv("PORT")
	server := setupServer(createHandler(), port)
	log.Fatal(runServer(server))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func saveEnvironment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Check Content-Type header
	contentType := r.Header.Get("Content-Type")
	if contentType != "" && contentType != "application/json" {
		http.Error(w, `{"error":"Content-Type must be application/json"}`, http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"failed to read request body"}`, http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		http.Error(w, `{"error":"empty request body"}`, http.StatusBadRequest)
		return
	}

	// Validate JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		http.Error(w, `{"error":"invalid JSON: `+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Prepare log entry
	record := map[string]interface{}{
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
		"environment": jsonData,
	}

	line, err := json.Marshal(record)
	if err != nil {
		http.Error(w, `{"error":"failed to marshal record"}`, http.StatusInternalServerError)
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
