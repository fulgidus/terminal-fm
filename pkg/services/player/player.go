// Package player provides audio playback functionality.
package player

import (
	"fmt"
	"os/exec"
	"sync"
	"syscall"

	"github.com/fulgidus/terminal-fm/pkg/services/radiobrowser"
)

// State represents the current playback state.
type State int

const (
	// StateStopped means no playback is active.
	StateStopped State = iota
	// StatePlaying means audio is currently playing.
	StatePlaying
	// StatePaused means playback is paused (not supported by all players).
	StatePaused
	// StateBuffering means the player is buffering.
	StateBuffering
)

// Player interface for audio playback.
type Player interface {
	// Play starts playing a radio station.
	Play(station *radiobrowser.Station) error
	// Stop stops the current playback.
	Stop() error
	// GetState returns the current playback state.
	GetState() State
	// GetCurrentStation returns the currently playing station, or nil.
	GetCurrentStation() *radiobrowser.Station
	// SetVolume sets the playback volume (0-100).
	SetVolume(volume int) error
	// GetVolume returns the current volume (0-100).
	GetVolume() int
}

// FFplayPlayer implements Player using ffplay.
type FFplayPlayer struct {
	mu             sync.RWMutex
	cmd            *exec.Cmd
	state          State
	currentStation *radiobrowser.Station
	volume         int
	ffplayPath     string
	processActive  bool
}

// NewFFplayPlayer creates a new ffplay-based player.
func NewFFplayPlayer(ffplayPath string) *FFplayPlayer {
	if ffplayPath == "" {
		ffplayPath = "ffplay"
	}

	return &FFplayPlayer{
		state:      StateStopped,
		volume:     70, // Default volume
		ffplayPath: ffplayPath,
	}
}

// Play starts playing a radio station.
func (p *FFplayPlayer) Play(station *radiobrowser.Station) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Stop any current playback
	if err := p.stopLocked(); err != nil {
		return fmt.Errorf("failed to stop current playback: %w", err)
	}

	if station == nil || station.URLResolved == "" {
		return fmt.Errorf("invalid station or URL")
	}

	// Build ffplay command
	// -nodisp: no video display
	// -loglevel quiet: suppress output
	// -autoexit: exit when playback ends
	// -volume: set volume (0-100)
	args := []string{
		"-nodisp",
		"-loglevel", "quiet",
		"-autoexit",
		"-volume", fmt.Sprintf("%d", p.volume),
		station.URLResolved,
	}

	p.cmd = exec.Command(p.ffplayPath, args...)

	// Start the player
	if err := p.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffplay: %w", err)
	}

	p.state = StatePlaying
	p.currentStation = station
	p.processActive = true

	// Monitor process in background
	cmd := p.cmd // Capture cmd for goroutine
	go func() {
		_ = cmd.Wait() // Ignore wait errors in goroutine

		// Cleanup after process exits
		p.mu.Lock()
		defer p.mu.Unlock()

		// Only clean up if this is still our active command
		if p.cmd == cmd {
			p.processActive = false
			p.state = StateStopped
			p.currentStation = nil
			p.cmd = nil
		}
	}()

	return nil
}

// Stop stops the current playback.
func (p *FFplayPlayer) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.stopLocked()
}

// stopLocked stops playback without acquiring the lock (internal use).
func (p *FFplayPlayer) stopLocked() error {
	// Mark as stopped first to prevent race conditions
	wasActive := p.processActive
	p.processActive = false
	p.state = StateStopped

	// Try to kill if we had an active process
	if wasActive && p.cmd != nil && p.cmd.Process != nil {
		// Send SIGTERM first (graceful)
		if err := p.cmd.Process.Signal(syscall.SIGTERM); err != nil {
			// If SIGTERM fails, try SIGKILL
			_ = p.cmd.Process.Signal(syscall.SIGKILL)
		}
		// Don't wait here - the goroutine will handle cleanup
	}

	p.currentStation = nil

	return nil
}

// GetState returns the current playback state.
func (p *FFplayPlayer) GetState() State {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state
}

// GetCurrentStation returns the currently playing station.
func (p *FFplayPlayer) GetCurrentStation() *radiobrowser.Station {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.currentStation
}

// SetVolume sets the playback volume (0-100).
// Note: This only affects new playback, not current playback.
func (p *FFplayPlayer) SetVolume(volume int) error {
	if volume < 0 || volume > 100 {
		return fmt.Errorf("volume must be between 0 and 100")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.volume = volume

	// If currently playing, restart with new volume
	if p.state == StatePlaying && p.currentStation != nil {
		station := p.currentStation
		if err := p.stopLocked(); err != nil {
			return err
		}

		// Restart playback with new volume
		// Note: We need to unlock before calling Play
		p.mu.Unlock()
		err := p.Play(station)
		p.mu.Lock()
		return err
	}

	return nil
}

// GetVolume returns the current volume.
func (p *FFplayPlayer) GetVolume() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.volume
}

// Cleanup forcefully stops playback and cleans up resources.
// Should be called when the session ends.
func (p *FFplayPlayer) Cleanup() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cmd != nil && p.cmd.Process != nil {
		// Force kill the process
		_ = p.cmd.Process.Signal(syscall.SIGKILL)
	}

	p.cmd = nil
	p.state = StateStopped
	p.currentStation = nil
	p.processActive = false

	return nil
}

// MpvPlayer implements Player using mpv (for future implementation).
type MpvPlayer struct {
	// TODO: Implement mpv-based player in future
}

// NewMpvPlayer creates a new mpv-based player.
func NewMpvPlayer(mpvPath string) *MpvPlayer {
	// TODO: Implement
	return &MpvPlayer{}
}
