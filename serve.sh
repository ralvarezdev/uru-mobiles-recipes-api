#!/bin/sh

# Load environment variables from .env file
set -a
if [ -f ./.env ]; then
  . ./.env
fi
set +a

# Check if port is provided via environment variables
if [ -z "$PORT" ]; then
    echo "Error: No port specified."
    exit 1
fi

# Execute the Go binary on port specified
exec ./bin/server/server -mode=prod -port=$PORT