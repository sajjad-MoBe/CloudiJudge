#!/bin/sh
set -e
echo "Starting Golang web server with 'code-runner' command..."
go run . code-runner "$@"