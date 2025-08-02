# HTTP Proxy Server dengan Go

HTTP Proxy Server yang ditulis dalam bahasa Go dengan autentikasi Basic Auth dan siap untuk deployment menggunakan Docker.

## Fitur

- ✅ HTTP/HTTPS Proxy Server
- ✅ Autentikasi Basic Auth (username/password)
- ✅ Support untuk HTTP CONNECT method (HTTPS tunneling)
- ✅ Konfigurasi melalui environment variables
- ✅ Docker support dengan multi-stage build
- ✅ Health check
- ✅ Logging request
- ✅ Non-root user di container

## Cara Menggunakan

### 1. Development (Local)

```bash
# Clone atau download project
cd go-proxy-server

# Install dependencies
go mod tidy

# Set environment variables (opsional)
export PROXY_USERNAME=admin
export PROXY_PASSWORD=mypassword
export PROXY_PORT=8080

# Run aplikasi
go run main.go
```

### 2. Build Binary

```bash
# Build binary
go build -o proxy-server

# Run binary
./proxy-server
```

### 3. Docker Deployment

#### Menggunakan Docker Build

```bash
# Build Docker image
docker build -t go-proxy-server .

# Run container dengan default credentials
docker run -d \
  --name proxy-server \
  -p 8080:8080 \
  go-proxy-server

# Run container dengan custom credentials
docker run -d \
  --name proxy-server \
  -p 8080:8080 \
  -e PROXY_USERNAME=myuser \
  -e PROXY_PASSWORD=mypassword \
  go-proxy-server
```

#### Menggunakan Docker Compose

```bash
# Copy environment file
cp .env.example .env

# Edit .env file dengan credentials yang diinginkan
nano .env

# Start dengan docker-compose
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

## Konfigurasi

Aplikasi menggunakan environment variables untuk konfigurasi:

| Variable | Default | Deskripsi |
|----------|---------|-----------|
| `PROXY_USERNAME` | `admin` | Username untuk autentikasi proxy |
| `PROXY_PASSWORD` | `password123` | Password untuk autentikasi proxy |
| `PROXY_PORT` | `8080` | Port server proxy |

## Cara Menggunakan Proxy

### 1. Konfigurasi Proxy di Browser/Aplikasi

- **Proxy Type**: HTTP
- **Server**: `localhost` (atau IP server Anda)
- **Port**: `8080` (atau port yang Anda konfigurasi)
- **Username**: sesuai `PROXY_USERNAME`
- **Password**: sesuai `PROXY_PASSWORD`

### 2. Test dengan curl

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

### 3. Test dengan wget

```bash
# Set proxy environment
export http_proxy=http://admin:mypassword@localhost:8080
export https_proxy=http://admin:mypassword@localhost:8080

