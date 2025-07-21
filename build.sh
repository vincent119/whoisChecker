#!/bin/bash

# WhoisChecker Cross-Platform Build Script

set -e

# Project information
APP_NAME="whoisChecker"
VERSION="v1.0.0"
BUILD_DIR="build"
SOURCE_FILE="./cmd/main.go"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔨 Building WhoisChecker for multiple platforms...${NC}"

# Create build directory
if [ -d "$BUILD_DIR" ]; then
    echo -e "${YELLOW}📁 Cleaning existing build directory...${NC}"
    rm -rf "$BUILD_DIR"
fi
mkdir -p "$BUILD_DIR"

# Get current git commit hash (if available)
if git rev-parse --git-dir > /dev/null 2>&1; then
    COMMIT_HASH=$(git rev-parse --short HEAD)
    BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
    LDFLAGS="-X main.Version=${VERSION} -X main.CommitHash=${COMMIT_HASH} -X main.BuildTime=${BUILD_TIME}"
else
    BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
    LDFLAGS="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"
fi

echo -e "${YELLOW}📦 Version: ${VERSION}${NC}"
echo -e "${YELLOW}🕒 Build Time: ${BUILD_TIME}${NC}"
if [ ! -z "$COMMIT_HASH" ]; then
    echo -e "${YELLOW}🔗 Commit: ${COMMIT_HASH}${NC}"
fi

# Build function
build_platform() {
    local os=$1
    local arch=$2
    local ext=$3
    
    local output_name="${APP_NAME}_${os}_${arch}${ext}"
    local output_path="${BUILD_DIR}/${output_name}"
    
    echo -e "${BLUE}🏗️  Building for ${os}/${arch}...${NC}"
    
    GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build \
        -ldflags="$LDFLAGS" \
        -o "$output_path" \
        "$SOURCE_FILE"
    
    if [ $? -eq 0 ]; then
        local size=$(du -h "$output_path" | cut -f1)
        echo -e "${GREEN}✅ Successfully built: ${output_name} (${size})${NC}"
    else
        echo -e "${RED}❌ Failed to build for ${os}/${arch}${NC}"
        return 1
    fi
}

# Build for different platforms
echo -e "\n${BLUE}🚀 Starting cross-platform compilation...${NC}\n"

# macOS
build_platform "darwin" "amd64" ""      # Intel Mac
build_platform "darwin" "arm64" ""      # Apple Silicon Mac

# Windows
build_platform "windows" "amd64" ".exe" # Windows 64-bit
build_platform "windows" "386" ".exe"   # Windows 32-bit

# Linux
build_platform "linux" "amd64" ""       # Linux 64-bit
build_platform "linux" "386" ""         # Linux 32-bit
build_platform "linux" "arm64" ""       # Linux ARM64 (Ubuntu on ARM)
build_platform "linux" "arm" ""         # Linux ARM (Raspberry Pi)

echo -e "\n${GREEN}🎉 All builds completed successfully!${NC}"

# List all built files
echo -e "\n${BLUE}📋 Built files:${NC}"
ls -la "$BUILD_DIR"

# Create checksums
echo -e "\n${BLUE}🔐 Generating checksums...${NC}"
cd "$BUILD_DIR"
if command -v sha256sum >/dev/null 2>&1; then
    sha256sum * > checksums.txt
elif command -v shasum >/dev/null 2>&1; then
    shasum -a 256 * > checksums.txt
else
    echo -e "${YELLOW}⚠️  No checksum utility found, skipping checksums${NC}"
fi
cd ..

echo -e "\n${GREEN}✨ Build process completed!${NC}"
echo -e "${BLUE}📁 All binaries are available in the '${BUILD_DIR}' directory${NC}"
