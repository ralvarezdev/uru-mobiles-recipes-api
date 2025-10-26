#!/bin/sh

# Set environment variable to enable CGO for cross-compilation
export CGO_ENABLED=1

echo "Compiling server..."
go mod tidy
go build -o bin/server/server ./cmd/server
echo "Compiling server... Done"