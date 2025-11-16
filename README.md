
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
ssh terminal.fm
```

That's it! No installation required. Just SSH and start listening.

---

## ✨ Features

### v1.0 (Current)
- 🌍 **30,000+ Radio Stations** - Browse stations by genre, country, and language
- 🎵 **Integrated Player** - Play/pause/stop with volume control
- 📊 **Real-time Metadata** - See what's playing, bitrate, and codec info
- ⭐ **Bookmarks System** - Save your favorite stations
- 🔍 **Search** - Find stations by name, genre, or tags
- 🌐 **Multilingual** - Interface in English and Italian (more coming soon)
- 📱 **Responsive UI** - Adapts to any terminal size
- 🔐 **Anonymous Access** - No registration required

### 🚧 Roadmap (v1.5+)
- 📈 Real-time spectrum analyzer (exploring WebRTC/client-side solutions)
- 📜 Listening history
- 💾 Stream recording
- 🎸 Last.fm scrobbling
- 📝 Lyrics display
- 👥 Multi-user listening rooms

## 🎮 Usage

### Navigation
```
↑/↓ or k/j     Navigate stations
Enter          Play selected station
Space          Pause/Resume
q              Quit
/              Search
b              Bookmarks
?              Help
```

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

Want to run your own instance? See [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) for detailed setup instructions.

### Requirements
- Go 1.21+
- VPS with SSH access (port 22)
- Domain name (optional)
- ~1GB RAM minimum

### Quick Deploy
```bash
git clone https://github.com/fulgidus/terminal-fm.git
cd terminal-fm
go build -o terminal-fm ./cmd/server
sudo ./terminal-fm --port 22 --host 0.0.0.0
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
go run ./cmd/server --dev
```

See [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for detailed development guide.

## 📖 Documentation

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

`ssh terminal.fm` and enjoy! 🎵
