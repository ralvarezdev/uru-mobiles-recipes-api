#!/bin/sh

echo "Compiling server..."
go mod tidy
go build -o bin/server/server ./cmd/server
echo "Compiling server... Done"