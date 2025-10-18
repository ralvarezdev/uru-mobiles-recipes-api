package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"path/filepath"
	"time"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gogrpcclientinterceptorauthapikey "github.com/ralvarezdev/go-grpc/client/interceptor/auth/api_key"
	gogrpcclientinterceptorauthjwt "github.com/ralvarezdev/go-grpc/client/interceptor/auth_verifier/jwt"
	gonetflagsport "github.com/ralvarezdev/go-net/flags/port"
	gonethttproute "github.com/ralvarezdev/go-net/http/route"
	pbauth "github.com/ralvarezdev/grpc-auth-proto-go"
	pbauthcompiled "github.com/ralvarezdev/grpc-auth-proto-go/compiled/ralvarezdev/auth"
	internalcookie "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/cookie"
	internalredis "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/databases/redis"
	internalsqlite "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/databases/sqlite"
	internalgrpcauth "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/grpc/auth"
	internaljson "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/json"
	internaljwt "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/jwt"
	internalloader "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/loader"
	internallogger "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/logger"
	internalmiddleware "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/middleware"
	internalprotojson "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/protojson"
	internalrabbitmq "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/rabbitmq"
	internalrouter "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/router"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	// ModeFlag is the mode flag
	ModeFlag = goflagsmode.NewFlag(
		goflagsmode.Dev,
		goflagsmode.AllowedModes,
	)

	// PortFlag is the port flag
	PortFlag = gonetflagsport.NewFlag(
		nil,
	)

	// Port is the port to listen on
	Port int
)

// init initializes the flags and calls the load functions
func init() {
	// Define the mode and port flags
	goflagsmode.SetFlag(ModeFlag)
	gonetflagsport.SetFlag(PortFlag)

	// Parse the flags
	flag.Parse()

	// Set the port variable
	port, err := PortFlag.Port()
	if err != nil {
		panic(err)
	}
	Port = port

	// Call the load functions
	internallogger.Load(ModeFlag)
	internalloader.Load(internallogger.Logger)
	internalcookie.Load(ModeFlag)
	internaljson.Load(ModeFlag)
	internalprotojson.Load(ModeFlag)
	internalredis.Load()
	internalsqlite.Load(internallogger.Logger)
	internaljwt.Load(
		ModeFlag,
		internalsqlite.TokenValidatorHandler,
		internallogger.Logger,
	)
	internalrabbitmq.Load(
		internaljwt.TokenValidator,
		internallogger.Logger,
	)
	internalmiddleware.Load(
		internaljson.Handler,
		internalprotojson.Handler,
		internalredis.RateLimiter,
		internaljwt.Validator,
		internallogger.Logger,
	)
	internalgrpcauth.Load()
}

//	@Title			... REST API
//	@Version		1.0
//	@Description	This is the REST API for the ... application.

//	@License.name	GPL-3.0
//	@License.url	http://www.gnu.org/licenses/gpl-3.0.html

//	@BasePath	/

// @securityDefinitions.apikey	CookieAuth
// @in							cookie
// @name						access_token
func main() {
	// Create a global context
	ctx := context.Background()
	defer ctx.Done()

	// Listen on the given port
	portListener, err := net.Listen(
		"tcp",
		fmt.Sprintf(":%d", Port),
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = portListener.Close(); err != nil {
			panic(err)
		}
	}()

	// Create the auth client JWT authentication interceptor
	authJWTInterceptor, err := gogrpcclientinterceptorauthjwt.NewInterceptor(
		pbauth.JWTInterceptions,
		nil,
		internallogger.Logger,
	)
	if err != nil {
		panic(err)
	}

	// Create the auth client API key authentication interceptor
	apiKeyInterceptor, err := gogrpcclientinterceptorauthapikey.NewInterceptor(
		pbauth.APIKeysInterceptions,
		internalgrpcauth.GRPCAuthAPIKey,
		internallogger.Logger,
	)

	// Create the gRPC auth client
	conn, err := grpc.NewClient(
		fmt.Sprintf(
			"%s:%d",
			internalgrpcauth.GRPCAuthHost,
			internalgrpcauth.GRPCAuthPort,
		),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithChainUnaryInterceptor(
			apiKeyInterceptor.Authenticate(),
			authJWTInterceptor.VerifyAuthentication(),
		),
	)
	if err != nil {
		panic(err)
	}
	defer func(conn *grpc.ClientConn) {
		if err = conn.Close(); err != nil {
			panic(err)
		}
	}(conn)

	// Create the auth client
	internalgrpcauth.Client = pbauthcompiled.NewAuthClient(conn)

	// Start the RabbitMQ service on a separate goroutine
	go func() {
		if err = internalrabbitmq.RabbitMQConsumerService.Start(
			ctx,
		); err != nil {
			internallogger.Logger.Error(
				"Error starting RabbitMQ service",
				slog.String("error", err.Error()),
			)
		}
	}()

	// Get the last sync time registered on the JWT sync service
	lastSyncTokensUpdateAt, err := internalsqlite.SyncService.GetLastSyncTokensUpdatedAt(ctx)
	if err != nil {
		panic(err)
	}

	// Update the last sync time on the JWT sync service
	if err = internalsqlite.SyncService.UpdateLastSyncTokensUpdateAt(
		ctx,
		time.Now(),
	); err != nil {
		panic(err)
	}

	// Load refresh tokens from gRPC auth client
	stream, err := internalgrpcauth.Client.ListTokensToService(
		ctx, &pbauthcompiled.ListTokensToServiceRequest{
			IssuedAfter: timestamppb.New(lastSyncTokensUpdateAt),
		},
	)
	if err != nil {
		panic(err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// Stream ended
			break
		}
		if err != nil {
			// Handle error
			panic(err)
		}

		// Get the refresh tokens and access tokens from the message
		refreshTokens := msg.GetRefreshTokens()
		accessTokens := msg.GetAccessTokens()

		// Insert the refresh tokens into the RabbitMQ service
		for i, refreshToken := range refreshTokens {
			if err = internalrabbitmq.RabbitMQConsumerService.AddRefreshToken(
				refreshToken.GetId(),
				refreshToken.GetExpiresAt().AsTime(),
			); err != nil {
				panic(err)
			}

			// Get the access token ID based on the index
			accessToken := accessTokens[i]

			// Insert the access tokens into the RabbitMQ service
			if err = internalrabbitmq.RabbitMQConsumerService.AddAccessToken(
				accessToken.GetId(),
				refreshToken.GetId(),
				accessToken.GetExpiresAt().AsTime(),
			); err != nil {
				panic(err)
			}
		}
	}

	// Create the main router
	router, err := gonethttproute.NewBaseRouter(
		ModeFlag,
		internaljson.Handler,
		internallogger.Logger,
	)
	if err != nil {
		panic(err)
	}

	// Log the serving of the Swagger UI
	absPath, err := filepath.Abs("./docs")
	if err != nil {
		panic(err)
	}
	internallogger.Logger.Info(
		"Serving Swagger UI",
		slog.String("docs_path", absPath),
	)

	// Serve the Swaggers docs
	router.ServeStaticFiles(
		"/docs/", absPath,
	)

	// Create the main router module
	if err = internalrouter.Module.Create(router); err != nil {
		panic(err)
	}

	// Serve the REST API server
	internallogger.Logger.Info(
		"REST API server started",
		slog.Int("port", Port),
	)
	if err = http.ListenAndServe(
		fmt.Sprintf(":%d", Port),
		router.Handler(),
	); err != nil {
		panic(err)
	}
}
