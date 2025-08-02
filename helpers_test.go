package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestHelper provides utilities for testing the proxy server
type TestHelper struct {
	ProxyServer   *ProxyServer
	BackendServer *httptest.Server
	ProxyHandler  http.Handler
}

// NewTestHelper creates a new test helper with a proxy and backend server
func NewTestHelper(username, password string) *TestHelper {
	// Create a simple backend server
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test-Backend", "true")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Backend response"))
	}))

	proxy := NewProxyServer(username, password, "8080")

	return &TestHelper{
		ProxyServer:   proxy,
		BackendServer: backendServer,
		ProxyHandler:  proxy,
	}
}

// Close cleans up test resources
func (th *TestHelper) Close() {
	if th.BackendServer != nil {
		th.BackendServer.Close()
	}
}

// CreateAuthenticatedRequest creates a request with proper proxy authentication
func (th *TestHelper) CreateAuthenticatedRequest(method, url string) *http.Request {
	req := httptest.NewRequest(method, url, nil)
	req.Header.Set("Proxy-Authorization", th.GetBasicAuth())
	return req
}

// CreateUnauthenticatedRequest creates a request without authentication
func (th *TestHelper) CreateUnauthenticatedRequest(method, url string) *http.Request {
	return httptest.NewRequest(method, url, nil)
}

// GetBasicAuth returns the basic auth header value for the proxy
func (th *TestHelper) GetBasicAuth() string {
	return CreateBasicAuth(th.ProxyServer.username, th.ProxyServer.password)
}

// CreateBasicAuth creates a basic auth header value
func CreateBasicAuth(username, password string) string {
	return "Basic " + base64Encode(username+":"+password)
}

// base64Encode encodes a string to base64
func base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func TestTestHelper(t *testing.T) {
	helper := NewTestHelper("testuser", "testpass")
	defer helper.Close()

	if helper.ProxyServer == nil {
		t.Error("ProxyServer should not be nil")
	}

	if helper.BackendServer == nil {
		t.Error("BackendServer should not be nil")
	}

	if helper.ProxyHandler == nil {
		t.Error("ProxyHandler should not be nil")
	}

	// Test authenticated request creation
	authReq := helper.CreateAuthenticatedRequest("GET", "http://example.com")
	if authReq.Header.Get("Proxy-Authorization") == "" {
		t.Error("Authenticated request should have Proxy-Authorization header")
	}

	// Test unauthenticated request creation
	unauthReq := helper.CreateUnauthenticatedRequest("GET", "http://example.com")
	if unauthReq.Header.Get("Proxy-Authorization") != "" {
		t.Error("Unauthenticated request should not have Proxy-Authorization header")
	}
}

func TestServerTimeout(t *testing.T) {
	// Create a slow backend server
	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sleep longer than the proxy timeout (30 seconds)
		time.Sleep(35 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Slow response"))
	}))
	defer slowServer.Close()

	proxy := NewProxyServer("admin", "password123", "8080")

	// Create request to slow server
	req := httptest.NewRequest("GET", slowServer.URL, nil)
	req.Header.Set("Proxy-Authorization", "Basic "+base64Encode("admin:password123"))

	// Create a context with timeout shorter than the server sleep
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	// This should timeout
	proxy.handleHTTP(w, req)

	// The request should fail due to context timeout or proxy timeout
	if w.Code == http.StatusOK {
		t.Error("Request should have failed due to timeout")
	}
}

func TestConcurrentRequests(t *testing.T) {
	helper := NewTestHelper("admin", "password123")
	defer helper.Close()

	// Number of concurrent requests to test
	numRequests := 10
	results := make(chan int, numRequests)

	// Start concurrent requests
	for i := 0; i < numRequests; i++ {
		go func() {
			req := helper.CreateAuthenticatedRequest("GET", helper.BackendServer.URL)
			w := httptest.NewRecorder()

			helper.ProxyHandler.ServeHTTP(w, req)
			results <- w.Code
		}()
	}

	// Collect results
	successCount := 0
	for i := 0; i < numRequests; i++ {
		statusCode := <-results
		if statusCode == http.StatusOK {
			successCount++
		}
	}

	// All requests should succeed
	if successCount != numRequests {
		t.Errorf("Expected %d successful requests, got %d", numRequests, successCount)
	}
}

func TestMemoryUsage(t *testing.T) {
	// Simple test to ensure we're not leaking memory during normal operation
	helper := NewTestHelper("admin", "password123")
	defer helper.Close()

	// Make multiple requests to check for obvious memory leaks
	for i := 0; i < 100; i++ {
		req := helper.CreateAuthenticatedRequest("GET", helper.BackendServer.URL)
		w := httptest.NewRecorder()

		helper.ProxyHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d failed with status %d", i, w.Code)
		}
	}

	// In a real scenario, you might want to use runtime.ReadMemStats()
	// to check for memory leaks, but for this simple test, we just ensure
	// all requests complete successfully
}

func TestLargeRequestBody(t *testing.T) {
	// Create a backend that echoes the request body size
	echoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		// Read and count bytes
		count := 0
		buffer := make([]byte, 1024)
		for {
			n, err := r.Body.Read(buffer)
			count += n
			if err != nil {
				break
			}
		}
		w.Write([]byte(fmt.Sprintf("Received %d bytes", count)))
	}))
	defer echoServer.Close()

	proxy := NewProxyServer("admin", "password123", "8080")

	// Create a large request body (1MB)
	largeBody := make([]byte, 1024*1024)
	for i := range largeBody {
		largeBody[i] = byte(i % 256)
	}

	req := httptest.NewRequest("POST", echoServer.URL, bytes.NewReader(largeBody))
	req.Header.Set("Proxy-Authorization", "Basic "+base64Encode("admin:password123"))
	req.Header.Set("Content-Type", "application/octet-stream")

	w := httptest.NewRecorder()
	proxy.handleHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// The response should indicate that all bytes were received
	expectedResponse := "Received 1048576 bytes"
	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response %s, got %s", expectedResponse, w.Body.String())
	}
}
