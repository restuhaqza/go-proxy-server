# Setup GitHub Repository

Panduan untuk setup repository GitHub dengan CI/CD yang telah dibuat.

## 1. Push ke GitHub

```bash
# Initialize git repository (jika belum)
git init
git add .
git commit -m "Initial commit: HTTP Proxy Server with CI/CD"

# Add remote repository
git remote add origin https://github.com/YOUR_USERNAME/go-proxy-server.git

# Push ke GitHub
git branch -M main
git push -u origin main
```

## 2. Setup GitHub Container Registry

GitHub Container Registry (GHCR) akan digunakan secara otomatis tanpa setup tambahan karena menggunakan `GITHUB_TOKEN` yang tersedia secara default.

Images akan tersedia di: `ghcr.io/YOUR_USERNAME/go-proxy-server`

## 3. Setup Docker Hub (Opsional)

Jika ingin push ke Docker Hub juga:

### 3.1 Buat Docker Hub Access Token
1. Login ke [Docker Hub](https://hub.docker.com)
2. Go to Account Settings > Security
3. Click "New Access Token"
4. Beri nama "GitHub Actions" dan copy token

### 3.2 Setup Repository Secrets
1. Di GitHub repository, go to Settings > Secrets and variables > Actions
2. Click "New repository secret"
3. Tambahkan secrets berikut:

```
DOCKERHUB_USERNAME = your-dockerhub-username
DOCKERHUB_TOKEN = your-dockerhub-access-token
```

## 4. Branch Protection Rules (Recommended)

1. Go to Settings > Branches
2. Click "Add rule"
3. Branch name pattern: `main`
4. Enable:
   - ✅ Require a pull request before merging
   - ✅ Require status checks to pass before merging
   - ✅ Require branches to be up to date before merging
   - ✅ Status checks: `test`, `build-and-push`

## 5. Enable GitHub Pages (Opsional)

Untuk dokumentasi otomatis:
1. Go to Settings > Pages
2. Source: Deploy from a branch
3. Branch: `main`, folder: `/ (root)`

## 6. Setup Repository Labels

Tambahkan labels untuk issue management:
```
bug - Something isn't working
enhancement - New feature or request
documentation - Improvements or additions to documentation
good first issue - Good for newcomers
help wanted - Extra attention is needed
security - Security related issues
```

## 7. Test CI/CD Pipeline

### 7.1 Test dengan Push
```bash
echo "# Test" >> README.md
git add README.md
git commit -m "test: trigger CI/CD"
git push origin main
```

### 7.2 Test Release
```bash
# Gunakan script release
./scripts/release.sh patch

# Atau manual
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

## 8. Monitor Workflows

1. Go to Actions tab di GitHub repository
2. Monitor workflow runs:
   - ✅ CI/CD Pipeline: Test dan build
   - ✅ Release: Create release dengan binaries
   - ✅ Security Scan: Daily security checks

## 9. Use Released Images

### GitHub Container Registry
```bash
# Pull latest
docker pull ghcr.io/YOUR_USERNAME/go-proxy-server:latest

# Run container
docker run -d \
  --name proxy-server \
  -p 8080:8080 \
  -e PROXY_USERNAME=admin \
  -e PROXY_PASSWORD=mypassword \
  ghcr.io/YOUR_USERNAME/go-proxy-server:latest
```

### Docker Hub (jika disetup)
```bash
# Pull latest
docker pull YOUR_USERNAME/go-proxy-server:latest

# Run container
docker run -d \
  --name proxy-server \
  -p 8080:8080 \
  -e PROXY_USERNAME=admin \
  -e PROXY_PASSWORD=mypassword \
  YOUR_USERNAME/go-proxy-server:latest
```

## 10. Development Workflow

### Setup Development Environment
```bash
./scripts/setup-dev.sh
```

### Create Feature
```bash
git checkout -b feature/new-feature
# Make changes
git add .
git commit -m "feat: add new feature"
git push origin feature/new-feature
# Create Pull Request di GitHub
```

### Release New Version
```bash
git checkout main
git pull origin main
./scripts/release.sh patch  # atau minor/major
```

## Troubleshooting

### Permission Denied untuk GITHUB_TOKEN
- GITHUB_TOKEN memiliki permission otomatis untuk push ke GHCR
- Pastikan repository settings > Actions > General > Workflow permissions diset ke "Read and write permissions"

### Docker Hub Push Failed
- Verify DOCKERHUB_USERNAME dan DOCKERHUB_TOKEN secrets
- Pastikan Docker Hub repository sudah dibuat atau set ke public

### Workflow Tidak Berjalan
- Cek file .github/workflows/ sudah ter-push ke repository
- Pastikan YAML syntax valid
- Cek Actions tab untuk error messages

### Security Scan Failures
- Review security alerts di Security tab
- Update dependencies jika ada vulnerability
- Whitelist false positives jika perlu

## Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Docker Hub](https://hub.docker.com)
- [Semantic Versioning](https://semver.org)
