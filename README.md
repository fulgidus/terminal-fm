
```
 _____ _____ ____  __  __ ___ _   _    _    _       _____ __  __ 
|_   _| ____|  _ \|  \/  |_ _| \ | |  / \  | |     |  ___|  \/  |
  | | |  _| | |_) | |\/| || ||  \| | / _ \ | |     | |_  | |\/| |
  | | | |___|  _ <| |  | || || |\  |/ ___ \| |___  |  _| | |  | |
  |_| |_____|_| \_\_|  |_|___|_| \_/_/   \_\_____| |_|   |_|  |_|
                                                                   
      Listen to 30,000+ radio stations from your terminal
```
> A TUI (Text User Interface) radio player accessible via SSH

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![Built with Bubbletea](https://img.shields.io/badge/Built%20with-Bubbletea-pink)](https://github.com/charmbracelet/bubbletea)


## 🚀 Quick Start

```bash
ssh terminal-radio.com
```

That's it! No installation required. Just SSH and start listening.

> **Note**: Currently in beta. IP access: `ssh 51.91.97.241`

---

## ✨ Features

### v1.0 (Current)
- 🌍 **30,000+ Radio Stations** - Access to Radio Browser community database
- 🎵 **Audio Playback** - FFplay-based player with volume control
- 📊 **Station Metadata** - Name, country, bitrate, codec, votes
- ⭐ **Bookmarks System** - SQLite-backed persistent favorites
- 🔍 **Interactive Search** - Search by name or country code with live results
- 🌐 **Multilingual** - Full i18n support (English and Italian)
- 📱 **Responsive TUI** - Adapts to any terminal size with styled components
- 🔐 **Anonymous SSH** - No authentication required
- 💾 **Local Storage** - Bookmarks saved to `~/.terminal-fm/terminal-fm.db`
- 🎨 **Beautiful UI** - Styled with Lipgloss (cyan/pink/purple theme)

### 🚧 Roadmap (v1.5+)
- 📈 Real-time spectrum analyzer (exploring WebRTC/client-side solutions)
- 📜 Listening history
- 💾 Stream recording
- 🎸 Last.fm scrobbling
- 📝 Lyrics display
- 👥 Multi-user listening rooms

## 🎮 Usage

### Keyboard Controls

**Navigation**
```
↑/↓ or k/j     Navigate stations
PgUp/PgDn      Fast scroll
Home/End       Jump to first/last
```

**Playback**
```
Enter/Space    Play selected station
s              Stop playback
+/-            Volume up/down (10% increments)
```

**Features**
```
a              Add/Remove bookmark
b              Toggle bookmarks view
/              Search stations
?              Show help
q or Ctrl+C    Quit
```

### Search
Press `/` to open search, then:
- Enter station name to search
- Enter 2-letter country code (e.g., `IT`, `US`, `UK`)
- Press `Tab` to switch between input and results
- Press `Enter` to execute search or play selected result
- Press `ESC` to return to browse view

### Filters
- Genre (Jazz, Rock, Electronic, Classical, etc.)
- Country (Italy, USA, UK, Germany, etc.)
- Language (Italian, English, Spanish, etc.)
- Bitrate (64kbps, 128kbps, 320kbps)

## 🏗️ Architecture

Terminal.FM uses a **client-streaming architecture** where audio streams directly from radio stations to your local machine, keeping server costs minimal and scalability infinite.

```
User Terminal → SSH (terminal.fm) → TUI + Metadata
                                  ↓
                         Stream URL + mpv player
                                  ↓
                    Radio Station (direct stream)
```

**Key components:**
- **SSH Server**: Charm Wish (Go) on port 22
- **TUI Framework**: Bubbletea
- **Audio Player**: mpv/ffplay (local)
- **API**: Radio Browser Community API
- **Storage**: SQLite (user preferences)

See [ARCHITECTURE.md](docs/ARCHITECTURE.md) for detailed technical documentation.

## 🛠️ Self-Hosting

Want to run your own instance? **Automated deployment in 5 minutes!**

### ⚡ Quick Deploy (Fully Automated)

```bash
# 1. SSH into your VPS
ssh ubuntu@YOUR_VPS_IP

# 2. Run automated setup
curl -fsSL https://raw.githubusercontent.com/fulgidus/terminal-fm/main/scripts/setup-vps.sh | bash

# 3. Setup GitHub Actions for auto-deploy on push
# See QUICKSTART.md for detailed instructions
```

**What you get:**
- ✅ Complete VPS setup (Go, FFmpeg, dependencies)
- ✅ Systemd service running on port 22
- ✅ Admin SSH moved to port 2222
- ✅ Firewall configured
- ✅ CI/CD ready with GitHub Actions
- ✅ Auto-deploy on every push to main

### 📋 Requirements
- Ubuntu 22.04 VPS (1GB RAM minimum)
- SSH access
- GitHub account (for CI/CD)
- Domain name (optional)

### 📚 Deployment Guides
- **[QUICKSTART.md](QUICKSTART.md)** - 5-minute automated setup
- **[docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)** - Detailed manual deployment

### 🔄 CI/CD Pipeline

Every push to `main` automatically:
1. Runs tests
2. Builds binary
3. Deploys to VPS
4. Restarts service
5. Verifies deployment

**See deployment workflow**: [`.github/workflows/deploy.yml`](.github/workflows/deploy.yml)

### 🧑‍💻 Development Mode

```bash
git clone https://github.com/fulgidus/terminal-fm.git
cd terminal-fm
go mod download

# Run in dev mode (mock data, no ffplay needed)
go run ./cmd/server --dev --port 2222

# Connect in another terminal
ssh localhost -p 2222
```

## 🤝 Contributing

We welcome contributions! Whether it's:
- 🐛 Bug reports
- ✨ Feature requests
- 🌐 Translations
- 💻 Code contributions
- 📚 Documentation improvements

Please read [CONTRIBUTING.md](CONTRIBUTING.md) before submitting PRs.

### Development Setup
```bash
git clone https://github.com/fulgidus/terminal-fm.git
cd terminal-fm
go mod download

# Run in development mode (uses mock data, no real API calls)
go run ./cmd/server --dev --port 2222

# In another terminal, connect
ssh localhost -p 2222
```

See [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for detailed development guide.

### Available Flags
```bash
--dev          Enable development mode (mock API client)
--port         SSH port to listen on (default: 22)
--host         Host address to bind to (default: 0.0.0.0)
--version      Show version information
```

## 📖 Documentation

- **[Quick Start](QUICKSTART.md)** - 5-minute automated deployment
- [Architecture Overview](docs/ARCHITECTURE.md)
- [Contributing Guidelines](CONTRIBUTING.md)
- [Deployment Guide](docs/DEPLOYMENT.md)
- [API Integration](docs/API.md)
- [Internationalization](docs/I18N.md)
- [Development Setup](docs/DEVELOPMENT.md)
- [Security Policy](SECURITY.md)

## 🌟 Why Terminal.FM?

- **Zero Installation**: Works on any machine with SSH
- **Lightweight**: Minimal resource usage
- **Privacy-Focused**: No tracking, no accounts, no data collection
- **Open Source**: GPLv3 licensed, community-driven
- **Cross-Platform**: Works on Linux, macOS, Windows (WSL), even mobile with SSH clients
- **Nostalgic**: Brings back the joy of radio in the terminal era

## 🙏 Credits

- **Radio Database**: [Radio Browser](https://www.radio-browser.info/) - Community-driven radio station directory
- **TUI Framework**: [Bubbletea](https://github.com/charmbracelet/bubbletea) by Charm
- **SSH Library**: [Wish](https://github.com/charmbracelet/wish) by Charm
- **Inspired by**: [terminal.shop](https://github.com/charmbracelet/termshop)

## 📜 License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## 💬 Community

- **Issues**: [GitHub Issues](https://github.com/fulgidus/terminal-fm/issues)
- **Discussions**: [GitHub Discussions](https://github.com/fulgidus/terminal-fm/discussions)

---

**Made with ❤️ for terminal enthusiasts and radio lovers**

`ssh terminal-radio.com` and enjoy! 🎵
