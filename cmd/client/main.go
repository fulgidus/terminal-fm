// Package main implements the terminal-radio client wrapper.
// This wrapper connects to the SSH server and intercepts OSC commands
// to control local audio playback.
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

const (
	defaultServer = "terminal-radio.com"
	oscPrefix     = "\033]8888;"
	oscSuffix     = "\007"
)

var (
	playerCmd *exec.Cmd
)

func main() {
	// Get server address from args or use default
	server := defaultServer
	sshArgs := []string{}

	if len(os.Args) > 1 {
		server = os.Args[1]
		// Check if it includes port (e.g., localhost:2222)
		if strings.Contains(server, ":") {
			parts := strings.Split(server, ":")
			server = parts[0]
			if len(parts) == 2 {
				sshArgs = append(sshArgs, "-p", parts[1])
			}
		}
	}

	sshArgs = append(sshArgs, server)

	// Setup signal handling for cleanup
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		cleanup()
		os.Exit(0)
	}()

	// Start SSH connection
	sshCmd := exec.Command("ssh", sshArgs...)
	sshCmd.Stdin = os.Stdin
	sshCmd.Stderr = os.Stderr

	// Get stdout pipe to intercept OSC commands
	stdout, err := sshCmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout pipe: %v", err)
	}

	// Start SSH
	if err := sshCmd.Start(); err != nil {
		log.Fatalf("Failed to start SSH: %v", err)
	}

	// Process output: intercept OSC commands and pass through the rest
	go processOutput(stdout)

	// Wait for SSH to finish
	if err := sshCmd.Wait(); err != nil {
		// Don't log error if it's just a normal exit
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != 0 {
				log.Printf("SSH exited with code %d", exitErr.ExitCode())
			}
		}
	}

	// Cleanup player on exit
	cleanup()
}

// processOutput reads from SSH stdout, intercepts OSC commands, and forwards the rest
func processOutput(r io.Reader) {
	reader := bufio.NewReader(r)
	var buffer strings.Builder

	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading output: %v", err)
			}
			break
		}

		// Write byte to stdout
		os.Stdout.Write([]byte{b})

		// Check for OSC sequence start
		buffer.WriteByte(b)
		content := buffer.String()

		// Keep last 100 bytes in buffer to detect sequences
		if buffer.Len() > 100 {
			content = content[buffer.Len()-100:]
			buffer.Reset()
			buffer.WriteString(content)
		}

		// Check if we have a complete OSC command
		if strings.Contains(content, oscPrefix) && strings.Contains(content, oscSuffix) {
			start := strings.Index(content, oscPrefix)
			end := strings.Index(content[start:], oscSuffix)
			if end != -1 {
				end += start
				cmd := content[start+len(oscPrefix) : end]
				handleCommand(cmd)
				// Clear buffer after processing command
				buffer.Reset()
			}
		}
	}
}

// handleCommand processes OSC commands from the server
func handleCommand(cmd string) {
	parts := strings.Split(cmd, ";")
	if len(parts) == 0 {
		return
	}

	action := parts[0]

	switch action {
	case "PLAY":
		if len(parts) >= 2 {
			url := parts[1]
			volume := "70"
			if len(parts) >= 3 {
				volume = parts[2]
			}
			playAudio(url, volume)
		}

	case "STOP":
		stopAudio()

	case "VOLUME":
		if len(parts) >= 2 {
			volume := parts[1]
			setVolume(volume)
		}
	}
}

// playAudio starts playing audio using mpv or ffplay
func playAudio(url, volume string) {
	// Stop any existing player
	stopAudio()

	// Try mpv first (most common and feature-rich)
	if isCommandAvailable("mpv") {
		playerCmd = exec.Command("mpv",
			"--no-video",
			"--really-quiet",
			fmt.Sprintf("--volume=%s", volume),
			url,
		)
	} else if isCommandAvailable("ffplay") {
		// Fallback to ffplay
		playerCmd = exec.Command("ffplay",
			"-nodisp",
			"-loglevel", "quiet",
			"-autoexit",
			"-volume", volume,
			url,
		)
	} else if isCommandAvailable("vlc") {
		// Last resort: VLC
		playerCmd = exec.Command("vlc",
			"--intf", "dummy",
			"--play-and-exit",
			"--volume", volume,
			url,
		)
	} else {
		log.Println("Warning: No audio player found (mpv, ffplay, or vlc). Please install one.")
		return
	}

	// Start player in background
	if err := playerCmd.Start(); err != nil {
		log.Printf("Failed to start player: %v", err)
		playerCmd = nil
	}
}

// stopAudio stops the current player
func stopAudio() {
	if playerCmd != nil && playerCmd.Process != nil {
		if err := playerCmd.Process.Signal(syscall.SIGTERM); err != nil {
			// Process might already be dead, try SIGKILL
			_ = playerCmd.Process.Kill()
		}
		_ = playerCmd.Wait() // Ignore wait errors
		playerCmd = nil
	}
}

// setVolume adjusts volume (requires restarting playback with current implementation)
func setVolume(volume string) {
	// For now, volume changes require restart
	// This is handled by the server sending STOP then PLAY with new volume
	_ = volume
}

// isCommandAvailable checks if a command exists in PATH
func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// cleanup stops audio on exit
func cleanup() {
	stopAudio()
}
