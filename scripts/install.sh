#!/bin/bash

set -e

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="elitecode"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Download URL
DOWNLOAD_URL="https://github.com/yourusername/elitecode/releases/latest/download/elitecode-${OS}-${ARCH}.tar.gz"

echo "Installing Elitecode CLI..."
echo "OS: $OS, Architecture: $ARCH"

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd $TMP_DIR

# Download and extract
echo "Downloading from $DOWNLOAD_URL..."
curl -L -o elitecode.tar.gz "$DOWNLOAD_URL"
tar -xzf elitecode.tar.gz

# Install binary
echo "Installing to $INSTALL_DIR..."
sudo mv "elitecode-${OS}-${ARCH}" "$INSTALL_DIR/$BINARY_NAME"
sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Cleanup
cd -
rm -rf $TMP_DIR

echo "Elitecode CLI installed successfully!"
echo "Run 'elitecode --help' to get started."