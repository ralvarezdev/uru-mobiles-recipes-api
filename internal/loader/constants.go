package loader

import (
	"log/slog"

	"github.com/joho/godotenv"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
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
//   - mode: The application mode
//   - logger: The logger (optional, can be nil)
func Load(mode *goflagsmode.Flag, logger *slog.Logger) {
	if logger != nil {
		Logger = logger.With(
			slog.String("component", "loader_env"),
		)
	}

	// Load the environment variables loader
	loader, err := goloaderenv.NewDefaultLoader(
		func() error {
			// Load .env file only if not in production mode
			if mode != nil && mode.IsProd() {
				return nil
			}
			return godotenv.Load()
		},
		Logger,
	)
	if err != nil {
		panic(err)
	}
	Loader = loader
}
