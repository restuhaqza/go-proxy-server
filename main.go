package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// ProxyServer represents the HTTP proxy server
type ProxyServer struct {
	username string
	password string
	port     string
}

// NewProxyServer creates a new proxy server instance
func NewProxyServer(username, password, port string) *ProxyServer {
	return &ProxyServer{
		username: username,
		password: password,
		port:     port,
	}
}

// authenticateRequest checks if the request has valid Basic Auth credentials
func (ps *ProxyServer) authenticateRequest(r *http.Request) bool {
	auth := r.Header.Get("Proxy-Authorization")
	if auth == "" {
		return false
	}

	// Check if it's Basic authentication
	if !strings.HasPrefix(auth, "Basic ") {
		return false
	}

	// Decode the base64 encoded credentials
	payload, err := base64.StdEncoding.DecodeString(auth[6:])
	if err != nil {
		return false
	}

	// Split username and password
	credentials := strings.SplitN(string(payload), ":", 2)
	if len(credentials) != 2 {
		return false
	}

	username, password := credentials[0], credentials[1]
	return username == ps.username && password == ps.password
}

// handleHTTP handles HTTP requests through the proxy
func (ps *ProxyServer) handleHTTP(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	if !ps.authenticateRequest(r) {
		w.Header().Set("Proxy-Authenticate", "Basic realm=\"Proxy Server\"")
		http.Error(w, "Proxy Authentication Required", http.StatusProxyAuthRequired)
		return
	}

	// Remove proxy-specific headers
	r.Header.Del("Proxy-Authorization")
	r.Header.Del("Proxy-Connection")

	// Create HTTP client
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create new request
	proxyReq, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}

	// Copy headers
	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	// Make the request
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Error making proxy request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Set status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error copying response body: %v", err)
	}
}

// handleHTTPS handles HTTPS CONNECT requests
func (ps *ProxyServer) handleHTTPS(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	if !ps.authenticateRequest(r) {
		w.Header().Set("Proxy-Authenticate", "Basic realm=\"Proxy Server\"")
		http.Error(w, "Proxy Authentication Required", http.StatusProxyAuthRequired)
		return
	}

	// Get the destination host
	destConn, err := net.DialTimeout("tcp", r.Host, 30*time.Second)
	if err != nil {
		http.Error(w, "Error connecting to destination", http.StatusBadGateway)
		return
	}
	defer destConn.Close()

	// Send 200 Connection established
	w.WriteHeader(http.StatusOK)

	// Get the underlying connection
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, "Error hijacking connection", http.StatusInternalServerError)
		return
	}
	defer clientConn.Close()

	// Start copying data between client and destination
	go func() {
		defer destConn.Close()
		defer clientConn.Close()
		io.Copy(destConn, clientConn)
	}()

	io.Copy(clientConn, destConn)
}

// ServeHTTP implements the http.Handler interface
func (ps *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.String())

	if r.Method == "CONNECT" {
		ps.handleHTTPS(w, r)
	} else {
		ps.handleHTTP(w, r)
	}
}

// Start starts the proxy server
func (ps *ProxyServer) Start() error {
	server := &http.Server{
		Addr:    ":" + ps.port,
		Handler: ps,
	}

	log.Printf("Starting HTTP Proxy Server on port %s", ps.port)
	log.Printf("Username: %s", ps.username)
	log.Printf("Server ready to accept connections...")

	return server.ListenAndServe()
}

func main() {
	// Get configuration from environment variables
	username := os.Getenv("PROXY_USERNAME")
	password := os.Getenv("PROXY_PASSWORD")
	port := os.Getenv("PROXY_PORT")

	// Set default values if not provided
	if username == "" {
		username = "admin"
	}
	if password == "" {
		password = "password123"
	}
	if port == "" {
		port = "8080"
	}

	// Validate configuration
	if username == "" || password == "" {
		log.Fatal("Username and password are required")
	}

	// Create and start proxy server
	proxy := NewProxyServer(username, password, port)

	fmt.Printf("=== HTTP Proxy Server ===\n")
	fmt.Printf("Port: %s\n", port)
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Password: %s\n", strings.Repeat("*", len(password)))
	fmt.Printf("========================\n\n")

	log.Fatal(proxy.Start())
}
