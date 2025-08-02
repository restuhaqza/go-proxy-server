#!/bin/bash

# Version release script
# Usage: ./scripts/release.sh [major|minor|patch]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default version bump type
BUMP_TYPE=${1:-patch}

echo -e "${GREEN}=== Go Proxy Server Release Script ===${NC}"

# Check if we're on main branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo -e "${RED}Error: You must be on the main branch to create a release${NC}"
    echo "Current branch: $CURRENT_BRANCH"
    exit 1
fi

# Check if working directory is clean
if [ -n "$(git status --porcelain)" ]; then
    echo -e "${RED}Error: Working directory is not clean${NC}"
    echo "Please commit or stash your changes before creating a release"
    git status --short
    exit 1
fi

# Pull latest changes
echo -e "${YELLOW}Pulling latest changes...${NC}"
git pull origin main

# Get current version
CURRENT_VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
echo "Current version: $CURRENT_VERSION"

# Remove 'v' prefix for calculation
VERSION_NUMBER=${CURRENT_VERSION#v}

# Split version into parts
IFS='.' read -ra VERSION_PARTS <<< "$VERSION_NUMBER"
MAJOR=${VERSION_PARTS[0]:-0}
MINOR=${VERSION_PARTS[1]:-0}
PATCH=${VERSION_PARTS[2]:-0}

# Calculate new version
case $BUMP_TYPE in
    major)
        NEW_MAJOR=$((MAJOR + 1))
        NEW_MINOR=0
        NEW_PATCH=0
        ;;
    minor)
        NEW_MAJOR=$MAJOR
        NEW_MINOR=$((MINOR + 1))
        NEW_PATCH=0
        ;;
    patch)
        NEW_MAJOR=$MAJOR
        NEW_MINOR=$MINOR
        NEW_PATCH=$((PATCH + 1))
        ;;
    *)
        echo -e "${RED}Error: Invalid bump type. Use major, minor, or patch${NC}"
        exit 1
        ;;
esac

NEW_VERSION="v${NEW_MAJOR}.${NEW_MINOR}.${NEW_PATCH}"

echo -e "${YELLOW}Bumping version from $CURRENT_VERSION to $NEW_VERSION${NC}"

# Confirm release
read -p "Do you want to create release $NEW_VERSION? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Release cancelled"
    exit 0
fi

# Update version in files if needed
# You can add version updates to specific files here
# For example, updating version in main.go or other files

# Create and push tag
echo -e "${YELLOW}Creating and pushing tag $NEW_VERSION...${NC}"
git tag -a "$NEW_VERSION" -m "Release $NEW_VERSION"
git push origin "$NEW_VERSION"

echo -e "${GREEN}✅ Release $NEW_VERSION created successfully!${NC}"
echo -e "${GREEN}GitHub Actions will automatically:${NC}"
echo "  - Build and test the code"
echo "  - Create Docker images for multiple platforms"
echo "  - Push to GitHub Container Registry"
echo "  - Create GitHub release with binaries"
echo "  - Run security scans"
echo ""
echo -e "${YELLOW}You can monitor the release progress at:${NC}"
echo "https://github.com/$(git config --get remote.origin.url | sed 's/.*github.com[:/]\([^/]*\/[^/]*\).*/\1/' | sed 's/\.git$//')/actions"
