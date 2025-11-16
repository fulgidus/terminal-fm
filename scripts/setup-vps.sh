#!/bin/bash
set -e

# Terminal-Radio VPS Setup Script
# Automatic setup for Ubuntu 22.04 VPS
# Usage: ./setup-vps.sh [--skip-ssh-config]

echo "=================================================="
echo "  Terminal-Radio VPS Automatic Setup"
echo "=================================================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="terminal-radio"
APP_DIR="/opt/terminal-radio"
APP_PORT=22
ADMIN_SSH_PORT=2222
GITHUB_REPO="https://github.com/fulgidus/terminal-fm.git"
GO_VERSION="1.21.5"

# Check if running as root
if [ "$EUID" -eq 0 ]; then 
    echo -e "${RED}❌ Do not run this script as root!${NC}"
    echo "Run as ubuntu user: ./setup-vps.sh"
    exit 1
fi

# Parse arguments
SKIP_SSH_CONFIG=false
if [ "$1" == "--skip-ssh-config" ]; then
    SKIP_SSH_CONFIG=true
    echo -e "${YELLOW}⚠️  Skipping SSH configuration (port will remain on 22)${NC}"
fi

echo -e "${GREEN}[1/10]${NC} Updating system packages..."
sudo apt update && sudo apt upgrade -y

echo ""
echo -e "${GREEN}[2/10]${NC} Installing dependencies (Go, FFmpeg, Git, Build tools)..."
sudo apt install -y build-essential git curl wget ffmpeg sqlite3 ufw

echo ""
echo -e "${GREEN}[3/10]${NC} Installing Go ${GO_VERSION}..."
if ! command -v go &> /dev/null || ! go version | grep -q "go${GO_VERSION}"; then
    cd /tmp
    wget -q https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
    rm go${GO_VERSION}.linux-amd64.tar.gz
    
    # Add Go to PATH
    if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
    fi
    export PATH=$PATH:/usr/local/go/bin
    
    echo -e "  ✓ Go $(go version) installed"
else
    echo -e "  ✓ Go already installed: $(go version)"
fi

echo ""
echo -e "${GREEN}[4/10]${NC} Creating dedicated user and directories..."
# Create user if doesn't exist
if ! id -u "$APP_NAME" &>/dev/null; then
    sudo useradd -r -s /bin/false -m -d "$APP_DIR" "$APP_NAME"
    echo -e "  ✓ User $APP_NAME created"
else
    echo -e "  ✓ User $APP_NAME already exists"
fi

# Create directories
sudo mkdir -p "$APP_DIR"
sudo mkdir -p "$APP_DIR/backups"
sudo chown -R ubuntu:ubuntu "$APP_DIR"

echo ""
echo -e "${GREEN}[5/10]${NC} Cloning repository..."
if [ ! -d "$APP_DIR/.git" ]; then
    git clone "$GITHUB_REPO" "$APP_DIR"
    echo -e "  ✓ Repository cloned"
else
    echo -e "  ✓ Repository already exists, pulling latest..."
    cd "$APP_DIR" && git pull origin main
fi

echo ""
echo -e "${GREEN}[6/10]${NC} Building application..."
cd "$APP_DIR"
go mod download
go build -ldflags="-s -w" -o terminal-radio ./cmd/server
sudo chown $APP_NAME:$APP_NAME "$APP_DIR/terminal-radio"
sudo chmod +x "$APP_DIR/terminal-radio"
echo -e "  ✓ Binary built: $(ls -lh terminal-radio | awk '{print $5}')"

echo ""
echo -e "${GREEN}[7/10]${NC} Creating systemd service..."
sudo tee /etc/systemd/system/$APP_NAME.service > /dev/null <<EOF
[Unit]
Description=Terminal-Radio - TUI Radio Player
After=network.target
Wants=network.target

[Service]
Type=simple
User=$APP_NAME
Group=$APP_NAME
WorkingDirectory=$APP_DIR
ExecStart=$APP_DIR/terminal-radio --port $APP_PORT --host 0.0.0.0
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=$APP_NAME

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=$APP_DIR
AmbientCapabilities=CAP_NET_BIND_SERVICE

# Resource limits
LimitNOFILE=10000
MemoryMax=512M
CPUQuota=50%

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable $APP_NAME
echo -e "  ✓ Systemd service created and enabled"

echo ""
echo -e "${GREEN}[8/10]${NC} Creating deployment script for CI/CD..."
sudo tee /usr/local/bin/${APP_NAME}-deploy.sh > /dev/null <<'EOFSCRIPT'
#!/bin/bash
set -e

APP_NAME="terminal-radio"
APP_DIR="/opt/terminal-radio"

echo "🚀 Starting Terminal-Radio deployment..."

# Stop service
echo "⏸️  Stopping service..."
sudo systemctl stop $APP_NAME || true

# Backup current binary
if [ -f "$APP_DIR/terminal-radio" ]; then
    echo "💾 Backing up current version..."
    sudo cp $APP_DIR/terminal-radio $APP_DIR/backups/terminal-radio.$(date +%Y%m%d-%H%M%S)
    # Keep only last 5 backups
    cd $APP_DIR/backups && ls -t terminal-radio.* 2>/dev/null | tail -n +6 | xargs -r rm
fi

# Pull latest code
echo "📥 Pulling latest code..."
cd $APP_DIR
sudo -u ubuntu git fetch origin
sudo -u ubuntu git reset --hard origin/main

# Build
echo "🔨 Building application..."
sudo -u ubuntu go build -ldflags="-s -w" -o terminal-radio ./cmd/server

# Set permissions
sudo chown $APP_NAME:$APP_NAME $APP_DIR/terminal-radio
sudo chmod +x $APP_DIR/terminal-radio

