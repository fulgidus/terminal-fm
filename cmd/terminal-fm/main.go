// Package main implements the Terminal.FM local TUI application.
package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fulgidus/terminal-fm/internal/config"
	"github.com/fulgidus/terminal-fm/pkg/i18n"
	"github.com/fulgidus/terminal-fm/pkg/services/player"
	"github.com/fulgidus/terminal-fm/pkg/services/radiobrowser"
	"github.com/fulgidus/terminal-fm/pkg/services/storage"
	"github.com/fulgidus/terminal-fm/pkg/ui"
)

var (
	version = "1.0.0"
	devMode = flag.Bool("dev", false, "Run in development mode with mock data")
	locale  = flag.String("locale", "", "Set locale (en or it)")
	showVer = flag.Bool("version", false, "Show version information")
)

func main() {
	flag.Parse()

	// Show version and exit
	if *showVer {
		fmt.Printf("Terminal.FM v%s\n", version)
		os.Exit(0)
	}

	// Load configuration
	cfg := config.New()
	cfg.DevMode = *devMode

	// Set locale if provided
	if *locale != "" {
		cfg.I18n.DefaultLocale = *locale
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}

	// Ensure data directory exists
	if err := cfg.EnsureDataDir(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create data directory: %v\n", err)
		os.Exit(1)
	}

	// Initialize Radio Browser API client
	var radioClient radiobrowser.Client
	var err error

	if cfg.DevMode {
		// Use mock client in development mode
		radioClient = radiobrowser.NewMockClient()
	} else {
		// Use real API client
		radioClient, err = radiobrowser.NewAPIClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize Radio Browser API: %v\n", err)
			os.Exit(1)
		}
	}

	// Initialize storage
	store, err := storage.NewStore(cfg.Storage.DBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize storage: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	// Initialize audio player
	var audioPlayer player.Player
	if cfg.Player.DefaultPlayer == "mpv" {
		// TODO: Implement MpvPlayer when ready
		audioPlayer = player.NewFFplayPlayer(cfg.Player.FFplayPath)
	} else {
		audioPlayer = player.NewFFplayPlayer(cfg.Player.FFplayPath)
	}

	// Create the TUI model
	model := ui.NewModel(radioClient, audioPlayer, store, cfg.I18n.DefaultLocale)

	// Initialize translator for startup messages
	tr := i18n.NewSimpleTranslator(cfg.I18n.DefaultLocale)

	// Create Bubbletea program
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the program
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}

	// Clean up resources
	if m, ok := finalModel.(ui.Model); ok {
		m.Cleanup()
	}

	// Show goodbye message
	fmt.Println(tr.T("goodbye"))
}
