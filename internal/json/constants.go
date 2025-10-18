package json

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
	gonethttphandlerjsonjsend "github.com/ralvarezdev/go-net/http/handler/json/jsend"
)

var (
	// Handler is the JSON handler in JSend format
	Handler gonethttphandler.Handler
)

// Load initializes the JSON encoder and decoder
func Load(mode *goflagsmode.Flag) {
	// Initialize the handler
	handler, err := gonethttphandlerjsonjsend.NewHandler(
		mode,
	)
	if err != nil {
		panic(err)
	}
	Handler = handler
}
