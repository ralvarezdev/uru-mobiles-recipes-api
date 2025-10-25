package auth

import (
	gogrpc "github.com/ralvarezdev/go-grpc"
	gonethttpcookie "github.com/ralvarezdev/go-net/http/cookie"
	gonethttpgrpc "github.com/ralvarezdev/go-net/http/grpc"
	pbauthcompiled "github.com/ralvarezdev/grpc-auth-proto-go/compiled/ralvarezdev/auth"

	internalcookie "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/cookie"
	internalloader "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/loader"
)

const (
	// EnvGRPCAuthHost is the environment variable for the gRPC auth host
	EnvGRPCAuthHost = "GRPC_AUTH_HOST"

	// EnvGRPCAuthPort is the environment variable for the gRPC auth port
	EnvGRPCAuthPort = "GRPC_AUTH_PORT"

	// EnvGRPCAuthAPIKey is the environment variable for the gRPC auth API key
	EnvGRPCAuthAPIKey = "GRPC_AUTH_API_KEY"
)

var (
	// Client is the gRPC auth client
	Client pbauthcompiled.AuthClient

	// GRPCAuthHost is the gRPC auth host
	GRPCAuthHost string

	// GRPCAuthPort is the gRPC auth port
	GRPCAuthPort int

	// GRPCAuthAPIKey is the gRPC auth API key
	GRPCAuthAPIKey string

	// AuthenticationParser is the gRPC auth authentication parser
	AuthenticationParser gonethttpgrpc.AuthenticationParser
)

// Load loads the auth gRPC client
func Load() {
	// Load the gRPC host and API key
	for env, dest := range map[string]*string{
		EnvGRPCAuthHost:   &GRPCAuthHost,
		EnvGRPCAuthAPIKey: &GRPCAuthAPIKey,
	} {
		if err := internalloader.Loader.LoadVariable(
			env,
			dest,
		); err != nil {
			panic(err)
		}
	}

	// Load the gRPC port
	if err := internalloader.Loader.LoadIntVariable(
		EnvGRPCAuthPort,
		&GRPCAuthPort,
	); err != nil {
		panic(err)
	}

	// Initialize the authentication parser options
	options := gonethttpgrpc.Options{
		MetadataKeysToCookiesAttributes: map[string]*gonethttpcookie.Attributes{
			gogrpc.AccessTokenMetadataKey:  internalcookie.AccessToken,
			gogrpc.RefreshTokenMetadataKey: internalcookie.RefreshToken,
		},
	}

	// Create the authentication parser
	authenticationParser, err := gonethttpgrpc.NewDefaultAuthenticationParser(
		&options,
	)
	if err != nil {
		panic(err)
	}
	AuthenticationParser = authenticationParser
}
