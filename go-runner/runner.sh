#!/bin/bash

if [ -z "$1" ] || [ -z "$2" ]; then
  echo "Usage: $0 <code_number> <time_limit_ms>"
  exit 1
fi

CODE=$1
TIME_LIMIT_MS=$2

TEMP_DIR=$(mktemp -d)
chmod 700 "$TEMP_DIR"

cp "/mnt/problem/$CODE.go" "$TEMP_DIR/code.go"
cp "/mnt/problem/input.txt" "$TEMP_DIR/input.txt"

go build -o "$TEMP_DIR/code" "$TEMP_DIR/code.go"
if [ $? -ne 0 ]; then
  echo "Compilation failed"
  rm -rf "$TEMP_DIR"
  exit 1
fi

timeout --signal=SIGKILL "${TIME_LIMIT_MS}ms" \
  /usr/bin/runuser -u appuser -- \
  /usr/bin/env -i \
  "$TEMP_DIR/code" < "$TEMP_DIR/input.txt" > "$TEMP_DIR/actual_output.txt" 2>/dev/null

EXIT_CODE=$?
if [ $EXIT_CODE -eq 137 ]; then
  echo "Time limit exceeded"
  rm -rf "$TEMP_DIR"
  exit 1
elif [ $EXIT_CODE -ne 0 ]; then
  echo "Runtime error"
  rm -rf "$TEMP_DIR"
  exit 1
fi

diff -w "$TEMP_DIR/actual_output.txt" "/mnt/problem/output.txt" > /dev/null
if [ $? -eq 0 ]; then
  echo "Output matches"
else
  echo "Output does not match"
fi

rm -rf "$TEMP_DIR"


docker run --rm \
  --cpus="1" \
  --memory="256m" \
  --network=none \
  --ulimit nofile=1024:1024 \
  --cap-drop=ALL \
  --security-opt=no-new-privileges \
  -v "$(pwd)/folders:/mnt/folders:ro" \
  go-code-runner 4 100