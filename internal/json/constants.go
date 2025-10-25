package json

import (
	"log/slog"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
	gonethttphandlerjsonjsend "github.com/ralvarezdev/go-net/http/handler/json/jsend"
)

var (
	// Handler is the JSON handler in JSend format
	Handler gonethttphandler.Handler
)

// Load initializes the JSON encoder and decoder
//
// Parameters:
//
//   - mode: the go-flags mode flag to determine if the environment is in debug mode
//   - logger: the logger instance
func Load(mode *goflagsmode.Flag, logger *slog.Logger) {
	// Initialize the handler
	handler, err := gonethttphandlerjsonjsend.NewHandler(
		mode,
		logger,
	)
	if err != nil {
		panic(err)
	}
	Handler = handler
}
