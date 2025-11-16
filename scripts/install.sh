#!/bin/bash
set -e

# Terminal-Radio Client Installer
# Installs the radio wrapper command

echo "╔════════════════════════════════════════════════╗"
echo "║     Terminal-Radio Client Installer           ║"
echo "╚════════════════════════════════════════════════╝"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
GITHUB_REPO="fulgidus/terminal-fm"
INSTALL_DIR="$HOME/.local/bin"
BINARY_NAME="radio"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    armv7l)
        ARCH="arm"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo -e "${BLUE}Detected: ${NC}$OS/$ARCH"
echo ""

# Check for required dependencies
echo -e "${BLUE}Checking dependencies...${NC}"

# Check if ssh is available
if ! command -v ssh &> /dev/null; then
    echo -e "${RED}✗ ssh not found${NC}"
    echo "  Please install OpenSSH client"
    exit 1
fi
echo -e "${GREEN}✓ ssh${NC}"

# Check for audio player
HAS_PLAYER=false
if command -v mpv &> /dev/null; then
    echo -e "${GREEN}✓ mpv${NC}"
    HAS_PLAYER=true
elif command -v ffplay &> /dev/null; then
    echo -e "${GREEN}✓ ffplay${NC}"
    HAS_PLAYER=true
elif command -v vlc &> /dev/null; then
    echo -e "${GREEN}✓ vlc${NC}"
    HAS_PLAYER=true
fi

if [ "$HAS_PLAYER" = false ]; then
    echo -e "${YELLOW}⚠  No audio player found${NC}"
    echo ""
    echo "Terminal-Radio requires one of these players:"
    echo "  • mpv (recommended): https://mpv.io/"
    echo "  • ffplay (part of ffmpeg)"
    echo "  • vlc"
    echo ""
    echo "Install instructions:"
    echo "  macOS:   brew install mpv"
    echo "  Ubuntu:  sudo apt install mpv"
    echo "  Fedora:  sudo dnf install mpv"
    echo "  Arch:    sudo pacman -S mpv"
    echo ""
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo ""

# Create install directory
mkdir -p "$INSTALL_DIR"

# Download or build client
echo -e "${BLUE}Installing terminal-radio client...${NC}"

# Check if Go is available for local build
if command -v go &> /dev/null; then
    echo -e "${BLUE}Building from source...${NC}"
    
    # Create temp directory
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"
    
    # Clone repo
    git clone --depth 1 "https://github.com/$GITHUB_REPO.git" . 2>/dev/null || {
        echo -e "${RED}Failed to clone repository${NC}"
        exit 1
    }
    
    # Build client
    go build -o "$INSTALL_DIR/$BINARY_NAME" ./cmd/client || {
        echo -e "${RED}Failed to build client${NC}"
        exit 1
    }
    
    # Cleanup
    cd -
    rm -rf "$TMP_DIR"
    
    echo -e "${GREEN}✓ Built from source${NC}"
else
    # Try to download pre-built binary (future enhancement)
    echo -e "${YELLOW}Go not found, attempting to download pre-built binary...${NC}"
    
    # For now, we'll just download the source and suggest installing Go
    echo -e "${RED}Pre-built binaries not yet available.${NC}"
    echo ""
    echo "Please install Go to build the client:"
    echo "  https://go.dev/doc/install"
    echo ""
    echo "Or install directly with:"
    echo "  go install github.com/$GITHUB_REPO/cmd/client@latest"
    echo "  mv \$(go env GOPATH)/bin/client ~/.local/bin/radio"
    exit 1
fi

# Make executable
chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Check if install dir is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    echo -e "${YELLOW}⚠  $INSTALL_DIR is not in your PATH${NC}"
    echo ""
    echo "Add this line to your ~/.bashrc or ~/.zshrc:"
    echo ""
    echo -e "${BLUE}  export PATH=\"\$PATH:$INSTALL_DIR\"${NC}"
    echo ""
    echo "Then run: source ~/.bashrc (or ~/.zshrc)"
    echo ""
fi

echo ""
echo -e "${GREEN}╔════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║         Installation Complete! 🎉             ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════╝${NC}"
echo ""
echo "Run terminal-radio with:"
echo -e "  ${BLUE}radio${NC}"
echo ""
echo "Or connect directly:"
echo -e "  ${BLUE}ssh terminal-radio.com${NC}"
echo ""
echo "Note: Using 'radio' command enables local audio playback"
echo "      Direct SSH connection will not play audio"
echo ""
