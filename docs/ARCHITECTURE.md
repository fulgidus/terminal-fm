# Architecture

This document describes the technical architecture of Terminal.FM.

## Table of Contents

- [Overview](#overview)
- [System Architecture](#system-architecture)
- [Components](#components)
- [Data Flow](#data-flow)
- [Technology Stack](#technology-stack)
- [Database Schema](#database-schema)
- [Deployment Architecture](#deployment-architecture)
- [Security Considerations](#security-considerations)
- [Performance & Scalability](#performance--scalability)

## Overview

Terminal.FM is a lightweight TUI radio player accessible via SSH. The architecture prioritizes:
- **Scalability**: Direct streaming keeps server load minimal
- **Simplicity**: Single Go binary with embedded assets
- **Privacy**: No user tracking, minimal data storage
- **Performance**: Efficient TUI rendering and API caching

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         User Terminal                           │
│                     $ ssh terminal.fm                           │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  SSH Client                                               │  │
│  │  - Handles terminal I/O                                   │  │
│  │  - Renders TUI (Bubbletea)                                │  │
│  └───────────────────────────────────────────────────────────┘  │
│                              │                                  │
│                              │ SSH Protocol (port 22)           │
│                              ▼                                  │
└─────────────────────────────────────────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────────┐
│                    Terminal.FM Server (VPS)                     │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  SSH Server Layer (Charm Wish)                            │  │
│  │  - Port 22 binding                                        │  │
│  │  - Anonymous authentication                               │  │
│  │  - Session management                                     │  │
│  │  - PTY allocation                                         │  │
│  └───────────────────────────────────────────────────────────┘  │
│                              │                                  │
│  ┌───────────────────────────▼───────────────────────────────┐  │
│  │  Application Layer (Bubbletea TUI)                        │  │
│  │  ┌─────────────┐  ┌──────────────┐  ┌─────────────────┐   │  │
│  │  │   Browser   │  │    Player    │  │   Bookmarks     │   │  │
│  │  │   Component │  │   Component  │  │   Component     │   │  │
│  │  └─────────────┘  └──────────────┘  └─────────────────┘   │  │
│  │  ┌─────────────┐  ┌──────────────┐  ┌─────────────────┐   │  │
│  │  │   Search    │  │   Metadata   │  │   Settings      │   │  │
│  │  │   Component │  │   Component  │  │   Component     │   │  │
│  │  └─────────────┘  └──────────────┘  └─────────────────┘   │  │
│  └───────────────────────────────────────────────────────────┘  │
│                              │                                  │
│  ┌───────────────────────────▼───────────────────────────────┐  │
│  │  Service Layer                                            │  │
│  │  ┌──────────────┐  ┌────────────┐  ┌──────────────────┐   │  │
│  │  │ Radio Browser│  │   Player   │  │   User Prefs     │   │  │
│  │  │ API Client   │  │   Service  │  │   Service        │   │  │
│  │  └──────────────┘  └────────────┘  └──────────────────┘   │  │
│  └───────────────────────────────────────────────────────────┘  │
│                              │                                  │
│  ┌───────────────────────────▼───────────────────────────────┐  │
│  │  Data Layer                                               │  │
│  │  ┌──────────────────┐  ┌────────────────────────────────┐ │  │
│  │  │   SQLite DB      │  │   In-Memory Cache (Redis?)     │ │  │
│  │  │   - Bookmarks    │  │   - Station lists              │ │  │
│  │  │   - User prefs   │  │   - API responses (TTL: 1h)    │ │  │
│  │  └──────────────────┘  └────────────────────────────────┘ │  │
│  └───────────────────────────────────────────────────────────┘  │
│                              │                                  │
└──────────────────────────────┼──────────────────────────────────┘
                               │
        ┌──────────────────────┼──────────────────────┐
        │                      │                      │
        ▼                      ▼                      ▼
┌───────────────┐    ┌──────────────────┐    ┌──────────────┐
│ Radio Browser │    │  mpv/ffplay      │    │ Radio Station│
│ API           │    │  (Local Process) │◄───│ Stream Server│
│ (External)    │    │  Spawned via SSH │    │ (External)   │
└───────────────┘    └──────────────────┘    └──────────────┘
```

## Components

### 1. SSH Server Layer (`pkg/ssh/`)

**Responsibility**: Handle SSH connections and authentication

**Implementation**: 
- Charm Wish library
- Anonymous access (no password required)
- Optional public key authentication for advanced features
- Session management with unique IDs
- PTY allocation for terminal control

**Key Files**:
```
pkg/ssh/
├── server.go          # Main SSH server
├── auth.go            # Authentication handlers
├── middleware.go      # Logging, rate limiting
└── session.go         # Session management
```

### 2. Application Layer (`pkg/ui/`)

**Responsibility**: TUI rendering and user interaction

**Implementation**: Bubbletea MVC pattern

**Components**:

#### Browser Component
- Lists stations from Radio Browser API
- Filterable by genre, country, language, bitrate
- Keyboard navigation (Vim-style + arrows)
- Pagination (50 stations per page)

#### Player Component
- Controls playback (play/pause/stop/volume)
- Displays current station info
- Shows connection status
- Spawns local mpv/ffplay process

#### Bookmarks Component
- User's saved stations
- CRUD operations on favorites
- Persistent storage per user (SQLite)

#### Search Component
- Real-time search across station names/tags
- Fuzzy matching support
- History of recent searches

#### Metadata Component
- Shows "Now Playing" information
- Displays bitrate, codec, sample rate
- Updates in real-time from stream metadata

#### Settings Component
- Language selection (i18n)
- Player preferences (mpv vs ffplay)
- Volume defaults
- Key binding customization

**Key Files**:
```
pkg/ui/
├── model.go           # Main Bubbletea model
├── update.go          # Update logic
├── view.go            # Rendering logic
├── components/
│   ├── browser.go
│   ├── player.go
│   ├── bookmarks.go
│   ├── search.go
│   ├── metadata.go
│   └── settings.go
└── styles/
    └── theme.go       # Lipgloss styling
```

### 3. Service Layer (`pkg/services/`)

**Responsibility**: Business logic and external integrations

#### Radio Browser Service (`radiobrowser/`)
- API client for community.radio-browser.info
- Caching layer (1 hour TTL for station lists)
- Error handling and retries
- Endpoints:
  - `GET /json/stations/search`
  - `GET /json/stations/bycountry/{country}`
  - `GET /json/stations/bygenre/{genre}`
  - `GET /json/stations/bylanguage/{language}`

#### Player Service (`player/`)
- Spawns mpv/ffplay as subprocess
- Manages audio stream lifecycle
- Volume control via IPC
- Metadata extraction from stream
- Graceful shutdown on disconnect

#### User Preferences Service (`userprefs/`)
- SQLite storage per user (identified by SSH fingerprint)
- Bookmark management
- Settings persistence
- Migration system for schema updates

**Key Files**:
```
pkg/services/
├── radiobrowser/
│   ├── client.go      # HTTP client
│   ├── models.go      # API response types
│   └── cache.go       # Caching logic
├── player/
│   ├── player.go      # Main player interface
│   ├── mpv.go         # mpv implementation
│   └── metadata.go    # Stream metadata parser
└── userprefs/
    ├── storage.go     # SQLite operations
    └── migrations.go  # Schema migrations
```

### 4. Data Layer (`pkg/storage/`)

**Responsibility**: Data persistence

#### SQLite Database
- Single file per deployment
- Schema version tracking
- Automatic migrations

#### In-Memory Cache (Optional)
- Consider Redis for multi-instance deployments
- For v1: Go map with mutex (single instance)

**Key Files**:
```
pkg/storage/
├── db.go              # Database initialization
├── models.go          # Data models
└── queries.go         # SQL queries
```

### 5. i18n Layer (`pkg/i18n/`)

**Responsibility**: Internationalization

**Implementation**: go-i18n library

**Supported Languages (v1)**:
- English (en-US)
- Italian (it-IT)

**Key Files**:
```
pkg/i18n/
├── i18n.go            # Main i18n logic
└── locales/
    ├── active.en.toml
    └── active.it.toml
```

## Data Flow

### 1. User Connects via SSH

```
User → ssh terminal.fm
       ↓
[SSH Server] Accept connection
       ↓
[Auth Middleware] Assign session ID
       ↓
[Application] Initialize Bubbletea model
       ↓
[Browser Component] Fetch stations (cached)
       ↓
[View] Render initial UI
```

### 2. User Searches for Station

```
User → Types "/jazz"
       ↓
[Search Component] Update search filter
       ↓
[Radio Browser Service] API call: GET /json/stations/search?name=jazz
       ↓
[Cache] Store results (1h TTL)
       ↓
[Browser Component] Update station list
       ↓
[View] Re-render with filtered results
```

### 3. User Plays Station

```
User → Press Enter on station
       ↓
[Player Component] Send play command
       ↓
[Player Service] Spawn mpv process with stream URL
       ↓
mpv → Connects directly to radio station
       ↓
[Metadata Component] Poll mpv for metadata
       ↓
[View] Display "Now Playing" info
```

### 4. User Bookmarks Station

```
User → Press 'b' on station
       ↓
[Bookmarks Component] Add bookmark
       ↓
[User Prefs Service] INSERT INTO bookmarks
       ↓
[SQLite] Persist to disk
       ↓
[View] Show confirmation
```

## Technology Stack

### Backend
- **Language**: Go 1.21+
- **SSH Server**: [Charm Wish](https://github.com/charmbracelet/wish) v1.3+
- **TUI Framework**: [Bubbletea](https://github.com/charmbracelet/bubbletea) v0.25+
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss) v0.9+
- **i18n**: [go-i18n](https://github.com/nicksnyder/go-i18n) v2.4+
- **Database**: SQLite 3
- **HTTP Client**: Go stdlib `net/http`

### External Dependencies
- **Audio Player**: mpv or ffplay (user's local system)
- **API**: [Radio Browser API](https://www.radio-browser.info/)

### Development Tools
- **Testing**: Go stdlib `testing`
- **Linting**: golangci-lint
- **CI/CD**: GitHub Actions
- **Deployment**: Systemd service or Docker

## Database Schema

### `users` table
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ssh_fingerprint TEXT UNIQUE NOT NULL,  -- User identifier
    language TEXT DEFAULT 'en-US',          -- Preferred language
    volume INTEGER DEFAULT 50,              -- Default volume (0-100)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_seen_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_fingerprint ON users(ssh_fingerprint);
```

### `bookmarks` table
```sql
CREATE TABLE bookmarks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    station_uuid TEXT NOT NULL,             -- Radio Browser station UUID
    station_name TEXT NOT NULL,
    station_url TEXT NOT NULL,
    added_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_bookmarks_user ON bookmarks(user_id);
CREATE INDEX idx_bookmarks_station ON bookmarks(station_uuid);
```

### `listening_history` table (future)
```sql
CREATE TABLE listening_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    station_uuid TEXT NOT NULL,
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    duration_seconds INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### `schema_migrations` table
```sql
CREATE TABLE schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Deployment Architecture

### Production Deployment (OVH VPS)

```
                    Internet
                        │
                        │ DNS: terminal.fm → <VPS_IP>
                        ▼
                ┌───────────────┐
                │   Firewall    │
                │  - Port 22    │ SSH (Terminal.FM)
                │  - Port 2222  │ SSH (Admin access)
                └───────┬───────┘
                        │
                ┌───────▼────────────────────────┐
                │     OVH VPS (Ubuntu 22.04)     │
                │  - 2 vCPU, 4GB RAM             │
                │  - 40GB SSD                    │
                │  - Unlimited bandwidth         │
                │                                │
                │  ┌──────────────────────────┐  │
                │  │  Systemd Service         │  │
                │  │  terminal-fm.service     │  │
                │  │  - Port 22               │  │
                │  │  - Auto-restart          │  │
                │  │  - Log to journald       │  │
                │  └──────────────────────────┘  │
                │                                │
                │  ┌──────────────────────────┐  │
                │  │  /var/lib/terminal-fm/   │  │
                │  │  - terminal-fm (binary)  │  │
                │  │  - data.db (SQLite)      │  │
                │  └──────────────────────────┘  │
                │                                │
                │  ┌──────────────────────────┐  │
                │  │  Monitoring (optional)   │  │
                │  │  - Prometheus exporter   │  │
                │  │  - Log rotation          │  │
                │  └──────────────────────────┘  │
                └────────────────────────────────┘
```

### Multi-Instance Deployment (Future)

For high availability:

```
                    Load Balancer
                    (HAProxy/nginx)
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
        ▼                 ▼                 ▼
    Instance 1        Instance 2        Instance 3
    (terminal-fm)     (terminal-fm)     (terminal-fm)
        │                 │                 │
        └─────────────────┴─────────────────┘
                          │
                          ▼
                  Shared PostgreSQL
                  (User prefs & bookmarks)
```

## Security Considerations

### SSH Security
- **Anonymous access**: Limited to read-only operations
- **Rate limiting**: Max 10 connections per IP per minute
- **Session timeout**: 30 minutes of inactivity
- **No shell access**: Forced command only (Terminal.FM app)
- **Admin SSH**: Separate port (2222) with key-based auth

### Data Security
- **No PII collection**: Only SSH fingerprint (hashed)
- **No passwords**: Anonymous or key-based auth
- **SQLite encryption**: Consider SQLCipher for sensitive data
- **HTTPS only**: All external API calls over TLS

### Input Validation
- **Station URLs**: Whitelist protocols (http, https, mms)
- **Search queries**: Sanitize to prevent injection
- **File paths**: No user-controllable paths

### Resource Limits
- **Max connections**: 100 concurrent users (v1)
- **Memory per session**: ~10MB
- **Database size**: Auto-rotate when >100MB
- **API rate limiting**: Respect Radio Browser limits

## Performance & Scalability

### Current Limits (v1 - Single VPS)
- **Concurrent users**: 100-200 (2 vCPU, 4GB RAM)
- **Bandwidth**: ~0 (direct streaming)
- **CPU**: <5% per user (TUI rendering only)
- **Memory**: ~10MB per user
- **Storage**: <1GB (database + binary)

### Bottlenecks
1. **Database**: SQLite write locks (bookmarks)
   - **Solution**: Connection pooling, WAL mode
2. **SSH connections**: File descriptor limits
   - **Solution**: Increase ulimit, use systemd limits
3. **API rate limits**: Radio Browser API
   - **Solution**: Aggressive caching (1h TTL)

### Optimization Strategies

#### Caching
- Station lists: 1 hour TTL
- Station metadata: 5 minutes TTL
- User preferences: In-memory after first load

#### Database Optimization
```sql
-- Enable WAL mode for better concurrency
PRAGMA journal_mode=WAL;

-- Increase cache size
PRAGMA cache_size=10000;

-- Use memory-mapped I/O
PRAGMA mmap_size=268435456; -- 256MB
```

#### Connection Pooling
- Max idle connections: 10
- Max open connections: 50
- Connection lifetime: 5 minutes

### Monitoring Metrics
- Active SSH connections
- Database query latency
- API response times
- Memory usage per session
- Error rates (by type)

### Scalability Path
1. **v1.0**: Single VPS, up to 200 users
2. **v1.5**: Horizontal scaling with shared PostgreSQL
3. **v2.0**: Kubernetes deployment with auto-scaling

---

For deployment instructions, see [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md).

For development setup, see [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md).
