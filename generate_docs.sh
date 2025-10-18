#!/bin/bash

echo "Generating Swagger docs..."
swag init -g ./cmd/server/main.go --parseDependency --parseInternal
echo "Generating Swagger docs... Done."