# Start service
echo "▶️  Starting service..."
sudo systemctl start $APP_NAME

# Wait and check status
sleep 3
if sudo systemctl is-active --quiet $APP_NAME; then
    echo "✅ Service is running!"
    sudo systemctl status $APP_NAME --no-pager -l
    echo ""
    echo "🎉 Deployment completed successfully!"
    exit 0
else
    echo "❌ Service failed to start!"
    sudo journalctl -u $APP_NAME -n 50 --no-pager
    exit 1
fi
EOFSCRIPT

sudo chmod +x /usr/local/bin/${APP_NAME}-deploy.sh

# Configure sudo permissions for deploy script
if ! sudo grep -q "${APP_NAME}-deploy" /etc/sudoers.d/${APP_NAME} 2>/dev/null; then
    echo "ubuntu ALL=(ALL) NOPASSWD: /usr/local/bin/${APP_NAME}-deploy.sh, /usr/bin/systemctl start ${APP_NAME}, /usr/bin/systemctl stop ${APP_NAME}, /usr/bin/systemctl restart ${APP_NAME}, /usr/bin/systemctl status ${APP_NAME}, /usr/bin/systemctl is-active ${APP_NAME}, /usr/bin/journalctl" | sudo tee /etc/sudoers.d/${APP_NAME} > /dev/null
    sudo chmod 440 /etc/sudoers.d/${APP_NAME}
fi
echo -e "  ✓ Deploy script created: /usr/local/bin/${APP_NAME}-deploy.sh"

echo ""
echo -e "${GREEN}[9/10]${NC} Configuring firewall..."
sudo ufw --force enable
sudo ufw allow $ADMIN_SSH_PORT/tcp comment 'Admin SSH'
sudo ufw allow $APP_PORT/tcp comment 'Terminal-Radio SSH'
sudo ufw status
echo -e "  ✓ Firewall configured"

if [ "$SKIP_SSH_CONFIG" = false ]; then
    echo ""
    echo -e "${GREEN}[10/10]${NC} Configuring SSH (moving to port $ADMIN_SSH_PORT)..."
    echo -e "${YELLOW}⚠️  IMPORTANT: Testing new SSH port before applying...${NC}"
    
    # Backup SSH config
    sudo cp /etc/ssh/sshd_config /etc/ssh/sshd_config.backup.$(date +%Y%m%d-%H%M%S)
    
    # Check if Port directive already exists
    if grep -q "^Port " /etc/ssh/sshd_config; then
        sudo sed -i "s/^Port .*/Port $ADMIN_SSH_PORT/" /etc/ssh/sshd_config
    else
        echo "Port $ADMIN_SSH_PORT" | sudo tee -a /etc/ssh/sshd_config > /dev/null
    fi
    
    # Test SSH config
    if sudo sshd -t; then
        echo -e "  ✓ SSH config is valid"
        echo -e "${YELLOW}"
        echo "  ⚠️  SSH will be moved to port $ADMIN_SSH_PORT"
        echo "  ⚠️  From now on, use: ssh -p $ADMIN_SSH_PORT ubuntu@<IP>"
        echo -e "${NC}"
        sudo systemctl restart sshd
        echo -e "  ✓ SSH restarted on port $ADMIN_SSH_PORT"
    else
        echo -e "${RED}  ❌ SSH config test failed! Reverting...${NC}"
        sudo mv /etc/ssh/sshd_config.backup.$(date +%Y%m%d-%H%M%S) /etc/ssh/sshd_config
        exit 1
    fi
else
    echo ""
    echo -e "${YELLOW}[10/10]${NC} SSH configuration skipped (--skip-ssh-config)"
fi

echo ""
echo -e "${GREEN}[FINAL]${NC} Starting Terminal-Radio service..."
sudo systemctl start $APP_NAME
sleep 2

if sudo systemctl is-active --quiet $APP_NAME; then
    echo -e "${GREEN}✅ Terminal-Radio is running!${NC}"
    sudo systemctl status $APP_NAME --no-pager
else
    echo -e "${RED}❌ Service failed to start!${NC}"
    echo "Check logs with: sudo journalctl -u $APP_NAME -n 50"
    exit 1
fi

echo ""
echo "=================================================="
echo -e "${GREEN}  ✨ Setup Completed Successfully! ✨${NC}"
echo "=================================================="
echo ""
echo "Next steps:"
echo ""
echo "1. 🔑 Setup GitHub Actions SSH key:"
echo "   - On your LOCAL machine:"
echo "     ssh-keygen -t ed25519 -C 'github-actions' -f ~/.ssh/terminal-radio-deploy"
echo "     cat ~/.ssh/terminal-radio-deploy.pub"
echo "   - Copy the public key"
echo "   - On THIS server:"
echo "     nano ~/.ssh/authorized_keys"
echo "     (paste the public key and save)"
echo ""
echo "2. 📋 Add GitHub Secret:"
echo "   - Go to: https://github.com/fulgidus/terminal-fm/settings/secrets/actions"
echo "   - Create secret: SSH_PRIVATE_KEY"
echo "   - Paste content of: ~/.ssh/terminal-radio-deploy (PRIVATE key)"
echo ""
echo "3. ✅ Test connection:"
if [ "$SKIP_SSH_CONFIG" = false ]; then
    echo "   ssh -p $ADMIN_SSH_PORT ubuntu@\$(curl -s ifconfig.me)"
else
    echo "   ssh ubuntu@\$(curl -s ifconfig.me)"
fi
echo ""
echo "4. 🎵 Test Terminal-Radio:"
echo "   ssh \$(curl -s ifconfig.me)"
echo "   (or wait for DNS: ssh terminal-radio.com)"
echo ""
echo "5. 🚀 Push to main branch → Auto-deploy!"
echo ""
