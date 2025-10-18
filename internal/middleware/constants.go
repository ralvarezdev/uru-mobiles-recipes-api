package middleware

import (
	"context"
	"log/slog"
	"net/http"

	gogrpcnethttp "github.com/ralvarezdev/go-grpc/client/net/http"
	gogrpcmd "github.com/ralvarezdev/go-grpc/metadata"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gojwttokenvalidor "github.com/ralvarezdev/go-jwt/token/validator"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
	gonethttpmiddlewareauth "github.com/ralvarezdev/go-net/http/middleware/auth"
	gonethttpmiddlewareauthgrpc "github.com/ralvarezdev/go-net/http/middleware/auth/grpc"
	gonethttpmiddlewareerrorhandler "github.com/ralvarezdev/go-net/http/middleware/error_handler"
	gonethttpmiddlewareratelimiter "github.com/ralvarezdev/go-net/http/middleware/rate_limiter/redis"
	gonethttpmiddlewaresizelimiter "github.com/ralvarezdev/go-net/http/middleware/size_limiter"
	gonethttpmiddlewarevalidator "github.com/ralvarezdev/go-net/http/middleware/validator"
	gonethttpresponsejsendgrpc "github.com/ralvarezdev/go-net/http/response/jsend/grpc"
	goratelimiter "github.com/ralvarezdev/go-rate-limiter/redis"
	pbauth "github.com/ralvarezdev/grpc-auth-proto-go"
	internalgrpcauth "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/grpc/auth"
	internalloader "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/loader"
	pbempty "google.golang.org/protobuf/types/known/emptypb"
)

const (
	// EnvBodyLimit is the environment variable key for the body limit
	EnvBodyLimit = "BODY_LIMIT"
)

var (
	// JWTOptions is the JWT options
	JWTOptions gonethttpmiddlewareauth.Options

	// BodyLimit is the API body limit
	BodyLimit int

	// HandleError is the API error handler middleware function
	HandleError func(next http.Handler) http.Handler

	// LimitBody is the API body limit middleware function
	LimitBody func(next http.Handler) http.Handler

	// ValidateJSON is the API request validator middleware function for JSON requests
	ValidateJSON func(
		body interface{},
		auxiliaryValidatorFns ...interface{},
	) func(next http.Handler) http.Handler

	// ValidateProtoJSON is the API request validator middleware function for ProtoJSON requests
	ValidateProtoJSON func(
		body interface{},
		auxiliaryValidatorFns ...interface{},
	) func(next http.Handler) http.Handler

	// LimitRequests is the API rate limiter middleware function
	LimitRequests func(next http.Handler) http.Handler

	// Authenticate is the JWT authentication middleware function
	Authenticate func(
		method string,
	) func(next http.Handler) http.Handler
)

// RefreshToken is the function to refresh JWT tokens
//
// Parameters:
//
//   - w: The HTTP response writer
//   - r: The HTTP request
//
// Returns:
//
//   - A map of tokens
//   - An error, if any
func RefreshToken(
	w http.ResponseWriter,
	r *http.Request,
) (map[gojwttoken.Token]string, error) {
	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return nil, err
	}

	// Call the gRPC service to refresh the token
	if _, err = internalgrpcauth.Client.RefreshToken(
		ctx,
		&pbempty.Empty{},
	); err != nil {
		return nil, gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Get the refreshed tokens from the context
	refreshToken, err := gogrpcmd.GetCtxMetadataRefreshToken(ctx)
	if err != nil {
		return nil, err
	}
	accessToken, err := gogrpcmd.GetCtxMetadataAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	// Create the tokens map
	tokens := make(map[gojwttoken.Token]string)
	tokens[gojwttoken.RefreshToken] = refreshToken
	tokens[gojwttoken.AccessToken] = accessToken
	return tokens, nil
}

// Load loads the API middlewares
//
// Parameters:
//
//   - jsonHandler: The JSON handler
//   - protoJSONHandler: The ProtoJSON handler
//   - rateLimiter: The rate limiter
//   - jwtValidator: The JWT validator
//   - logger: The logger
func Load(
	jsonHandler gonethttphandler.Handler,
	protoJSONHandler gonethttphandler.Handler,
	rateLimiter goratelimiter.RateLimiter,
	jwtValidator gojwttokenvalidor.Validator,
	logger *slog.Logger,
) {
	// Load the body limit
	if err := internalloader.Loader.LoadIntVariable(
		EnvBodyLimit,
		&BodyLimit,
	); err != nil {
		panic(err)
	}

	// Create API error handler middleware
	errorHandler, err := gonethttpmiddlewareerrorhandler.NewMiddleware(
		jsonHandler,
	)
	if err != nil {
		panic(err)
	}
	HandleError = errorHandler.HandleError

	// Create API body limit middleware
	sizeLimiter := gonethttpmiddlewaresizelimiter.NewMiddleware()
	LimitBody = sizeLimiter.Limit(int64(BodyLimit))

	// Create API request validator middleware for JSON requests
	jsonValidator, err := gonethttpmiddlewarevalidator.NewMiddleware(
		jsonHandler,
		nil,
		nil,
		logger,
	)
	if err != nil {
		panic(err)
	}
	ValidateJSON = jsonValidator.Validate

	// Create API request validator middleware for ProtoJSON requests
	protoJSONValidator, err := gonethttpmiddlewarevalidator.NewMiddleware(
		protoJSONHandler,
		nil,
		nil,
		logger,
	)
	if err != nil {
		panic(err)
	}
	ValidateProtoJSON = protoJSONValidator.Validate

	// Create API rate limiter middleware
	rateLimiterMiddleware, err := gonethttpmiddlewareratelimiter.NewMiddleware(
		jsonHandler,
		rateLimiter,
		logger,
	)
	if err != nil {
		panic(err)
	}
	LimitRequests = rateLimiterMiddleware.Limit()

	// Initialize JWT options
	cookieRefreshTokenName := gojwttoken.RefreshToken.String()
	cookieAccessTokenName := gojwttoken.AccessToken.String()
	JWTOptions = gonethttpmiddlewareauth.Options{
		CookieRefreshTokenName: &cookieRefreshTokenName,
		CookieAccessTokenName:  &cookieAccessTokenName,
		RefreshTokenFn:         RefreshToken,
	}

	// Create JWT authentication middleware
	authenticator, err := gonethttpmiddlewareauth.NewMiddleware(
		jsonHandler,
		jwtValidator,
		&JWTOptions,
	)
	if err != nil {
		panic(err)
	}

	// Create JWT authentication middleware
	grpcAuthenticator, err := gonethttpmiddlewareauthgrpc.NewMiddleware(
		pbauth.JWTInterceptions,
		jsonHandler,
		authenticator,
		logger,
	)
	if err != nil {
		panic(err)
	}
	Authenticate = grpcAuthenticator.AuthenticateFromCookie
}
