package ssh

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/fulgidus/terminal-fm/internal/config"
	"github.com/fulgidus/terminal-fm/pkg/services/radiobrowser"
	"github.com/fulgidus/terminal-fm/pkg/ui"
	
	tea "github.com/charmbracelet/bubbletea"
)

// Server wraps the SSH server with configuration.
type Server struct {
	config      *config.Config
	sshServer   *ssh.Server
	radioClient radiobrowser.Client
}

// NewServer creates a new SSH server instance.
func NewServer(cfg *config.Config, radioClient radiobrowser.Client) (*Server, error) {
	s := &Server{
		config:      cfg,
		radioClient: radioClient,
	}
	
	// Create the Bubbletea handler that will be called for each SSH session
	teaHandler := func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
		// Get locale from environment or use default
		locale := cfg.I18n.DefaultLocale
		
		// Check environment variables for locale
		for _, env := range sess.Environ() {
			// Environment comes as "KEY=value" format
			if len(env) >= 5 && env[:5] == "LANG=" {
				envLang := env[5:]
				// Simple locale detection - just check if it starts with "it"
				if len(envLang) >= 2 && envLang[:2] == "it" {
					locale = "it"
				}
				break
			}
		}
		
		// Create a new UI model for this session
		model := ui.NewModel(radioClient, locale)
		
		return model, []tea.ProgramOption{
			tea.WithAltScreen(),       // Use alternate screen buffer
			tea.WithMouseCellMotion(), // Enable mouse support
		}
	}
	
	// Build the SSH server with middleware stack
	srv, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)),
		wish.WithHostKeyPath(cfg.Server.KeyPath),
		wish.WithPublicKeyAuth(NoPublicKeyAuth),
		wish.WithPasswordAuth(NoPasswordAuth),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			logging.Middleware(),
		),
		wish.WithIdleTimeout(time.Duration(cfg.Server.IdleTimeout)*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH server: %w", err)
	}
	
	s.sshServer = srv
	return s, nil
}

// Start begins listening for SSH connections.
func (s *Server) Start() error {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	
	log.Printf("Starting SSH server on %s:%d", s.config.Server.Host, s.config.Server.Port)
	
	// Start server in goroutine
	go func() {
		if err := s.sshServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	
	// Wait for shutdown signal
	<-done
	log.Println("Shutting down SSH server...")
	
	return s.Shutdown()
}

// Shutdown gracefully stops the SSH server.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := s.sshServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown SSH server: %w", err)
	}
	
	log.Println("SSH server stopped")
	return nil
}
