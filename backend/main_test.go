package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	healthCheck(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", rr.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if resp["status"] != "ok" {
		t.Errorf("unexpected response: %v", resp)
	}
}

func TestSaveEnvironmentSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "env.log")
	os.Setenv("ENV_LOG_FILE", logPath)
	defer os.Unsetenv("ENV_LOG_FILE")

	body := `{"foo":"bar"}`
	req := httptest.NewRequest(http.MethodPost, "/environments", strings.NewReader(body))
	rr := httptest.NewRecorder()

	saveEnvironment(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201 got %d", rr.Code)
	}
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 log line got %d", len(lines))
	}
	var rec map[string]json.RawMessage
	if err := json.Unmarshal([]byte(lines[0]), &rec); err != nil {
		t.Fatalf("invalid log JSON: %v", err)
	}
	if !strings.Contains(string(rec["environment"]), "foo") {
		t.Errorf("log missing environment: %s", rec["environment"])
	}
}

func TestSaveEnvironmentEmpty(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/environments", http.NoBody)
	rr := httptest.NewRecorder()

	saveEnvironment(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 got %d", rr.Code)
	}
	b, _ := io.ReadAll(rr.Body)
	if len(b) == 0 {
		t.Errorf("expected error body")
	}
}
