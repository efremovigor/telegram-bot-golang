#!/bin/sh
set -e
echo "Build app"

cd /app
go mod download
go build -o /app/build/go_app
echo "Running app"
exec /app/build/go_app
