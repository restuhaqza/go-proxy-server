package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewProxyServer(t *testing.T) {
	username := "testuser"
	password := "testpass"
	port := "8080"

	proxy := NewProxyServer(username, password, port)

	if proxy.username != username {
		t.Errorf("Expected username %s, got %s", username, proxy.username)
	}
	if proxy.password != password {
		t.Errorf("Expected password %s, got %s", password, proxy.password)
	}
	if proxy.port != port {
		t.Errorf("Expected port %s, got %s", port, proxy.port)
	}
}

func TestAuthenticateRequest(t *testing.T) {
	proxy := NewProxyServer("admin", "password123", "8080")

	tests := []struct {
		name           string
		auth           string
		expectedResult bool
	}{
		{
			name:           "Valid credentials",
			auth:           "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:password123")),
			expectedResult: true,
		},
		{
			name:           "Invalid username",
			auth:           "Basic " + base64.StdEncoding.EncodeToString([]byte("wrong:password123")),
			expectedResult: false,
		},
		{
			name:           "Invalid password",
			auth:           "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:wrong")),
			expectedResult: false,
		},
		{
			name:           "No auth header",
			auth:           "",
			expectedResult: false,
		},
		{
			name:           "Invalid auth type",
			auth:           "Bearer token123",
			expectedResult: false,
		},
		{
			name:           "Malformed basic auth",
			auth:           "Basic invalid-base64",
			expectedResult: false,
		},
		{
			name:           "Basic auth without colon",
			auth:           "Basic " + base64.StdEncoding.EncodeToString([]byte("adminpassword")),
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com", nil)
			if tt.auth != "" {
				req.Header.Set("Proxy-Authorization", tt.auth)
			}

			result := proxy.authenticateRequest(req)
			if result != tt.expectedResult {
				t.Errorf("Expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestHandleHTTP_Authentication(t *testing.T) {
	proxy := NewProxyServer("admin", "password123", "8080")

	// Test without authentication
	t.Run("No authentication", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com", nil)
		w := httptest.NewRecorder()

		proxy.handleHTTP(w, req)

		if w.Code != http.StatusProxyAuthRequired {
			t.Errorf("Expected status %d, got %d", http.StatusProxyAuthRequired, w.Code)
		}

		if w.Header().Get("Proxy-Authenticate") == "" {
			t.Error("Expected Proxy-Authenticate header to be set")
		}
	})

	// Test with invalid authentication
	t.Run("Invalid authentication", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com", nil)
		req.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("wrong:credentials")))
		w := httptest.NewRecorder()

		proxy.handleHTTP(w, req)

		if w.Code != http.StatusProxyAuthRequired {
			t.Errorf("Expected status %d, got %d", http.StatusProxyAuthRequired, w.Code)
		}
	})
}

