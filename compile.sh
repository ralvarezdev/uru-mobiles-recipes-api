#!/bin/sh

echo "Compiling server..."
go mod tidy
CGO_ENABLED=1 go build -o bin/server/server ./cmd/server
echo "Compiling server... Done"