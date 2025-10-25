package jwt

import (
	"log/slog"
	"strings"

	godatabasessql "github.com/ralvarezdev/go-databases/sql"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gojwttokenclaims "github.com/ralvarezdev/go-jwt/token/claims"
	gojwttokenclaimssqlite "github.com/ralvarezdev/go-jwt/token/claims/sqlite"
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

	// TokenValidator is the JWT token validator
	TokenValidator gojwttokenclaims.TokenValidator

	// ClaimsValidator is the JWT claims validator instance
	ClaimsValidator gojwttokenclaims.ClaimsValidator

	// Validator is the JWT validator instance
	Validator gojwttokenvalidator.Validator
)

// Load initializes the JWT validator
//
// Parameters:
//
//   - mode: the go-flags mode flag to determine if the environment is in debug mode
//   - tokenValidatorSQLiteService: the SQLite service for the token validator
//   - logger: the logger instance
func Load(
	mode *goflagsmode.Flag,
	tokenValidatorSQLiteService godatabasessql.Service,
	logger *slog.Logger,
) {
	// Load the JWT public key from environment variable
	if err := internalloader.Loader.LoadVariable(
		EnvJWTPublicKey,
		&JWTPublicKey,
	); err != nil {
		panic(err)
	}
	JWTPublicKey = strings.ReplaceAll(JWTPublicKey, `\n`, "\n")

	// Initialize the JWT token validator with SQLite
	tokenValidator, err := gojwttokenclaimssqlite.NewTokenValidator(
		tokenValidatorSQLiteService,
		logger,
	)
	if err != nil {
		panic(err)
	}
	TokenValidator = tokenValidator

	// Initialize the JWT claims validator
	claimsValidator, err := gojwttokenclaims.NewDefaultClaimsValidator(
		TokenValidator,
	)
	if err != nil {
		panic(err)
	}
	ClaimsValidator = claimsValidator

	// Initialize the JWT validator
	validator, err := gojwttokenvalidator.NewEd25519Validator(
		[]byte(JWTPublicKey),
		ClaimsValidator,
		mode,
	)
	if err != nil {
		panic(err)
	}
	Validator = validator
}
