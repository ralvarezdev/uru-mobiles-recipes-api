package logger

import (
	"log/slog"
	"os"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
)

var (
	// Options are the options for the logger
	Options *slog.HandlerOptions

	// Logger is the logger for the application
	Logger *slog.Logger
)

// Load loads the constants from the environment variables
//
// Parameters:
//
//   - mode: the go-flags mode flag to determine if the environment is in debug mode
func Load(mode *goflagsmode.Flag) {
	// Check if the environment is in debug mode
	if mode != nil && mode.IsDebug() {
		Options = &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}
	} else {
		Options = &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true,
		}
	}

	// Create a new logger
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, Options))
}
