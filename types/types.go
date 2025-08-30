// Package types provides shared data structures used across the password manager application.
// These types facilitate communication between different packages and manage application state.
package types

// Options represents the current state and configuration of the application session.
// It tracks login status, quit requests, and error messages for the user interface.
type Options struct {
	// LoggedIn indicates whether the user has successfully authenticated
	LoggedIn bool
	
	// Quit signals that the user wants to exit the application (via Ctrl+C or Esc)
	Quit bool
	
	// ErrorMessage holds validation or authentication error messages to display to the user
	ErrorMessage string
}