# 🎉 Terminal-Radio - Deployment Package Ready!

## ✅ What's Been Created

### 1. **Automated VPS Setup** (`scripts/setup-vps.sh`)
One-command installation script that:
- Installs Go 1.21, FFmpeg, and all dependencies
- Creates dedicated user and directories
- Clones repository and builds application
- Sets up systemd service (auto-start on boot)
- Configures firewall (UFW)
- Moves admin SSH to port 2222
- Starts Terminal-Radio on port 22

### 2. **CI/CD Pipeline** (`.github/workflows/deploy.yml`)
GitHub Actions workflow that automatically:
- Runs tests on every push to main
- Runs linter (golangci-lint)
- Builds optimized binary
- Deploys to VPS via SSH
- Restarts service
- Verifies deployment success

### 3. **Monitoring Dashboard** (`scripts/monitor.sh`)
Real-time monitoring script showing:
- Service status and uptime
- Active connections count
- CPU and memory usage
- System resources
- Recent logs

### 4. **Quick Start Guide** (`QUICKSTART.md`)
Step-by-step guide for 5-minute deployment

### 5. **Updated Documentation**
- README.md updated with CI/CD info and terminal-radio.com
- Links to new automated deployment workflow

---

## 🚀 How to Deploy (Complete Workflow)

### Phase 1: VPS Setup (5 minutes)

```bash
# 1. SSH into your VPS
ssh ubuntu@51.91.97.241

# 2. Download and run setup script
wget https://raw.githubusercontent.com/fulgidus/terminal-fm/main/scripts/setup-vps.sh
chmod +x setup-vps.sh
./setup-vps.sh

# Script will:
# ✅ Install everything
# ✅ Build and start Terminal-Radio
# ✅ Configure firewall
# ✅ Move SSH to port 2222
# ⚠️  After this, use: ssh -p 2222 ubuntu@51.91.97.241
```

### Phase 2: GitHub Actions Setup (3 minutes)

#### Step 1: Generate Deploy Key (on your local machine)
```bash
ssh-keygen -t ed25519 -C "github-actions" -f ~/.ssh/terminal-radio-deploy
# Don't set passphrase (press Enter twice)

# Display public key
cat ~/.ssh/terminal-radio-deploy.pub
```

#### Step 2: Add Public Key to VPS
```bash
# SSH to VPS (NOTE: Now using port 2222!)
ssh -p 2222 ubuntu@51.91.97.241

# Add public key
nano ~/.ssh/authorized_keys
# Paste the public key at the end
# Save: Ctrl+X, Y, Enter

# Test it works
exit
ssh -p 2222 -i ~/.ssh/terminal-radio-deploy ubuntu@51.91.97.241
```

#### Step 3: Add Private Key to GitHub
1. Go to: https://github.com/fulgidus/terminal-fm/settings/secrets/actions
2. Click **"New repository secret"**
3. Name: `SSH_PRIVATE_KEY`
4. Value: 
   ```bash
   cat ~/.ssh/terminal-radio-deploy
   # Copy the ENTIRE output (including BEGIN/END lines)
   ```
5. Paste and click **"Add secret"**

### Phase 3: Test Auto-Deploy (1 minute)

```bash
# On your local machine, in the terminal-fm directory

# Make a small change
echo "# CI/CD Active!" >> README.md

# Commit and push
git add README.md
git commit -m "test: verify CI/CD deployment"
git push origin main

# Watch the magic happen! 🎉
# Go to: https://github.com/fulgidus/terminal-fm/actions
# You'll see the deployment running in real-time
```

### Phase 4: Configure DNS (Async)

While waiting for GitHub Actions to finish:

1. **Go to OVH Manager**: https://www.ovh.it/manager/
2. **Domains** → **terminal-radio.com** → **DNS Zone**
3. **Add A Record**:
   ```
   Type: A
   Subdomain: @ (or leave empty)
   Target: 51.91.97.241
   TTL: 3600
   ```
4. **Save** and wait (1-48 hours, usually 1-2 hours)

Test DNS propagation:
```bash
nslookup terminal-radio.com
# Should return: 51.91.97.241
```

---

## 🎵 Test Terminal-Radio

```bash
# Via IP (works immediately)
ssh 51.91.97.241

# Via domain (after DNS propagates)
ssh terminal-radio.com
```

---

## 📊 Monitoring

