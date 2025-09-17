#!/usr/bin/env bash

# This tells Railway to use Go 1.21 and build from the correct directory
echo "Building Go service..."

# Install Go dependencies
go mod download
go mod tidy

# Build the service
go build -o app ./cmd/server

echo "Build complete!"