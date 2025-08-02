#!/bin/bash

# Local development setup script

set -e

echo "=== Go Proxy Server Development Setup ==="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    echo "Visit: https://golang.org/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if ! printf '%s\n%s' "$REQUIRED_VERSION" "$GO_VERSION" | sort -C -V; then
    echo "❌ Go version $GO_VERSION is too old. Please upgrade to Go $REQUIRED_VERSION or later."
    exit 1
fi

echo "✅ Go version $GO_VERSION detected"

# Check if Docker is installed (optional)
if command -v docker &> /dev/null; then
    echo "✅ Docker detected"
    DOCKER_AVAILABLE=true
else
    echo "⚠️  Docker not found. Docker-based development will not be available."
    DOCKER_AVAILABLE=false
fi

# Install development dependencies
echo "📦 Installing development dependencies..."

# Go tools
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Download project dependencies
echo "📦 Downloading project dependencies..."
go mod download
go mod verify

# Create .env file from example if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from example..."
    cp .env.example .env
    echo "✅ .env file created. Please review and modify as needed."
fi

# Setup git hooks (optional)
if [ -d .git ]; then
    echo "🔧 Setting up git hooks..."
    
    # Pre-commit hook
    cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
# Run tests and linting before commit

set -e

echo "Running pre-commit checks..."

# Run tests
go test ./...

# Run static analysis
if command -v staticcheck &> /dev/null; then
    staticcheck ./...
fi

# Run go vet
go vet ./...

# Check formatting
if [ -n "$(gofmt -l .)" ]; then
    echo "Code is not properly formatted. Run 'go fmt ./...'"
    exit 1
fi

echo "All pre-commit checks passed!"
EOF
    
    chmod +x .git/hooks/pre-commit
    echo "✅ Git pre-commit hook installed"
fi

echo ""
echo "🎉 Development environment setup complete!"
echo ""
echo "Next steps:"
echo "1. Review and modify .env file with your preferred settings"
echo "2. Run the application:"
echo "   go run main.go"
echo ""
echo "3. Or build and run:"
echo "   go build -o proxy-server ."
echo "   ./proxy-server"
echo ""

if [ "$DOCKER_AVAILABLE" = true ]; then
    echo "4. Or use Docker:"
    echo "   docker-compose up --build"
    echo ""
fi

echo "5. Test the proxy:"
echo "   curl --proxy http://admin:password123@localhost:8080 http://httpbin.org/ip"
echo ""
echo "For more information, see README.md"
