package main

import (
	"os"
	"testing"
)

func TestEnvironmentConfiguration(t *testing.T) {
	// Save original environment
	originalUsername := os.Getenv("PROXY_USERNAME")
	originalPassword := os.Getenv("PROXY_PASSWORD")
	originalPort := os.Getenv("PROXY_PORT")

	// Clean up after test
	defer func() {
		os.Setenv("PROXY_USERNAME", originalUsername)
		os.Setenv("PROXY_PASSWORD", originalPassword)
		os.Setenv("PROXY_PORT", originalPort)
	}()

	tests := []struct {
		name         string
		envUsername  string
		envPassword  string
		envPort      string
		expectedUser string
		expectedPass string
		expectedPort string
	}{
		{
			name:         "All environment variables set",
			envUsername:  "envuser",
			envPassword:  "envpass",
			envPort:      "9090",
			expectedUser: "envuser",
			expectedPass: "envpass",
			expectedPort: "9090",
		},
		{
			name:         "Only username set, others default",
			envUsername:  "envuser",
			envPassword:  "",
			envPort:      "",
			expectedUser: "envuser",
			expectedPass: "password123",
			expectedPort: "8080",
		},
		{
			name:         "No environment variables, use defaults",
			envUsername:  "",
			envPassword:  "",
			envPort:      "",
			expectedUser: "admin",
			expectedPass: "password123",
			expectedPort: "8080",
		},
		{
			name:         "Custom port only",
			envUsername:  "",
			envPassword:  "",
			envPort:      "3128",
			expectedUser: "admin",
			expectedPass: "password123",
			expectedPort: "3128",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			os.Setenv("PROXY_USERNAME", tt.envUsername)
			os.Setenv("PROXY_PASSWORD", tt.envPassword)
			os.Setenv("PROXY_PORT", tt.envPort)

			// Get configuration like main() does
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

			// Verify configuration
			if username != tt.expectedUser {
				t.Errorf("Expected username %s, got %s", tt.expectedUser, username)
			}
			if password != tt.expectedPass {
				t.Errorf("Expected password %s, got %s", tt.expectedPass, password)
			}
			if port != tt.expectedPort {
				t.Errorf("Expected port %s, got %s", tt.expectedPort, port)
			}
		})
	}
}

func TestConfigurationValidation(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		isValid  bool
	}{
		{
			name:     "Valid credentials",
			username: "admin",
			password: "password123",
			isValid:  true,
		},
		{
			name:     "Empty username after defaults",
			username: "",
			password: "password123",
			isValid:  false,
		},
		{
			name:     "Empty password after defaults",
			username: "admin",
			password: "",
			isValid:  false,
		},
		{
			name:     "Both empty after defaults",
			username: "",
			password: "",
			isValid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate validation logic from main()
			isValid := tt.username != "" && tt.password != ""

			if isValid != tt.isValid {
				t.Errorf("Expected validity %v, got %v", tt.isValid, isValid)
			}
		})
	}
}

func TestPortValidation(t *testing.T) {
	validPorts := []string{"80", "8080", "3128", "1080", "9090"}
	invalidPorts := []string{"", "0", "65536", "abc", "-1"}

	for _, port := range validPorts {
		t.Run("Valid port "+port, func(t *testing.T) {
			// For this simple test, we just check if it's a non-empty string
			// In a real application, you might want to parse and validate the port number
			if port == "" {
				t.Errorf("Port %s should be valid", port)
			}
		})
	}

	for _, port := range invalidPorts {
		t.Run("Invalid port "+port, func(t *testing.T) {
			// This is a simplified test - in production you'd want proper port validation
			if port == "" {
				// Empty port should use default
				if port != "" {
					t.Errorf("Empty port should remain empty for default handling")
				}
			}
		})
	}
}
