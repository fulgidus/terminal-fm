// Package ssh provides SSH server functionality for Terminal.FM.
package ssh

import (
	"github.com/charmbracelet/ssh"
)

// NoPasswordAuth is a password handler that accepts any username/password.
func NoPasswordAuth(ctx ssh.Context, password string) bool {
	// Accept any password (including empty) for anonymous access
	return true
}

// NoPublicKeyAuth is a public key handler that accepts any key.
func NoPublicKeyAuth(ctx ssh.Context, key ssh.PublicKey) bool {
	// Accept any public key for anonymous access
	return true
}
