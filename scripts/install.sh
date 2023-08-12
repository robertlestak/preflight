#!/bin/bash
set -e

# Get the installation directory from environment variable or use default
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Get the caller's OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
fi

# GitHub repository and release URL
REPO_OWNER="robertlestak"
REPO_NAME="preflight"
API_URL="https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest"

# Fetch the latest release information
release_info=$(curl -s "$API_URL")
asset_name="preflight_${OS}_${ARCH}"

# Find the asset with the desired name
asset_url=$(echo "$release_info" | jq -r ".assets[] | select(.name == \"$asset_name\") | .browser_download_url")

if [ -z "$asset_url" ]; then
    echo "No matching binary release found for your OS and architecture."
    exit 1
fi

# Download the asset
echo "Downloading $asset_name..."
curl -LJO "$asset_url"

# Get the downloaded filename
downloaded_file=$(basename "$asset_url")

# Make the downloaded binary executable
chmod +x "$downloaded_file"
# Move the downloaded binary to the installation directory
renamed_file="preflight"
mv "$downloaded_file" "$INSTALL_DIR/$renamed_file"

echo "Binary installed in $INSTALL_DIR."
echo "Installation complete."
