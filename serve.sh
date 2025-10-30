#!/bin/sh
set -e

# Load environment variables from .env file
set -a
if [ -f ./.env ]; then
  . ./.env
fi
set +a

# Check if the mode variable is set, default to 'dev' if not
if [ -z "$MODE" ]; then
  MODE="dev"
fi

# Execute the Go binary on port specified
exec ./bin/server/server -mode=$MODE -port=8080