# Test request
wget http://httpbin.org/ip
```

## Monitoring

### Health Check

Aplikasi memiliki built-in health check yang bisa digunakan untuk monitoring:

```bash
# Check health (hanya untuk monitoring, bukan untuk proxy)
curl http://localhost:8080
```

### Logs

Aplikasi akan menampilkan log untuk setiap request:

```
2024/01/01 12:00:00 127.0.0.1:12345 GET http://example.com
2024/01/01 12:00:01 127.0.0.1:12346 CONNECT example.com:443
```

## Keamanan

- ✅ Autentikasi Basic Auth required untuk semua request
- ✅ Container berjalan sebagai non-root user
- ✅ Minimal Alpine Linux image
- ✅ No unnecessary packages installed

## Troubleshooting

### Error "Proxy Authentication Required"

Pastikan username dan password yang digunakan sesuai dengan konfigurasi server.

### Connection Timeout

Aplikasi memiliki timeout 30 detik untuk setiap request. Untuk request yang lebih lama, sesuaikan timeout di kode.

### Port Already in Use

Ganti port di environment variable `PROXY_PORT` atau gunakan port lain saat menjalankan container.

## CI/CD dan Versioning

### GitHub Actions Workflows

Project ini dilengkapi dengan GitHub Actions workflows untuk otomatisasi:

#### 1. CI/CD Pipeline (`ci-cd.yml`)
- **Trigger**: Push ke `main`/`develop`, Pull Request, atau Tag
- **Proses**:
  - ✅ Test dan build aplikasi
  - ✅ Security scan dengan staticcheck
  - ✅ Build multi-platform Docker images
  - ✅ Push ke GitHub Container Registry
  - ✅ Vulnerability scan dengan Trivy

#### 2. Release Pipeline (`release.yml`)
- **Trigger**: Push tag dengan format `v*` (contoh: `v1.0.0`)
- **Proses**:
  - ✅ Build binary untuk multiple platforms (Linux, macOS, Windows)
  - ✅ Create checksums
  - ✅ Build dan push Docker images
  - ✅ Create GitHub release dengan binaries
  - ✅ Generate release notes otomatis

#### 3. Security Scan (`security.yml`)
- **Trigger**: Schedule harian atau manual
- **Proses**:
  - ✅ Dependency scanning dengan Gosec
  - ✅ Container vulnerability scanning
  - ✅ Upload hasil ke GitHub Security tab

#### 4. Docker Hub Release (`dockerhub.yml`)
- **Trigger**: Push tag `v*` atau manual
- **Proses**:
  - ✅ Build dan push ke Docker Hub
  - ✅ Update description otomatis

### Setup GitHub Repository

1. **Secrets yang dibutuhkan**:
```bash
# Untuk Docker Hub (opsional)
DOCKERHUB_USERNAME=your-dockerhub-username
DOCKERHUB_TOKEN=your-dockerhub-access-token
```

2. **Branch Protection Rules** (Recommended):
   - Require pull request reviews
   - Require status checks (CI tests)
   - Require branches to be up to date

### Versioning dan Release

#### Menggunakan Script Release
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
# Create dan push tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# GitHub Actions akan otomatis:
# - Build aplikasi untuk multiple platforms
# - Create Docker images
# - Create GitHub release
# - Generate release notes
```

### Registry Images

#### GitHub Container Registry
```bash
# Pull latest
docker pull ghcr.io/your-username/go-proxy-server:latest

# Pull specific version
docker pull ghcr.io/your-username/go-proxy-server:v1.0.0
```

#### Docker Hub (opsional)
```bash
# Pull latest
docker pull your-username/go-proxy-server:latest

# Pull specific version
docker pull your-username/go-proxy-server:v1.0.0
```

## Development

### Quick Setup
```bash
# Setup development environment
./scripts/setup-dev.sh

# Install git hooks, dependencies, dan tools
```

### Struktur Project

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

### Development Workflow

1. **Setup environment**:
```bash
./scripts/setup-dev.sh
```

2. **Create feature branch**:
```bash
git checkout -b feature/your-feature
```

3. **Development**:
```bash
# Run locally
go run main.go

# Or with Docker
docker-compose up --build
```

4. **Testing**:
```bash
# Run tests
go test ./...

# Run linting
staticcheck ./...
golangci-lint run
```

5. **Commit dan Push**:
```bash
git add .
git commit -m "feat: add your feature"
git push origin feature/your-feature
```

6. **Create Pull Request**:
   - GitHub akan otomatis run CI checks
   - Review dan merge ke main

7. **Release**:
```bash
# Setelah merge ke main
git checkout main
git pull origin main
./scripts/release.sh patch  # atau minor/major
```

### Monitoring dan Maintenance

- **Dependabot**: Update dependencies otomatis setiap Senin
- **Security Scans**: Berjalan harian untuk vulnerability detection
- **Health Checks**: Built-in health check di container
- **Logging**: Request logging untuk monitoring

### Menambah Fitur

1. Edit `main.go` untuk menambah logika baru
2. Tambah tests jika diperlukan
3. Update dokumentasi di README.md
4. Create Pull Request dengan deskripsi yang jelas

## License

MIT License
