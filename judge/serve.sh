#!/bin/sh
set -e
echo "Starting Golang web server with 'serve' command..."
# go run . serve "$@"
./judge serve "$@"