# Testing Guide for Terminal.FM

## Automated Tests

### Running All Tests
```bash
go test ./...
```

### Running Specific Package Tests
```bash
# Player tests
go test -v ./pkg/services/player

# Storage tests
go test -v ./pkg/services/storage

# Radio Browser API tests
go test -v ./pkg/services/radiobrowser
```

### Running Tests with Coverage
```bash
go test -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Manual Testing

### 1. Basic Server Startup
```bash
# Build the binary
go build -o terminal-fm ./cmd/server

# Start in development mode
./terminal-fm --dev --port 2222
```

**Expected**: Server starts without errors and logs "Starting SSH server on 0.0.0.0:2222"

### 2. SSH Connection
```bash
# In another terminal
ssh localhost -p 2222
```

**Expected**: You should see the Terminal.FM TUI with a list of stations

### 3. Player Lifecycle Testing

#### Test 3a: Start Playback
1. Connect via SSH
2. Use arrow keys to select a station
3. Press `Enter` or `Space` to play

**Expected**: 
- Station should start playing
- Status bar shows "Playing: [Station Name]"
- One ffplay process should be running: `ps aux | grep ffplay`

#### Test 3b: Stop Playback
1. While a station is playing
2. Press `s` to stop

**Expected**:
- Audio stops
- Status bar shows "Stopped"
- No ffplay processes running: `ps aux | grep ffplay`

#### Test 3c: Switch Stations
1. Play a station
2. Select a different station with arrow keys
3. Press `Enter` to play the new station

**Expected**:
- Old station stops
- New station starts playing
- Only ONE ffplay process running
- No orphaned processes

#### Test 3d: Volume Control
1. While playing a station
2. Press `+` or `=` to increase volume
3. Press `-` to decrease volume

**Expected**:
- Volume changes take effect (station restarts with new volume)
- No orphaned processes

### 4. Cleanup Testing

#### Test 4a: Normal Quit
1. Connect and play a station
2. Press `q` to quit

**Expected**:
- SSH session ends cleanly
- No ffplay processes remain: `ps aux | grep ffplay`

#### Test 4b: Ctrl+C Quit
1. Connect and play a station
2. Press `Ctrl+C` to force quit

**Expected**:
- SSH session ends
- No ffplay processes remain

#### Test 4c: SSH Disconnect
1. Connect and play a station
2. Close terminal window or kill SSH client

**Expected**:
- Server logs "SSH session ended, cleaning up player..."
- No ffplay processes remain after ~1 second

### 5. Concurrent Sessions
1. Open 3-4 SSH sessions simultaneously
2. Play different stations in each session
3. Check process count: `ps aux | grep ffplay | wc -l`

**Expected**:
- Number of ffplay processes matches number of playing sessions
- Each session can control its own playback independently

4. Quit all sessions
5. Verify no processes remain: `ps aux | grep ffplay`

### 6. Search Functionality
1. Press `/` to open search
2. Enter a station name or 2-letter country code (e.g., "BBC" or "US")
3. Press `Enter` to search
4. Press `Tab` to focus results
5. Use arrow keys to select a result
6. Press `Enter` to play

**Expected**:
- Search returns relevant results
- Playing from search results works
- No process leaks

### 7. Bookmarks
1. Select a station
2. Press `a` to add bookmark
3. Press `b` to view bookmarks
4. Verify station appears in bookmarks

**Expected**:
- Bookmarks persist in `~/.terminal-fm/terminal-fm.db`
- Bookmarks view shows saved stations

### 8. Stress Testing

#### Rapid Station Switching
1. Connect via SSH
2. Rapidly press `Down` + `Enter` repeatedly for 10-20 seconds

**Expected**:
- No crashes
- Only one ffplay process at a time
- No zombie processes

#### Long-Running Session
1. Connect and play a station
2. Leave it running for 30+ minutes
3. Press `s` to stop

**Expected**:
- Clean stop after long playback
- No memory leaks (check with `top` or `htop`)

## Debugging Commands

### Check Running Processes
```bash
# Find all ffplay processes
ps aux | grep ffplay | grep -v grep

# Count ffplay processes
ps aux | grep ffplay | grep -v grep | wc -l

# Kill orphaned ffplay processes (if any)
killall ffplay
```

### Check Database
```bash
# View bookmarks in SQLite
sqlite3 ~/.terminal-fm/terminal-fm.db "SELECT * FROM bookmarks;"
```

### Monitor Server Logs
```bash
# Server logs to stdout/stderr
./terminal-fm --dev --port 2222 2>&1 | tee server.log
```

### Network Testing
```bash
# Test without network
sudo iptables -A OUTPUT -p tcp --dport 80 -j DROP
./terminal-fm --dev --port 2222

# Restore network
sudo iptables -D OUTPUT -p tcp --dport 80 -j DROP
```

## Common Issues

### Issue: Orphaned ffplay processes
**Symptoms**: Multiple ffplay processes remain after quit
**Fix**: 
1. Kill processes: `killall ffplay`
2. Check player cleanup code in `pkg/services/player/player.go`
3. Verify `Cleanup()` is called in `pkg/ui/model.go`

### Issue: "address already in use"
**Symptoms**: Server won't start on port 2222
**Fix**: Kill existing server: `pkill terminal-fm` or change port

### Issue: No audio output
**Symptoms**: Station shows as playing but no sound
**Fix**:
1. Check ffplay is installed: `which ffplay`
2. Test ffplay directly: `ffplay -nodisp -autoexit [stream_url]`
3. Check system audio: `pactl list sinks` or `aplay -l`

### Issue: Database locked
**Symptoms**: "database is locked" errors
**Fix**: Close all instances: `pkill terminal-fm && rm ~/.terminal-fm/terminal-fm.db-*`

## Test Checklist

- [ ] Server starts in dev mode
- [ ] SSH connection succeeds
- [ ] Browse stations with arrow keys
- [ ] Play station with Enter/Space
- [ ] Stop station with 's' key
- [ ] Volume up/down with +/-
- [ ] Switch between stations (no orphaned processes)
- [ ] Search by name
- [ ] Search by country code
- [ ] Add bookmark
- [ ] View bookmarks
- [ ] Quit with 'q' (cleanup verified)
- [ ] Quit with Ctrl+C (cleanup verified)
- [ ] SSH disconnect (cleanup verified)
- [ ] Multiple concurrent sessions
- [ ] Rapid station switching (stress test)
- [ ] No ffplay processes remain after all tests

## Success Criteria

1. **No Orphaned Processes**: After every test, `ps aux | grep ffplay` should return no results
2. **Clean Shutdown**: Server logs should show "SSH session ended, cleaning up player..."
3. **Responsive UI**: All keyboard commands work without lag
4. **Stable Playback**: Audio plays without interruption (except during switching)
5. **No Crashes**: Server and client remain stable through all tests
