<div align="center">

# 🚀 HTTP Proxy Server with Go

**A high-performance HTTP/HTTPS proxy server written in Go with Basic Auth authentication and production-ready Docker deployment**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Supported-2496ED?style=flat&logo=docker)](https://www.docker.com)
[![GitHub Actions](https://img.shields.io/badge/CI%2FCD-GitHub%20Actions-2088FF?style=flat&logo=github-actions)](https://github.com/features/actions)

A lightweight, secure, and easy-to-deploy proxy server perfect for development, testing, and production environments.

</div>

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🌐 **HTTP/HTTPS Proxy** | Full support for both HTTP and HTTPS protocols |
| 🔐 **Basic Auth** | Secure username/password authentication for all requests |
| 🔗 **CONNECT Method** | Native support for HTTPS tunneling |
| ⚙️ **Environment Config** | Simple configuration via environment variables |
| 🐳 **Docker Ready** | Multi-stage Docker build for optimized images |
| 💚 **Health Check** | Built-in health monitoring endpoint |
| 📝 **Request Logging** | Detailed logging for monitoring and debugging |
| 👤 **Non-root User** | Security-focused container running as non-root |

---

## 📖 Usage

### 🖥️ 1. Development (Local)

```bash
# Clone or download project
cd go-proxy-server

# Install dependencies
go mod tidy

# Set environment variables (optional)
export PROXY_USERNAME=admin
export PROXY_PASSWORD=mypassword
export PROXY_PORT=8080

# Run application
go run main.go
```

### 🔨 2. Build Binary

```bash
# Build binary
go build -o proxy-server

# Run binary
./proxy-server
```

### 🐳 3. Docker Deployment

#### Using Docker Build

```bash
# Build Docker image
docker build -t go-proxy-server .

# Run container with default credentials
docker run -d \
  --name proxy-server \
  -p 8080:8080 \
  go-proxy-server

# Run container with custom credentials
docker run -d \
  --name proxy-server \
  -p 8080:8080 \
  -e PROXY_USERNAME=myuser \
  -e PROXY_PASSWORD=mypassword \
  go-proxy-server
```

#### Using Docker Compose

```bash
# Copy environment file
cp .env.example .env

# Edit .env file with desired credentials
nano .env

# Start with docker-compose
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

---

## ⚙️ Configuration

The application uses environment variables for configuration:

| Variable | Default | Description |
|----------|---------|-------------|
| `PROXY_USERNAME` | `admin` | Username for proxy authentication |
| `PROXY_PASSWORD` | `password123` | Password for proxy authentication |
| `PROXY_PORT` | `8080` | Proxy server port |

---

## 🔌 How to Use Proxy

### 🌍 1. Configure Proxy in Browser/Application

| Setting | Value |
|---------|-------|
| **Proxy Type** | HTTP |
| **Server** | `localhost` (or your server IP) |
| **Port** | `8080` (or your configured port) |
| **Username** | As per `PROXY_USERNAME` |
| **Password** | As per `PROXY_PASSWORD` |

### 🧪 2. Test with curl

```bash
# Test HTTP request
curl -v \
  --proxy http://admin:mypassword@localhost:8080 \
  http://httpbin.org/ip

# Test HTTPS request
curl -v \
  --proxy http://admin:mypassword@localhost:8080 \
  https://httpbin.org/ip
```

### 📥 3. Test with wget

```bash
# Set proxy environment
export http_proxy=http://admin:mypassword@localhost:8080
export https_proxy=http://admin:mypassword@localhost:8080

# Test request
wget http://httpbin.org/ip
```

---

## 📊 Monitoring

### 🏥 Health Check

The application has a built-in health check that can be used for monitoring:

```bash
# Check health (for monitoring only, not for proxy)
curl http://localhost:8080
```

### 📋 Logs

The application will display logs for each request:

```bash
2024/01/01 12:00:00 127.0.0.1:12345 GET http://example.com
2024/01/01 12:00:01 127.0.0.1:12346 CONNECT example.com:443
```

---

## 🔒 Security

| Security Feature | Description |
|-----------------|-------------|
| 🔐 **Basic Auth** | Authentication required for all requests |
| 👤 **Non-root User** | Container runs as non-root user |
| 🏔️ **Alpine Linux** | Minimal, secure base image |
| 📦 **No Extra Packages** | Only necessary dependencies installed |

---

## 🛠️ Troubleshooting

### ❌ Error "Proxy Authentication Required"

**Solution**: Make sure the username and password used match the server configuration.

### ⏱️ Connection Timeout

**Solution**: The application has a 30-second timeout for each request. For longer requests, adjust the timeout in the code.

### 🔌 Port Already in Use

**Solution**: Change the port in the `PROXY_PORT` environment variable or use a different port when running the container.

---

## 🔄 CI/CD and Versioning

### 🤖 GitHub Actions Workflows

This project is equipped with GitHub Actions workflows for automation:

#### 1. CI/CD Pipeline (`ci-cd.yml`)
- **Trigger**: Push to `main`/`develop`, Pull Request, or Tag
- **Process**:
  - ✅ Test and build application
  - ✅ Security scan with staticcheck
  - ✅ Build multi-platform Docker images
  - ✅ Push to GitHub Container Registry
  - ✅ Vulnerability scan with Trivy

#### 2. Release Pipeline (`release.yml`)
- **Trigger**: Push tag with format `v*` (example: `v1.0.0`)
- **Process**:
  - ✅ Build binary for multiple platforms (Linux, macOS, Windows)
  - ✅ Create checksums
  - ✅ Build and push Docker images
  - ✅ Create GitHub release with binaries
  - ✅ Generate automatic release notes

#### 3. Security Scan (`security.yml`)
- **Trigger**: Daily schedule or manual
- **Process**:
  - ✅ Dependency scanning with Gosec
  - ✅ Container vulnerability scanning
  - ✅ Upload results to GitHub Security tab

#### 4. Docker Hub Release (`dockerhub.yml`)
- **Trigger**: Push tag `v*` or manual
- **Process**:
  - ✅ Build and push to Docker Hub
  - ✅ Update description automatically

### 🔑 Setup GitHub Repository

1. **Required Secrets**:
```bash
# For Docker Hub (optional)
DOCKERHUB_USERNAME=your-dockerhub-username
DOCKERHUB_TOKEN=your-dockerhub-access-token
```

2. **Branch Protection Rules** (Recommended):
   - Require pull request reviews
   - Require status checks (CI tests)
   - Require branches to be up to date

### 📦 Versioning and Release

#### Using Release Script
```bash
# Patch release (1.0.0 -> 1.0.1)
./scripts/release.sh patch

# Minor release (1.0.1 -> 1.1.0)
./scripts/release.sh minor

# Major release (1.1.0 -> 2.0.0)
./scripts/release.sh major
```

#### Manual Release
```bash
# Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# GitHub Actions will automatically:
# - Build application for multiple platforms
# - Create Docker images
# - Create GitHub release
# - Generate release notes
```

### 🐋 Registry Images

#### GitHub Container Registry
```bash
# Pull latest
docker pull ghcr.io/your-username/go-proxy-server:latest

# Pull specific version
docker pull ghcr.io/your-username/go-proxy-server:v1.0.0
```

#### Docker Hub (optional)
```bash
# Pull latest
docker pull your-username/go-proxy-server:latest

# Pull specific version
docker pull your-username/go-proxy-server:v1.0.0
```

---

## 💻 Development

### 🚀 Quick Setup
```bash
# Setup development environment
./scripts/setup-dev.sh

# Install git hooks, dependencies, and tools
```

### 📁 Project Structure

```
go-proxy-server/
├── .github/
│   ├── workflows/           # GitHub Actions workflows
│   │   ├── ci-cd.yml       # Main CI/CD pipeline
│   │   ├── release.yml     # Release automation
│   │   ├── security.yml    # Security scanning
│   │   └── dockerhub.yml   # Docker Hub publishing
│   ├── ISSUE_TEMPLATE/     # Issue templates
│   ├── dependabot.yml     # Dependency updates
│   └── pull_request_template.md
├── scripts/
│   ├── release.sh          # Release automation script
│   └── setup-dev.sh        # Development setup
├── main.go                 # Main application
├── go.mod                  # Go modules
├── Dockerfile              # Docker configuration
├── docker-compose.yml      # Docker Compose
├── .env.example            # Environment example
├── .gitignore             # Git ignore
└── README.md              # Documentation
```

### 🔄 Development Workflow

#### 1. Setup environment
```bash
./scripts/setup-dev.sh
```

#### 2. Create feature branch
```bash
git checkout -b feature/your-feature
```

#### 3. Development
```bash
# Run locally
go run main.go

# Or with Docker
docker-compose up --build
```

#### 4. Testing
```bash
# Run tests
go test ./...

# Run linting
staticcheck ./...
golangci-lint run
```

#### 5. Commit and Push
```bash
git add .
git commit -m "feat: add your feature"
git push origin feature/your-feature
```

#### 6. Create Pull Request
   - GitHub will automatically run CI checks
   - Review and merge to main

#### 7. Release
```bash
# After merge to main
git checkout main
git pull origin main
./scripts/release.sh patch  # or minor/major
```

### 🔍 Monitoring and Maintenance

| Feature | Description |
|---------|-------------|
| 🤖 **Dependabot** | Automatically updates dependencies every Monday |
| 🔒 **Security Scans** | Run daily for vulnerability detection |
| 💚 **Health Checks** | Built-in health check in container |
| 📝 **Logging** | Request logging for monitoring |

### ➕ Adding Features

1. Edit [`main.go`](main.go) to add new logic
2. Add tests if needed
3. Update documentation in [`README.md`](README.md)
4. Create Pull Request with clear description

---

<div align="center">

## 📄 License

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

**Made with ❤️ using Go**

[⬆ Back to Top](#-http-proxy-server-with-go)

</div>

</div>
