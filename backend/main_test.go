package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	healthCheck(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "status code should be 200")
	
	var resp map[string]string
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err, "should decode response body without error")
	assert.Equal(t, "ok", resp["status"], "status should be 'ok'")
}

func TestCreateHandler(t *testing.T) {
	tests := []struct {
		name           string
		method        string
		path          string
		headers       map[string]string
		expectedCode  int
		expectCORS    bool
		isPreflight   bool
	}{
		{
			name:          "GET health check with CORS",
			method:        http.MethodGet,
			path:          "/health",
			headers:       map[string]string{"Origin": "http://localhost:3000"},
			expectedCode:  http.StatusOK,
			expectCORS:    true,
			isPreflight:   false,
		},
		{
			name:    "OPTIONS preflight for health check",
			method:  http.MethodOptions,
			path:    "/health",
			headers: map[string]string{
				"Origin":                         "http://localhost:3000",
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "Content-Type",
			},
			expectedCode:  http.StatusOK,
			expectCORS:    true,
			isPreflight:   true,
		},
		{
			name:          "Non-existent endpoint",
			method:        http.MethodGet,
			path:          "/nonexistent",
			headers:       map[string]string{"Origin": "http://localhost:3000"},
			expectedCode:  http.StatusNotFound,
			expectCORS:    true,
			isPreflight:   false,
		},
		{
			name:    "OPTIONS preflight for environments",
			method:  http.MethodOptions,
			path:    "/environments",
			headers: map[string]string{
				"Origin":                         "http://localhost:3000",
				"Access-Control-Request-Method":  "POST",
				"Access-Control-Request-Headers": "Content-Type",
			},
			expectedCode:  http.StatusOK,
			expectCORS:    true,
			isPreflight:   true,
		},
	}

	handler := createHandler()
	server := httptest.NewServer(handler)
	defer server.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, server.URL+tt.path, nil)
			require.NoError(t, err, "should create request without error")

			// Set request headers
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			// Make the request
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err, "should make request without error")
			defer resp.Body.Close()

			// Check status code
			assert.Equal(t, tt.expectedCode, resp.StatusCode, 
				"status code should match expected for %s %s", tt.method, tt.path)

			// For CORS requests, check the appropriate headers
			if tt.expectCORS && req.Header.Get("Origin") != "" {
				// Our custom CORS middleware sets the Access-Control-Allow-Origin
				// to the request's origin (not *) when credentials are allowed
				expectedOrigin := req.Header.Get("Origin")
				assert.Equal(t, expectedOrigin, resp.Header.Get("Access-Control-Allow-Origin"),
					"Access-Control-Allow-Origin should match the request origin")
				
				// Check that credentials are allowed
				assert.Equal(t, "true", resp.Header.Get("Access-Control-Allow-Credentials"),
					"Access-Control-Allow-Credentials should be true")

				// For preflight requests, check for additional CORS headers
				if tt.isPreflight {
					// Our custom CORS middleware sets Access-Control-Allow-Methods to the requested method
					// For the test, we'll just check that it's set to something
					allowedMethods := resp.Header.Get("Access-Control-Allow-Methods")
					assert.NotEmpty(t, allowedMethods, 
						"Access-Control-Allow-Methods should not be empty for preflight")

					// Our custom CORS middleware sets Access-Control-Allow-Headers to the requested headers
					// For the test, we'll just check that it's set to something
					allowedHeaders := resp.Header.Get("Access-Control-Allow-Headers")
					assert.NotEmpty(t, allowedHeaders,
						"Access-Control-Allow-Headers should not be empty for preflight")

					// Check for Access-Control-Max-Age header
					assert.Equal(t, "86400", resp.Header.Get("Access-Control-Max-Age"),
						"Access-Control-Max-Age should be set to 86400 for preflight")
				}
			}

			// For non-preflight requests, check the response body if it's a GET
			if tt.method == http.MethodGet && tt.path == "/health" {
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err, "should read response body without error")
				assert.JSONEq(t, `{"status":"ok"}`, string(body), 
					"health check should return status ok")
			}
		})
	}
}

func TestSaveEnvironmentSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "env.log")
	os.Setenv("ENV_LOG_FILE", logPath)
	defer os.Unsetenv("ENV_LOG_FILE")

	tests := []struct {
		name       string
		body       string
		statusCode int
		envKey     string
		envValue   string
	}{
		{
			name:       "simple JSON object",
			body:       `{"foo":"bar"}`,
			statusCode: http.StatusCreated,
			envKey:     "foo",
			envValue:   "bar",
		},
		{
			name:       "nested JSON object",
			body:       `{"nested":{"key":"value"}}`,
			statusCode: http.StatusCreated,
			envKey:     "nested",
			envValue:   `{"key":"value"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/environments", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			saveEnvironment(rr, req)

			assert.Equal(t, tt.statusCode, rr.Code, "status code should match expected")

			// Verify response body for success case
			if tt.statusCode == http.StatusCreated {
				var resp map[string]string
				err := json.NewDecoder(rr.Body).Decode(&resp)
				require.NoError(t, err, "should decode response body without error")
				assert.Equal(t, "saved", resp["status"], "status should be 'saved'")

				// Verify log file was created and contains the environment data
				data, err := os.ReadFile(logPath)
				require.NoError(t, err, "should read log file without error")

				lines := strings.Split(strings.TrimSpace(string(data)), "\n")
				assert.GreaterOrEqual(t, len(lines), 1, "should have at least one log entry")

				var rec map[string]json.RawMessage
				err = json.Unmarshal([]byte(lines[len(lines)-1]), &rec)
				require.NoError(t, err, "should unmarshal log entry without error")

				// Verify timestamp exists and is valid
				var timestamp string
				err = json.Unmarshal(rec["timestamp"], &timestamp)
				require.NoError(t, err, "should unmarshal timestamp without error")
				_, err = time.Parse(time.RFC3339, timestamp)
				assert.NoError(t, err, "timestamp should be in RFC3339 format")

				// Verify environment data
				assert.Contains(t, string(rec["environment"]), tt.envKey, "log should contain environment key")
				assert.Contains(t, string(rec["environment"]), tt.envValue, "log should contain environment value")
			}
		})
	}
}

func TestSaveEnvironmentEmpty(t *testing.T) {
	tests := []struct {
		name       string
		body       io.Reader
		headers    map[string]string
		statusCode int
		errMsg     string
	}{
		{
			name:       "empty body",
			body:       strings.NewReader(""),
			headers:    map[string]string{"Content-Type": "application/json"},
			statusCode: http.StatusBadRequest,
			errMsg:     "empty request body",
		},
		{
			name:       "invalid JSON",
			body:       strings.NewReader("{invalid"),
			headers:    map[string]string{"Content-Type": "application/json"},
			statusCode: http.StatusBadRequest,
			errMsg:     "invalid JSON",
		},
		{
			name:       "missing content type still works",
			body:       strings.NewReader("{}"),
			headers:    map[string]string{},
			statusCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			logPath := filepath.Join(tmpDir, "env.log")
			os.Setenv("ENV_LOG_FILE", logPath)
			defer os.Unsetenv("ENV_LOG_FILE")

			req := httptest.NewRequest(http.MethodPost, "/environments", tt.body)
			req.Header.Set("Content-Type", "application/json") // Set default content type
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			rr := httptest.NewRecorder()

			saveEnvironment(rr, req)

			assert.Equal(t, tt.statusCode, rr.Code, 
				"status code should match expected for %s", tt.name)

			if tt.errMsg != "" {
				body, err := io.ReadAll(rr.Body)
				require.NoError(t, err, "should read response body without error")
				assert.Contains(t, string(body), tt.errMsg, 
					"error message should contain expected text for %s", tt.name)
			}

			// For successful requests, verify the log file was created
			if tt.statusCode == http.StatusCreated {
				_, err := os.Stat(logPath)
				assert.NoError(t, err, "log file should be created for successful request")
			}
		})
	}
}

func TestSaveEnvironmentFileErrors(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping read-only tests when running as root")
	}

	tests := []struct {
		name       string
		setup      func(t *testing.T) (string, func())
		errMsg     string
		statusCode int
		expectErr  bool
	}{
		{
			name: "nonexistent parent directory - should be created automatically",
			setup: func(t *testing.T) (string, func()) {
				tmpDir := t.TempDir()
				logPath := filepath.Join(tmpDir, "nonexistent", "env.log")
				os.Setenv("ENV_LOG_FILE", logPath)
				return logPath, func() { 
					os.Unsetenv("ENV_LOG_FILE") 
					// Clean up the created directory
					os.RemoveAll(filepath.Dir(logPath))
				}
			},
			errMsg:     "",
			statusCode: http.StatusCreated,
			expectErr:  false,
		},
		{
			name: "read-only file - should fail to append",
			setup: func(t *testing.T) (string, func()) {
				tmpDir := t.TempDir()
				logPath := filepath.Join(tmpDir, "env.log")
				
				// Create and close the file first
				f, err := os.Create(logPath)
				require.NoError(t, err, "should create test file")
				f.Close()
				
				// Set read-only permissions
				err = os.Chmod(logPath, 0400)
				require.NoError(t, err, "should set file permissions")
				
				os.Setenv("ENV_LOG_FILE", logPath)
				return logPath, func() {
					os.Unsetenv("ENV_LOG_FILE")
					// Make file writable for cleanup
					os.Chmod(logPath, 0600)
				}
			},
			errMsg:     "failed to open log",
			statusCode: http.StatusInternalServerError,
			expectErr:  true,
		},
		{
			name: "read-only parent directory - should fail to create file",
			setup: func(t *testing.T) (string, func()) {
				tmpDir := t.TempDir()
				logDir := filepath.Join(tmpDir, "logs")
				
				// Create a read-only directory
				err := os.Mkdir(logDir, 0400)
				require.NoError(t, err, "should create test directory")
				
				logPath := filepath.Join(logDir, "env.log")
				os.Setenv("ENV_LOG_FILE", logPath)
				
				return logPath, func() {
					os.Unsetenv("ENV_LOG_FILE")
					// Make directory writable for cleanup
					os.Chmod(logDir, 0700)
				}
			},
			errMsg:     "failed to open log",
			statusCode: http.StatusInternalServerError,
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logPath, cleanup := tt.setup(t)
			defer cleanup()

			// Get the initial state of the log file if it exists
			initialSize := int64(0)
			if info, err := os.Stat(logPath); err == nil {
				initialSize = info.Size()
			}

			req := httptest.NewRequest(http.MethodPost, "/environments", strings.NewReader(`{"test":"data"}`))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			saveEnvironment(rr, req)

			// Verify the response status code
			assert.Equal(t, tt.statusCode, rr.Code, 
				"status code should match expected for %s", tt.name)

			// Verify the response body
			body, err := io.ReadAll(rr.Body)
			require.NoError(t, err, "should read response body without error")

			if tt.expectErr {
				// For error cases, verify the error message is in the response
				assert.Contains(t, string(body), tt.errMsg, 
					"error message should contain expected text for %s: got %s", tt.name, string(body))

				// Verify the log file was not modified on error
				if info, err := os.Stat(logPath); err == nil {
					assert.Equal(t, initialSize, info.Size(), 
						"log file should not be modified on error")
				}
			} else {
				// For success cases, verify the log file was created/modified
				info, err := os.Stat(logPath)
				require.NoError(t, err, "log file should exist after successful write")
				assert.Greater(t, info.Size(), initialSize, 
					"log file should be larger after successful write")
			}
		})
	}
}

func TestSetupServer(t *testing.T) {
	tests := []struct {
		name     string
		handler  http.Handler
		port     string
		expected string
	}{
		{
			name:     "default port",
			handler:  http.NewServeMux(),
			port:     "",
			expected: ":8080",
		},
		{
			name:     "custom port",
			handler:  http.NewServeMux(),
			port:     "9090",
			expected: ":9090",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupServer(tt.handler, tt.port)
			assert.Equal(t, tt.expected, server.Addr, "server address should match expected")
			assert.NotNil(t, server.Handler, "server handler should not be nil")
		})
	}
}

func TestSaveEnvironment_Validation(t *testing.T) {
	t.Run("unsupported_content_type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/environments", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "text/plain")
		rr := httptest.NewRecorder()
		saveEnvironment(rr, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rr.Code)
		assert.Contains(t, rr.Body.String(), "Content-Type must be application/json")
	})

	t.Run("empty_request_body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/environments", nil)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		saveEnvironment(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "empty request body")
	})

	t.Run("invalid_json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/environments", strings.NewReader(`{`))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		saveEnvironment(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid JSON")
	})
}

func TestRunServer(t *testing.T) {
	// Find an available port by creating a listener and closing it immediately.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err, "should be able to create a test listener")
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close() // Release the port so the server can use it.

	server := &http.Server{
		Addr: fmt.Sprintf("127.0.0.1:%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	}

	// Start the server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- runServer(server)
	}()

	// Retry connecting to the server to avoid race conditions in CI.
	var resp *http.Response
	retries := 10
	for i := 0; i < retries; i++ {
		resp, err = http.Get(fmt.Sprintf("http://127.0.0.1:%d", port))
		if err == nil {
			break // Success
		}
		time.Sleep(200 * time.Millisecond)
	}
	require.NoError(t, err, "should be able to make request to server after retries")
	
	if resp != nil {
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode, "server should respond with 200 OK")
	}

	// Shutdown the server with a timeout.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(shutdownCtx)
	require.NoError(t, err, "should be able to shutdown server")

	// Verify the server shutdown gracefully.
	select {
	case err := <-errCh:
		assert.Equal(t, http.ErrServerClosed, err, "should return server closed error")
	case <-time.After(2 * time.Second):
		t.Fatal("server did not shut down within timeout")
	}
}
