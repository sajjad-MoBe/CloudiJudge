#!/bin/sh
set -e
echo "start create admin process..."
# go run . create-admin "$@"
./judge create-admin "$@"