```bash
# SSH into VPS (admin port)
ssh -p 2222 ubuntu@51.91.97.241

# Run monitoring dashboard
/opt/terminal-radio/scripts/monitor.sh

# Or check service manually
sudo systemctl status terminal-radio
sudo journalctl -u terminal-radio -f
```

---

## 🔄 Daily Workflow

From now on, deploying updates is **automatic**:

```bash
# 1. Make changes locally
vim pkg/ui/view.go

# 2. Commit and push
git add .
git commit -m "feat: add new feature"
git push origin main

# 3. GitHub Actions automatically:
#    - Runs tests
#    - Builds binary
#    - Deploys to VPS
#    - Restarts service
#
# That's it! ✨
```

Watch deployments at: https://github.com/fulgidus/terminal-fm/actions

---

## 🛠️ Useful Commands

### On VPS (ssh -p 2222 ubuntu@51.91.97.241)

```bash
# Service management
sudo systemctl status terminal-radio
sudo systemctl restart terminal-radio
sudo journalctl -u terminal-radio -f

# Manual deployment
sudo /usr/local/bin/terminal-radio-deploy.sh

# Check connections
sudo ss -tnp | grep ':22' | grep ESTAB

# Monitoring dashboard
/opt/terminal-radio/scripts/monitor.sh
```

### On Local Machine

```bash
# Test SSH connection (users)
ssh 51.91.97.241

# Test SSH connection (admin)
ssh -p 2222 ubuntu@51.91.97.241

# Watch GitHub Actions
# Go to: https://github.com/fulgidus/terminal-fm/actions

# Test DNS
nslookup terminal-radio.com
dig terminal-radio.com
```

---

## 🔐 Security Notes

- ✅ Admin SSH on port 2222 (not exposed to users)
- ✅ Terminal-Radio on port 22 (public access)
- ✅ UFW firewall enabled
- ✅ Dedicated system user with minimal permissions
- ✅ Systemd hardening (NoNewPrivileges, PrivateTmp, etc.)
- ✅ Resource limits (512MB RAM, 50% CPU)

---

## 🐛 Troubleshooting

### GitHub Actions Fails

```bash
# Check SSH key works manually
ssh -p 2222 -i ~/.ssh/terminal-radio-deploy ubuntu@51.91.97.241

# Check secret is added correctly on GitHub
# Go to: https://github.com/fulgidus/terminal-fm/settings/secrets/actions
```

### Service Won't Start

```bash
ssh -p 2222 ubuntu@51.91.97.241
sudo journalctl -u terminal-radio -n 50
sudo systemctl status terminal-radio
```

### Cannot Connect

```bash
# Check service is running
sudo systemctl status terminal-radio

# Check firewall
sudo ufw status

# Test locally on VPS
ssh localhost
```

---

## 📁 Project Structure

```
terminal-fm/
├── .github/
│   └── workflows/
│       └── deploy.yml          # CI/CD pipeline
├── scripts/
│   ├── setup-vps.sh            # Automated VPS setup
│   └── monitor.sh              # Monitoring dashboard
├── cmd/server/                 # Application entry point
├── pkg/                        # Application code
├── QUICKSTART.md               # Quick deploy guide
└── README.md                   # Updated with CI/CD info
```

---

## 🎯 Next Steps After DNS Propagates

1. **Test with domain**:
   ```bash
   ssh terminal-radio.com
   ```

2. **Update README** (optional):
   ```bash
   # Change from:
   > **Note**: Currently in beta. IP access: `ssh 51.91.97.241`
   
   # To:
   # (remove the note entirely once DNS works)
   ```

3. **Share with the world!** 🎉
   - Post on Reddit (r/linux, r/golang, r/commandline)
   - Share on Hacker News
   - Tweet about it
   - Add to awesome-tui lists

---

## 📝 Commit Message

This commit includes:
- ✅ Fully automated VPS setup script
- ✅ GitHub Actions CI/CD pipeline
- ✅ Monitoring dashboard
- ✅ Complete documentation
- ✅ Zero-downtime deployment system

**Everything is ready for production! 🚀**

---

## 🎉 Summary

You now have:
- 📦 **One-command VPS setup** (5 minutes)
- 🔄 **Automatic deployment** on every push to main
- 📊 **Real-time monitoring** dashboard
- 🔐 **Production-ready security** configuration
- 📚 **Complete documentation**

**Total setup time: ~10 minutes**
**Future deployments: Automatic on git push**

Enjoy your Terminal-Radio! 🎵