func TestHandleHTTP_ValidRequest(t *testing.T) {
	// Create a test server to simulate the target
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"message": "success", "method": "`+r.Method+`"}`)
	}))
	defer targetServer.Close()

	proxy := NewProxyServer("admin", "password123", "8080")

	tests := []struct {
		name   string
		method string
	}{
		{"GET request", "GET"},
		{"POST request", "POST"},
		{"PUT request", "PUT"},
		{"DELETE request", "DELETE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body io.Reader
			if tt.method == "POST" || tt.method == "PUT" {
				body = strings.NewReader(`{"test": "data"}`)
			}

			req := httptest.NewRequest(tt.method, targetServer.URL, body)
			req.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:password123")))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			proxy.handleHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
			}

			// Check if proxy-specific headers are removed
			if req.Header.Get("Proxy-Authorization") != "" {
				t.Error("Proxy-Authorization header should be removed")
			}

			// Check response
			responseBody := w.Body.String()
			if !strings.Contains(responseBody, tt.method) {
				t.Errorf("Response should contain method %s", tt.method)
			}
		})
	}
}

func TestHandleHTTP_InvalidURL(t *testing.T) {
	proxy := NewProxyServer("admin", "password123", "8080")

	req := httptest.NewRequest("GET", "http://invalid-url-that-does-not-exist.local", nil)
	req.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:password123")))
	w := httptest.NewRecorder()

	proxy.handleHTTP(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("Expected status %d, got %d", http.StatusBadGateway, w.Code)
	}
}

func TestHandleHTTPS_Authentication(t *testing.T) {
	proxy := NewProxyServer("admin", "password123", "8080")

	// Test CONNECT without authentication
	t.Run("CONNECT without auth", func(t *testing.T) {
		req := httptest.NewRequest("CONNECT", "example.com:443", nil)
		w := httptest.NewRecorder()

		proxy.handleHTTPS(w, req)

		if w.Code != http.StatusProxyAuthRequired {
			t.Errorf("Expected status %d, got %d", http.StatusProxyAuthRequired, w.Code)
		}
	})

	// Test CONNECT with invalid authentication
	t.Run("CONNECT with invalid auth", func(t *testing.T) {
		req := httptest.NewRequest("CONNECT", "example.com:443", nil)
		req.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("wrong:credentials")))
		w := httptest.NewRecorder()

		proxy.handleHTTPS(w, req)

		if w.Code != http.StatusProxyAuthRequired {
			t.Errorf("Expected status %d, got %d", http.StatusProxyAuthRequired, w.Code)
		}
	})
}

func TestHandleHTTPS_InvalidHost(t *testing.T) {
	proxy := NewProxyServer("admin", "password123", "8080")

	req := httptest.NewRequest("CONNECT", "invalid-host-that-does-not-exist.local:443", nil)
	req.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:password123")))
	w := httptest.NewRecorder()

	proxy.handleHTTPS(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("Expected status %d, got %d", http.StatusBadGateway, w.Code)
	}
}

func TestServeHTTP(t *testing.T) {
	proxy := NewProxyServer("admin", "password123", "8080")

	// Test HTTP method routing
	t.Run("HTTP GET routing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com", nil)
		w := httptest.NewRecorder()

		proxy.ServeHTTP(w, req)

		// Should be handled by handleHTTP and require auth
		if w.Code != http.StatusProxyAuthRequired {
			t.Errorf("Expected status %d, got %d", http.StatusProxyAuthRequired, w.Code)
		}
	})

	// Test CONNECT method routing
	t.Run("CONNECT routing", func(t *testing.T) {
		req := httptest.NewRequest("CONNECT", "example.com:443", nil)
		w := httptest.NewRecorder()

		proxy.ServeHTTP(w, req)

		// Should be handled by handleHTTPS and require auth
		if w.Code != http.StatusProxyAuthRequired {
			t.Errorf("Expected status %d, got %d", http.StatusProxyAuthRequired, w.Code)
		}
	})
}

func TestHeaderHandling(t *testing.T) {
	// Create a test server that echoes headers
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Echo-User-Agent", r.Header.Get("User-Agent"))
		w.Header().Set("Echo-Custom-Header", r.Header.Get("Custom-Header"))
		w.WriteHeader(http.StatusOK)
	}))
	defer targetServer.Close()

	proxy := NewProxyServer("admin", "password123", "8080")

	req := httptest.NewRequest("GET", targetServer.URL, nil)
	req.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:password123")))
	req.Header.Set("User-Agent", "TestAgent/1.0")
	req.Header.Set("Custom-Header", "CustomValue")
	req.Header.Set("Proxy-Connection", "keep-alive")
	w := httptest.NewRecorder()

	proxy.handleHTTP(w, req)

	// Check that custom headers are preserved
	if w.Header().Get("Echo-User-Agent") != "TestAgent/1.0" {
		t.Error("User-Agent header was not properly forwarded")
	}
	if w.Header().Get("Echo-Custom-Header") != "CustomValue" {
		t.Error("Custom header was not properly forwarded")
	}

	// Check that proxy-specific headers were removed
	if req.Header.Get("Proxy-Authorization") != "" {
		t.Error("Proxy-Authorization header should be removed")
	}
	if req.Header.Get("Proxy-Connection") != "" {
		t.Error("Proxy-Connection header should be removed")
	}
}

func TestRequestBodyHandling(t *testing.T) {
	// Create a test server that echoes the request body
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}))
	defer targetServer.Close()

	proxy := NewProxyServer("admin", "password123", "8080")

	testBody := `{"test": "data", "number": 123}`
	req := httptest.NewRequest("POST", targetServer.URL, bytes.NewReader([]byte(testBody)))
	req.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:password123")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	proxy.handleHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	responseBody := w.Body.String()
	if responseBody != testBody {
		t.Errorf("Expected body %s, got %s", testBody, responseBody)
	}
}

func BenchmarkAuthenticateRequest(b *testing.B) {
	proxy := NewProxyServer("admin", "password123", "8080")
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:password123")))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		proxy.authenticateRequest(req)
	}
}

func BenchmarkHandleHTTP(b *testing.B) {
	// Create a simple test server
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	}))
	defer targetServer.Close()

	proxy := NewProxyServer("admin", "password123", "8080")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", targetServer.URL, nil)
		req.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:password123")))
		w := httptest.NewRecorder()

		proxy.handleHTTP(w, req)
	}
}

// Integration test for the complete proxy flow
func TestProxyIntegration(t *testing.T) {
	// Create a test backend server
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Backend-Header", "backend-value")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Backend response: %s %s", r.Method, r.URL.Path)
	}))
	defer backendServer.Close()

	// Create proxy server
	proxy := NewProxyServer("testuser", "testpass", "0")

	// Create request to backend through proxy
	req, err := http.NewRequest("GET", backendServer.URL+"/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set proxy authorization
	req.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("testuser:testpass")))

	// Make request through proxy
	w := httptest.NewRecorder()
	proxy.ServeHTTP(w, req)

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check that backend headers are preserved
	if w.Header().Get("Backend-Header") != "backend-value" {
		t.Error("Backend header was not preserved")
	}

	// Check response body
	expectedBody := "Backend response: GET /test"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}
