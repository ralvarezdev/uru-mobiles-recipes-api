package loader

import (
	"log/slog"

	"github.com/joho/godotenv"
	goloaderenv "github.com/ralvarezdev/go-loader/env"
)

var (
	// Loader is the environment variables loader
	Loader goloaderenv.Loader

	// Logger is the logger for the loader
	Logger *slog.Logger
)

// Load loads the loader
//
// Parameters:
//
//   - logger: The logger (optional, can be nil)
func Load(logger *slog.Logger) {
	if logger != nil {
		Logger = logger.With(
			slog.String("component", "loader_env"),
		)
	}

	// Load the environment variables loader
	loader, err := goloaderenv.NewDefaultLoader(
		func() error {
			return godotenv.Load()
		},
		Logger,
	)
	if err != nil {
		panic(err)
	}
	Loader = loader
}
