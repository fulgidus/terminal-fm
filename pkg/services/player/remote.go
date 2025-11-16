// Package player provides audio playback functionality.
package player

import (
	"fmt"
	"io"
	"sync"

	"github.com/fulgidus/terminal-fm/pkg/services/radiobrowser"
)

// RemotePlayer implements Player by sending control commands to a remote client.
// It uses custom OSC sequences to communicate with the terminal-radio client wrapper.
type RemotePlayer struct {
	mu             sync.RWMutex
	state          State
	currentStation *radiobrowser.Station
	volume         int
	writer         io.Writer
}

// Custom OSC sequences for terminal-radio client communication
const (
	oscPrefix = "\033]8888;"
	oscSuffix = "\007"
)

// NewRemotePlayer creates a new remote player that sends commands to the client wrapper.
func NewRemotePlayer(writer io.Writer) *RemotePlayer {
	return &RemotePlayer{
		state:  StateStopped,
		volume: 70,
		writer: writer,
	}
}

// Play starts playing a radio station by sending a PLAY command to the client.
func (p *RemotePlayer) Play(station *radiobrowser.Station) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if station == nil || station.URLResolved == "" {
		return fmt.Errorf("invalid station or URL")
	}

	if p.writer == nil {
		return fmt.Errorf("no output writer configured")
	}

	// Stop any current playback first (ignore errors)
	_ = p.sendCommand("STOP")

	// Send PLAY command with URL and volume
	cmd := fmt.Sprintf("PLAY;%s;%d", station.URLResolved, p.volume)
	if err := p.sendCommand(cmd); err != nil {
		return fmt.Errorf("failed to send play command: %w", err)
	}

	p.state = StatePlaying
	p.currentStation = station

	return nil
}

// Stop stops the current playback by sending a STOP command to the client.
func (p *RemotePlayer) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state == StateStopped {
		return nil
	}

	if err := p.sendCommand("STOP"); err != nil {
		return fmt.Errorf("failed to send stop command: %w", err)
	}

	p.state = StateStopped
	p.currentStation = nil

	return nil
}

// GetState returns the current playback state.
func (p *RemotePlayer) GetState() State {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state
}

// GetCurrentStation returns the currently playing station.
func (p *RemotePlayer) GetCurrentStation() *radiobrowser.Station {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.currentStation
}

// SetVolume sets the playback volume (0-100) by sending a VOLUME command.
func (p *RemotePlayer) SetVolume(volume int) error {
	if volume < 0 || volume > 100 {
		return fmt.Errorf("volume must be between 0 and 100")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.volume = volume

	// If currently playing, send volume update
	if p.state == StatePlaying {
		cmd := fmt.Sprintf("VOLUME;%d", volume)
		if err := p.sendCommand(cmd); err != nil {
			return fmt.Errorf("failed to send volume command: %w", err)
		}
	}

	return nil
}

// GetVolume returns the current volume.
func (p *RemotePlayer) GetVolume() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.volume
}

// Cleanup stops playback and cleans up resources.
func (p *RemotePlayer) Cleanup() error {
	return p.Stop()
}

// sendCommand sends a command to the client wrapper using OSC sequences.
// Format: \033]8888;COMMAND\007
func (p *RemotePlayer) sendCommand(cmd string) error {
	if p.writer == nil {
		return fmt.Errorf("no writer available")
	}

	msg := oscPrefix + cmd + oscSuffix
	_, err := p.writer.Write([]byte(msg))
	return err
}
