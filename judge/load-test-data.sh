#!/bin/sh
set -e
echo "start generation of test data..."
# go run . load-test-data "$@"
./judge load-test-data "$@"
