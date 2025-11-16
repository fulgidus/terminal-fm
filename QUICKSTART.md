# 🚀 Quick Deploy Guide - Terminal-Radio

Complete automated deployment in **5 minutes**!

## ✅ Prerequisites

- VPS with Ubuntu 22.04 (OVH, Hetzner, DigitalOcean, etc.)
- SSH access to VPS
- Domain name (optional but recommended)

## 📋 One-Command Setup

### Step 1: Setup VPS (5 minutes)

```bash
# SSH into your VPS
ssh ubuntu@51.91.97.241

# Download and run setup script
curl -fsSL https://raw.githubusercontent.com/fulgidus/terminal-fm/main/scripts/setup-vps.sh | bash
```

**What it does:**
- ✅ Installs Go 1.21, FFmpeg, dependencies
- ✅ Creates dedicated user and directories
- ✅ Clones repository and builds application
- ✅ Sets up systemd service
- ✅ Configures firewall
- ✅ Moves SSH admin to port 2222
- ✅ Starts Terminal-Radio on port 22

### Step 2: Setup GitHub Actions (3 minutes)

#### 2a. Generate SSH Key (on your local machine)

```bash
# Generate dedicated deploy key
ssh-keygen -t ed25519 -C "github-actions@terminal-radio" -f ~/.ssh/terminal-radio-deploy

# Don't set a passphrase (press Enter)

# Display public key
cat ~/.ssh/terminal-radio-deploy.pub
```

#### 2b. Add Public Key to VPS

```bash
# SSH into VPS (note the new port!)
ssh -p 2222 ubuntu@51.91.97.241

# Add public key
nano ~/.ssh/authorized_keys
# Paste the public key and save (Ctrl+X, Y, Enter)
```

#### 2c. Add Private Key to GitHub Secrets

1. Go to: https://github.com/fulgidus/terminal-fm/settings/secrets/actions
2. Click **"New repository secret"**
3. Name: `SSH_PRIVATE_KEY`
4. Value: Content of `~/.ssh/terminal-radio-deploy` (the PRIVATE key)
   ```bash
   # Display private key to copy
   cat ~/.ssh/terminal-radio-deploy
   ```
5. Click **"Add secret"**

### Step 3: Test Deployment

```bash
# Make a small change
echo "# Test CI/CD" >> README.md

# Commit and push
git add .
git commit -m "test: CI/CD deployment"
git push origin main

# Watch deployment
# Go to: https://github.com/fulgidus/terminal-fm/actions
```

### Step 4: Configure DNS (if you have a domain)

1. Go to your DNS provider (OVH, Cloudflare, etc.)
2. Add A record:
   ```
   Type: A
   Name: @ (or subdomain)
   Value: 51.91.97.241
   TTL: 3600
   ```
3. Wait for propagation (1-48 hours, usually 1-2 hours)

## 🎵 Test Terminal-Radio

```bash
# Via IP
ssh 51.91.97.241

# Via domain (after DNS propagation)
ssh terminal-radio.com
```

## 📊 Monitoring

```bash
# SSH into VPS (admin port)
ssh -p 2222 ubuntu@51.91.97.241

# Run monitoring dashboard
./opt/terminal-radio/scripts/monitor.sh
```

## 🔄 Update Process

**Automatic (via CI/CD):**
```bash
# Just push to main!
git push origin main
```

**Manual:**
```bash
ssh -p 2222 ubuntu@51.91.97.241
sudo /usr/local/bin/terminal-radio-deploy.sh
```

## 🛠️ Useful Commands

```bash
# Check service status
sudo systemctl status terminal-radio

# View logs
sudo journalctl -u terminal-radio -f

# Restart service
sudo systemctl restart terminal-radio

# Check active connections
sudo ss -tnp | grep ':22' | grep ESTAB

# Check resource usage
htop
```

## 🐛 Troubleshooting

### Service won't start

```bash
# Check logs
sudo journalctl -u terminal-radio -n 50

# Test binary manually
sudo -u terminal-radio /opt/terminal-radio/terminal-radio --port 2223
```

### Cannot connect via SSH

```bash
# Check firewall
sudo ufw status

# Check service
sudo systemctl status terminal-radio

# Test from VPS itself
ssh localhost
```

### GitHub Actions failing

```bash
# Check SSH key is added correctly
ssh -p 2222 -i ~/.ssh/terminal-radio-deploy ubuntu@51.91.97.241

# Check deploy script permissions
ls -l /usr/local/bin/terminal-radio-deploy.sh
```

## 📚 Documentation

- [Full Deployment Guide](docs/DEPLOYMENT.md)
- [Architecture](docs/ARCHITECTURE.md)
- [Development Guide](docs/DEVELOPMENT.md)

## 🎉 You're Done!

Your Terminal-Radio is now:
- ✅ Live and accessible via SSH
- ✅ Automatically deployed on every push to main
- ✅ Monitored and logged
- ✅ Secured with firewall

**Enjoy your radio! 🎵**
