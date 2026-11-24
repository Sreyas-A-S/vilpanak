#!/usr/bin/env bash
set -e

# Load environment if present
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

echo "Starting Vilpanak Real Estate Server..."
go run ./cmd/web/main.go
