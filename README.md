# uru-mobiles-recipes-api

Recipes HTTP API ("Cooking REST API") for the Mobile Development college course, written in Go. It is the backend consumed by `uru-mobiles-recipes-app`, delegating auth to a separate `grpc-auth` gRPC service and deployed together with it via `uru-mobiles-recipes`.

Go module: `github.com/ralvarezdev/uru-mobiles-recipes-api` (Go 1.25.1).

## Architecture

Built around a module/router tree (`go-net` `http.Module`):

- **gRPC auth client** (`internal/grpc/auth`) — talks to `grpc-auth` via API-key and JWT interceptors.
- **SQLite** (`internal/databases/sqlite`) — local token sync/validation store.
- **Redis** (`internal/databases/redis`) — backs the rate limiter.
- **RabbitMQ** (`internal/rabbitmq`) — consumes a token queue to keep SQLite in sync with `grpc-auth`.
- **Cookie-based JWT auth** — access tokens read from the `access_token` cookie.
- **Middlewares** — body limiting, rate limiting, error handling, security headers (`go-security-headers`).
- **Docs** (`internal/router/docs`) — generated Swagger UI.

On startup the server connects to SQLite, opens the gRPC auth client, starts the RabbitMQ consumer, replays tokens issued by `grpc-auth` since the last sync, then starts the HTTP server.

## API

All routes under `/api/v1`: `GET /ping` (health check); auth routes proxied to `grpc-auth` (`/auth/signup`, `/login`, `/refresh-token`, `/logout`, `/refresh-tokens`); user routes proxied to `grpc-auth` (`/user/profile`, `/username`). Responses follow JSend. A `recipe` domain model (`internal/router/api/v1/recipe/model.go`) defines `Recipe` and `Group`, but as of this snapshot no recipe-specific routes or SQL tables are wired in — `sql/create_tables.sql` and `sql/create_functions.sql` are empty placeholders.

## Configuration

Environment variables (see `.env.example`): gRPC auth (`GRPC_AUTH_HOST/PORT/API_KEY/JWT_PUBLIC_KEY`), `BODY_LIMIT`, Redis (`REDIS_ADDRESS/USERNAME/PASSWORD/DB`), RabbitMQ (`RABBITMQ_URL`, `RABBITMQ_TOKEN_QUEUE_NAME`), rate limiter (`RATE_LIMITER_MAX_REQUESTS/PERIOD`). CLI flags `-mode` and `-port` are also accepted, set by `serve.sh`.

## Running

```bash
cp .env.example .env
./compile.sh             # builds ./bin/server/server
./serve.sh                # runs with -mode and -port=8080
```

With Docker:

```bash
docker build --secret id=GITHUB_TOKEN,src=./.github_token -t recipes-api .
docker run -p 8080:8080 --env-file .env recipes-api
```

## Documentation

API docs are generated with `swag` (`generate_docs.sh`) and served via `internal/router/docs` as Swagger UI. Linting via `.golangci.yml`.

## License

GNU General Public License v3.0 (see `LICENSE`).
