package main

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHealthEndpoint(t *testing.T) {
    handler := createHandler()
    server := httptest.NewServer(handler)
    defer server.Close()

    resp, err := http.Get(server.URL + "/health")
    if err != nil {
        t.Fatalf("failed to call health endpoint: %v", err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("expected status 200, got %d", resp.StatusCode)
    }
    var body map[string]string
    if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
        t.Fatalf("failed to decode response: %v", err)
    }
    if body["status"] != "ok" {
        t.Fatalf("unexpected body: %v", body)
    }
}
