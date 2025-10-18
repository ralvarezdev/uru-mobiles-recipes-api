#!/bin/sh

# Load environment variables from .env file
set -a
if [ -f ./.env ]; then
  . ./.env
fi
set +a

# Execute the Go binary on port specified
exec ./bin/server/server -mode=prod -port=8080