#!/usr/bin/env bash

set -euo pipefail

REPO="ByteTheCookies/CookieFarm"
INSTALL_DIR="/usr/local/bin"
TMP_DIR="$(mktemp -d)"
ASSET_NAME="cookieclient"
FINAL_NAME="cookieclient"

REQUIRED_CMDS=("curl" "jq")

echo "🔍 Checking for required tools..."
for cmd in "${REQUIRED_CMDS[@]}"; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "❌ Error: '$cmd' is not installed. Please install it first."
    exit 1
  fi
done

# Fetch latest version info from GitHub
echo "🌐 Checking latest release..."
LATEST_VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | jq -r '.tag_name')
ASSET_URL=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | jq -r '.assets[].browser_download_url' | grep "$ASSET_NAME" || true)

if [[ -z "$ASSET_URL" || -z "$LATEST_VERSION" ]]; then
  echo "❌ Error: Unable to fetch latest release or asset."
  exit 1
fi

echo "⬇️ Downloading from: $ASSET_URL"
curl -L "$ASSET_URL" -o "$TMP_DIR/$ASSET_NAME"

echo "🚚 Installing to $INSTALL_DIR..."
sudo mv "$TMP_DIR/$ASSET_NAME" "$INSTALL_DIR/$FINAL_NAME"
sudo chmod +x "$INSTALL_DIR/$FINAL_NAME"

echo "🔧 Setting up configuration..."
CONFIG_DIR="$HOME/.config/cookiefarm"
mkdir -p "$CONFIG_DIR"
echo "# Configuration directory for CookieFarm" > "$CONFIG_DIR/.readme"

echo "📄 Creating default configuration file..."
cat <<EOF > "$CONFIG_DIR/config.yml"
address: "localhost"
port: 8080
https: false
nickname: "cookieguest"
EOF

echo "✅ Installation complete! Run '$FINAL_NAME --help' to get started. Enjoy farming 🍪"

rm -rf "$TMP_DIR"
