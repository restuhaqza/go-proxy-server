#!/bin/bash

# Build script for Go Proxy Server

set -e

echo "=== Building Go Proxy Server ==="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed"
    exit 1
fi

# Get dependencies
echo "Getting dependencies..."
go mod tidy

# Build the application
echo "Building application..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o proxy-server .

echo "Build completed successfully!"
echo "Binary: ./proxy-server"

# Make the binary executable
chmod +x proxy-server

echo "Ready to run with: ./proxy-server"
