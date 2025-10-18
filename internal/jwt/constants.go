package jwt

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gojwttokenclaims "github.com/ralvarezdev/go-jwt/token/claims"
	gojwttokenvalidator "github.com/ralvarezdev/go-jwt/token/validator"
	internalloader "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/loader"
)

const (
	// EnvJWTPublicKey is the environment variable name for the JWT public key
	EnvJWTPublicKey = "GRPC_AUTH_JWT_PUBLIC_KEY"
)

var (
	// JWTPublicKey is the JWT public key
	JWTPublicKey string

	// Validator is the JWT validator instance
	Validator gojwttokenvalidator.Validator
)

// Load initializes the JWT validator
//
// Parameters:
//
//   - mode: the go-flags mode flag to determine if the environment is in debug mode
//   - claimsValidator: the JWT token claims validator
func Load(mode *goflagsmode.Flag, claimsValidator gojwttokenclaims.Validator) {
	// Load the JWT public key from environment variable
	if err := internalloader.Loader.LoadVariable(
		EnvJWTPublicKey,
		&JWTPublicKey,
	); err != nil {
		panic(err)
	}

	// Initialize the JWT validator
	validator, err := gojwttokenvalidator.NewEd25519Validator(
		[]byte(JWTPublicKey),
		claimsValidator,
		mode,
	)
	if err != nil {
		panic(err)
	}
	Validator = validator
}
