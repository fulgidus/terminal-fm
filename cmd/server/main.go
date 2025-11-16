// Package main is the entry point for Terminal.FM server.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fulgidus/terminal-fm/internal/config"
	"github.com/fulgidus/terminal-fm/pkg/services/radiobrowser"
	sshserver "github.com/fulgidus/terminal-fm/pkg/ssh"
)

var (
	// Version information (set by build flags)
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Parse command-line flags
	var (
		showVersion = flag.Bool("version", false, "Show version information")
		devMode     = flag.Bool("dev", false, "Run in development mode (port 2222)")
		port        = flag.Int("port", 0, "SSH port to listen on (overrides config)")
		host        = flag.String("host", "", "Host address to bind to (overrides config)")
	)
	flag.Parse()

	// Show version and exit
	if *showVersion {
		fmt.Printf("Terminal.FM %s (commit: %s, built: %s)\n", version, commit, date)
		os.Exit(0)
	}

	// Load configuration
	cfg := config.New()
	cfg.DevMode = *devMode

	// Apply command-line overrides
	if *devMode {
		cfg.Server.Port = 2222
		log.Println("Running in development mode")
	}
	if *port > 0 {
		cfg.Server.Port = *port
	}
	if *host != "" {
		cfg.Server.Host = *host
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Ensure data directory exists
	if err := cfg.EnsureDataDir(); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Initialize radio browser client (using mock for now)
	radioClient := radiobrowser.NewMockClient()
	log.Println("Using mock Radio Browser API client")

	// Create SSH server
	srv, err := sshserver.NewServer(cfg, radioClient)
	if err != nil {
		log.Fatalf("Failed to create SSH server: %v", err)
	}

	// Start server (blocks until shutdown signal)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
