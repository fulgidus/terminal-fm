# Deployment Guide

This guide walks you through deploying Terminal.FM on an OVH VPS or any other server.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Server Setup](#server-setup)
- [OVH VPS Setup](#ovh-vps-setup)
- [Domain Configuration](#domain-configuration)
- [Building the Application](#building-the-application)
- [Systemd Service](#systemd-service)
- [Security Hardening](#security-hardening)
- [Monitoring](#monitoring)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### Hardware Requirements

**Minimum (10-20 concurrent users)**:
- 1 vCPU
- 1GB RAM
- 10GB SSD
- 1 Mbps bandwidth

**Recommended (50-100 concurrent users)**:
- 2 vCPU
- 4GB RAM
- 40GB SSD
- 10 Mbps bandwidth

### Software Requirements

- Ubuntu 22.04 LTS (or compatible Linux distribution)
- Go 1.21+ (for building)
- systemd (for service management)
- OpenSSH server (for admin access on port 2222)

### Domain Requirements

- Domain name pointing to your server IP
- DNS A record configured

## Server Setup

### 1. Initial Server Configuration

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install required packages
sudo apt install -y build-essential git curl wget

# Install Go (if not already installed)
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify Go installation
go version
```

### 2. Create Dedicated User

```bash
# Create terminal-fm user
sudo useradd -r -s /bin/false -m -d /var/lib/terminal-fm terminal-fm

# Create necessary directories
sudo mkdir -p /var/lib/terminal-fm
sudo mkdir -p /var/log/terminal-fm
sudo chown -R terminal-fm:terminal-fm /var/lib/terminal-fm
sudo chown -R terminal-fm:terminal-fm /var/log/terminal-fm
```

### 3. Firewall Configuration

```bash
# Install UFW (if not already installed)
sudo apt install -y ufw

# Allow SSH on port 22 (Terminal.FM)
sudo ufw allow 22/tcp comment 'Terminal.FM SSH'

# Allow SSH on port 2222 (Admin access)
sudo ufw allow 2222/tcp comment 'Admin SSH'

# Enable firewall
sudo ufw enable

# Check status
sudo ufw status
```

## OVH VPS Setup

### 1. Order VPS

1. Go to [OVH VPS](https://www.ovhcloud.com/it/vps/)
2. Select VPS plan:
   - **Starter**: €3.50/month (for testing)
   - **Value**: €7/month (recommended for production)
3. Choose Ubuntu 22.04 LTS
4. Complete order

### 2. Initial Access

```bash
# SSH into your VPS (use credentials from OVH email)
ssh ubuntu@<YOUR_VPS_IP>

# Change root password
sudo passwd

# Update system
sudo apt update && sudo apt upgrade -y
```

### 3. Configure Admin SSH on Port 2222

```bash
# Edit SSH config
sudo nano /etc/ssh/sshd_config

# Add these lines at the end:
Port 2222
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes

# Restart SSH service
sudo systemctl restart sshd

# Test connection on new port (DO NOT CLOSE CURRENT SESSION)
# In a new terminal:
ssh -p 2222 ubuntu@<YOUR_VPS_IP>
```

## Domain Configuration

### 1. Register Domain

**Option A: .fm domain (expensive but thematic)**
- Cost: ~€80-100/year
- Registrars: Namecheap, Gandi, Hover

**Option B: Alternative domains (cheaper)**
- terminal-fm.com: ~€10/year
- terminal-fm.org: ~€10/year
- terminal-fm.dev: ~€12/year

### 2. Configure DNS

Add an A record pointing to your VPS IP:

```
Type: A
Name: @ (or terminal-fm if subdomain)
Value: <YOUR_VPS_IP>
TTL: 3600
```

Wait for DNS propagation (up to 48 hours, usually ~1 hour):

```bash
# Test DNS resolution
dig terminal.fm
# or
nslookup terminal.fm
```

## Building the Application

### 1. Clone Repository

```bash
# As ubuntu user
cd ~
git clone https://github.com/fulgidus/terminal-fm.git
cd terminal-fm
```

### 2. Build Binary

```bash
# Build for production
go build -ldflags="-s -w" -o terminal-fm ./cmd/server

# Verify binary
./terminal-fm --version
```

### 3. Install Binary

```bash
# Copy binary to system location
sudo cp terminal-fm /var/lib/terminal-fm/terminal-fm
sudo chown terminal-fm:terminal-fm /var/lib/terminal-fm/terminal-fm
sudo chmod +x /var/lib/terminal-fm/terminal-fm
```

## Systemd Service

### 1. Create Service File

```bash
sudo nano /etc/systemd/system/terminal-fm.service
```

Add this configuration:

```ini
[Unit]
Description=Terminal.FM - TUI Radio Player
After=network.target
Wants=network.target

[Service]
Type=simple
User=terminal-fm
Group=terminal-fm
WorkingDirectory=/var/lib/terminal-fm
ExecStart=/var/lib/terminal-fm/terminal-fm --port 22 --host 0.0.0.0
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=terminal-fm

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/terminal-fm
AmbientCapabilities=CAP_NET_BIND_SERVICE

# Resource limits
LimitNOFILE=10000
MemoryMax=512M
CPUQuota=50%

[Install]
WantedBy=multi-user.target
```

### 2. Enable and Start Service

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service (start on boot)
sudo systemctl enable terminal-fm

# Start service
sudo systemctl start terminal-fm

# Check status
sudo systemctl status terminal-fm

# View logs
sudo journalctl -u terminal-fm -f
```

### 3. Verify Deployment

```bash
# Test SSH connection
ssh terminal.fm
# or
ssh <YOUR_VPS_IP>
```

## Security Hardening

### 1. SSH Security

```bash
# Edit sshd_config (port 2222 - admin access)
sudo nano /etc/ssh/sshd_config

# Recommended settings:
Protocol 2
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes
MaxAuthTries 3
MaxSessions 10
ClientAliveInterval 300
ClientAliveCountMax 2
```

### 2. Fail2Ban Setup

```bash
# Install fail2ban
sudo apt install -y fail2ban

# Create jail for Terminal.FM
sudo nano /etc/fail2ban/jail.d/terminal-fm.conf
```

Add:

```ini
[terminal-fm]
enabled = true
port = 22
filter = sshd
logpath = /var/log/auth.log
maxretry = 10
bantime = 3600
findtime = 600
```

```bash
# Restart fail2ban
sudo systemctl restart fail2ban

# Check status
sudo fail2ban-client status terminal-fm
```

### 3. Automatic Security Updates

```bash
# Install unattended-upgrades
sudo apt install -y unattended-upgrades

# Enable automatic updates
sudo dpkg-reconfigure -plow unattended-upgrades
```

## Monitoring

### 1. Basic Monitoring

```bash
# View service status
sudo systemctl status terminal-fm

# View logs
sudo journalctl -u terminal-fm -n 100 --no-pager

# Follow logs in real-time
sudo journalctl -u terminal-fm -f

# Monitor resource usage
htop
```

### 2. Check Active Connections

```bash
# List SSH connections on port 22
sudo ss -tnp | grep ':22'

# Count active connections
sudo ss -tnp | grep ':22' | wc -l
```

### 3. Database Monitoring

```bash
# Check database size
du -h /var/lib/terminal-fm/data.db

# Backup database
sudo cp /var/lib/terminal-fm/data.db /var/lib/terminal-fm/backups/data-$(date +%Y%m%d).db
```

### 4. Log Rotation

```bash
# Create log rotation config
sudo nano /etc/logrotate.d/terminal-fm
```

Add:

```
/var/log/terminal-fm/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0644 terminal-fm terminal-fm
}
```

## Troubleshooting

### Service Won't Start

```bash
# Check status
sudo systemctl status terminal-fm

# View full logs
sudo journalctl -u terminal-fm -n 50

# Check port binding
sudo lsof -i :22

# Test binary manually
sudo -u terminal-fm /var/lib/terminal-fm/terminal-fm --port 2223
```

### Port 22 Already in Use

```bash
# Check what's using port 22
sudo lsof -i :22

# If OpenSSH is using it, reconfigure SSH first
# Move OpenSSH to port 2222 (see "Configure Admin SSH" section)
```

### Cannot Connect

```bash
# Test from server
ssh localhost

# Check firewall
sudo ufw status

# Check DNS
dig terminal.fm

# Test direct IP
ssh <YOUR_VPS_IP>

# Check service
sudo systemctl status terminal-fm
```

### High Memory Usage

```bash
# Check memory
free -h

# Check per-process memory
ps aux | grep terminal-fm

# Restart service
sudo systemctl restart terminal-fm
```

### Database Locked Errors

```bash
# Stop service
sudo systemctl stop terminal-fm

# Check database integrity
sqlite3 /var/lib/terminal-fm/data.db "PRAGMA integrity_check;"

# Enable WAL mode
sqlite3 /var/lib/terminal-fm/data.db "PRAGMA journal_mode=WAL;"

# Restart service
sudo systemctl start terminal-fm
```

## Updating Terminal.FM

```bash
# Stop service
sudo systemctl stop terminal-fm

# Backup database
sudo cp /var/lib/terminal-fm/data.db /var/lib/terminal-fm/backups/

# Pull latest code
cd ~/terminal-fm
git pull origin main

# Rebuild
go build -ldflags="-s -w" -o terminal-fm ./cmd/server

# Replace binary
sudo cp terminal-fm /var/lib/terminal-fm/terminal-fm
sudo chown terminal-fm:terminal-fm /var/lib/terminal-fm/terminal-fm

# Start service
sudo systemctl start terminal-fm

# Check status
sudo systemctl status terminal-fm
```

## Backup Strategy

### Automated Backups

```bash
# Create backup script
sudo nano /usr/local/bin/backup-terminal-fm.sh
```

Add:

```bash
#!/bin/bash
BACKUP_DIR="/var/lib/terminal-fm/backups"
DATE=$(date +%Y%m%d-%H%M%S)

mkdir -p $BACKUP_DIR

# Backup database
cp /var/lib/terminal-fm/data.db $BACKUP_DIR/data-$DATE.db

# Keep only last 7 days
find $BACKUP_DIR -name "data-*.db" -mtime +7 -delete

echo "Backup completed: $DATE"
```

```bash
# Make executable
sudo chmod +x /usr/local/bin/backup-terminal-fm.sh

# Add to crontab (daily at 3 AM)
sudo crontab -e
```

Add:

```
0 3 * * * /usr/local/bin/backup-terminal-fm.sh >> /var/log/terminal-fm/backup.log 2>&1
```

## Cost Estimation

### Monthly Costs

- **VPS**: €7/month (OVH Value)
- **Domain**: €0.83-8.33/month (€10-100/year)
- **Total**: ~€8-15/month

### Scaling Costs

- 200+ users: Upgrade to VPS (€15/month)
- 500+ users: Consider multi-instance + load balancer

---

For development setup, see [DEVELOPMENT.md](DEVELOPMENT.md).

For architecture details, see [../ARCHITECTURE.md](../ARCHITECTURE.md).
