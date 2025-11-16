// Package player provides audio playback functionality.
package player

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"
	"syscall"

	"github.com/fulgidus/terminal-fm/pkg/services/radiobrowser"
)

// StreamingPlayer implements Player by streaming audio data through an io.Writer.
// This allows audio to be sent through SSH sessions or other transports.
type StreamingPlayer struct {
	mu             sync.RWMutex
	cmd            *exec.Cmd
	state          State
	currentStation *radiobrowser.Station
	volume         int
	writer         io.Writer
	ffmpegPath     string
	processActive  bool
}

// NewStreamingPlayer creates a new streaming player that writes PCM audio to the given writer.
func NewStreamingPlayer(writer io.Writer, ffmpegPath string) *StreamingPlayer {
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}

	return &StreamingPlayer{
		state:      StateStopped,
		volume:     70,
		writer:     writer,
		ffmpegPath: ffmpegPath,
	}
}

// Play starts playing a radio station by streaming audio to the writer.
func (p *StreamingPlayer) Play(station *radiobrowser.Station) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Stop any current playback
	if err := p.stopLocked(); err != nil {
		return fmt.Errorf("failed to stop current playback: %w", err)
	}

	if station == nil || station.URLResolved == "" {
		return fmt.Errorf("invalid station or URL")
	}

	if p.writer == nil {
		return fmt.Errorf("no output writer configured")
	}

	// Use ffmpeg to transcode stream to PCM audio (CD quality: 44.1kHz, 16-bit, stereo)
	// -i: input URL
	// -f s16le: output format (signed 16-bit little-endian PCM)
	// -ar 44100: sample rate
	// -ac 2: stereo (2 channels)
	// -af volume: adjust volume
	// pipe:1: output to stdout
	volumeFilter := fmt.Sprintf("volume=%.2f", float64(p.volume)/100.0)
	args := []string{
		"-i", station.URLResolved,
		"-f", "s16le",
		"-ar", "44100",
		"-ac", "2",
		"-af", volumeFilter,
		"-loglevel", "error", // Only show errors
		"pipe:1",
	}

	p.cmd = exec.Command(p.ffmpegPath, args...)
	p.cmd.Stdout = p.writer
	p.cmd.Stderr = nil // Discard stderr to avoid polluting TUI

	// Start the transcoding process
	if err := p.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	p.state = StatePlaying
	p.currentStation = station
	p.processActive = true

	log.Printf("Started streaming: %s (%s)", station.Name, station.URLResolved)

	// Monitor process in background
	cmd := p.cmd
	go func() {
		err := cmd.Wait()

		p.mu.Lock()
		defer p.mu.Unlock()

		// Only clean up if this is still our active command
		if p.cmd == cmd {
			if err != nil {
				log.Printf("Streaming ended with error: %v", err)
			} else {
				log.Printf("Streaming ended normally")
			}

			p.processActive = false
			p.state = StateStopped
			p.currentStation = nil
			p.cmd = nil
		}
	}()

	return nil
}

// Stop stops the current playback.
func (p *StreamingPlayer) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.stopLocked()
}

// stopLocked stops playback without acquiring the lock (internal use).
func (p *StreamingPlayer) stopLocked() error {
	wasActive := p.processActive
	p.processActive = false
	p.state = StateStopped

	if wasActive && p.cmd != nil && p.cmd.Process != nil {
		// Send SIGTERM first (graceful)
		if err := p.cmd.Process.Signal(syscall.SIGTERM); err != nil {
			// If SIGTERM fails, try SIGKILL
			p.cmd.Process.Signal(syscall.SIGKILL)
		}
	}

	p.currentStation = nil
	return nil
}

// GetState returns the current playback state.
func (p *StreamingPlayer) GetState() State {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state
}

// GetCurrentStation returns the currently playing station.
func (p *StreamingPlayer) GetCurrentStation() *radiobrowser.Station {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.currentStation
}

// SetVolume sets the playback volume (0-100).
// Note: Requires restarting playback to take effect.
func (p *StreamingPlayer) SetVolume(volume int) error {
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
		p.mu.Unlock()
		err := p.Play(station)
		p.mu.Lock()
		return err
	}

	return nil
}

// GetVolume returns the current volume.
func (p *StreamingPlayer) GetVolume() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.volume
}

// Cleanup forcefully stops playback and cleans up resources.
func (p *StreamingPlayer) Cleanup() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Signal(syscall.SIGKILL)
	}

	p.cmd = nil
	p.state = StateStopped
	p.currentStation = nil
	p.processActive = false

	return nil
}
