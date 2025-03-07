package types

// Constants used across the application
const (
	// PreviousContextSymbol represents the symbol used to switch to the previous context
	PreviousContextSymbol = "-"
)

// Config represents the application configuration
type Config struct {
	// ConfigDir is the directory containing talosconfig files
	ConfigDir string
	// LogLevel is the logging level
	LogLevel string
	// LogFormat is the logging format (json, text)
	LogFormat string
}
