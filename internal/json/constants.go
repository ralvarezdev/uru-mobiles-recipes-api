package json

import (
	"log/slog"
	"path/filepath"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	goloaderfilesystem "github.com/ralvarezdev/go-loader/filesystem"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
	gonethttphandlerjsonjsend "github.com/ralvarezdev/go-net/http/handler/json/jsend"
)

const (
	// SwaggerJSONDefinitionsRelativePath is the relative path for the swagger.json definitions
	SwaggerJSONDefinitionsRelativePath = "docs/swagger.json"
)

var (
	// Handler is the JSON handler in JSend format
	Handler gonethttphandler.Handler

	// SwaggerJSONDefinitionsFilePath is the file path for the swagger.json definitions
	SwaggerJSONDefinitionsFilePath string

	// SwaggerJSONDefinitions is the swagger.json definitions content
	SwaggerJSONDefinitions []byte
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

	// Load swagger.json definitions
	goModPath, err := goloaderfilesystem.GetExecutableGoModPath()
	if err != nil {
		panic(err)
	}

	// Get the swagger.json definitions file path based on the go.mod path
	SwaggerJSONDefinitionsFilePath = filepath.Join(
		filepath.Dir(goModPath),
		SwaggerJSONDefinitionsRelativePath,
	)

	// Load the swagger.json definitions content
	fileContent, err := goloaderfilesystem.ReadFile(SwaggerJSONDefinitionsFilePath)
	if err != nil {
		panic(err)
	}
	SwaggerJSONDefinitions = fileContent
}
