package protojson

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
	gonethttphandlerprotojsonjsend "github.com/ralvarezdev/go-net/http/handler/protojson/jsend"
)

var (
	// Handler is the ProtoJSON handler in JSend format
	Handler gonethttphandler.Handler
)

// Load initializes the JSON encoder and decoder
//
// Parameters:
//
//   - mode: the go-flags mode flag to determine if the environment is in debug mode
func Load(mode *goflagsmode.Flag) {
	// Initialize the handler
	handler, err := gonethttphandlerprotojsonjsend.NewHandler(
		mode,
	)
	if err != nil {
		panic(err)
	}
	Handler = handler
}
