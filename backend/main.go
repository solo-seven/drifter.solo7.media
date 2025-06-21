package main

import (
        "encoding/json"
        "log"
        "net/http"

        "github.com/gorilla/handlers"
        "github.com/gorilla/mux"
)

func createHandler() http.Handler {
        router := mux.NewRouter()
        router.HandleFunc("/health", healthCheck).Methods("GET")

        allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
        allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

        return handlers.CORS(allowedOrigins, allowedMethods)(router)
}

func main() {
        log.Println("Drifter backend starting on :8080")
        log.Fatal(http.ListenAndServe(":8080", createHandler()))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
