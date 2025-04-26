#!/bin/bash

if [ -z "$1" ] || [ -z "$2" ]; then
  echo "Usage: $0 <time_limit_s> <memory_limit_mb>"
  exit 1
fi

TIME_LIMIT_S=$1
MEM_LIMIT_MB=$2

TEMP_DIR=$(mktemp -d)

cp "/mnt/problem/code.go" "$TEMP_DIR/code.go"
cp "/mnt/problem/input.txt" "$TEMP_DIR/input.txt"

go build -o "$TEMP_DIR/code" "$TEMP_DIR/code.go"
if [ $? -ne 0 ]; then
  echo "Compilation failed"
  rm -rf "$TEMP_DIR"
  exit 1
fi

ulimit -v $((MEM_LIMIT_MB * 1024 * 1024))

timeout -s KILL "$TIME_LIMIT_S" "$TEMP_DIR/code" < "$TEMP_DIR/input.txt" > "$TEMP_DIR/actual_output.txt" 2>/dev/null
EXIT_CODE=$?
if [ $EXIT_CODE -eq 137 ]; then
  echo "Time limit exceeded"
  rm -rf "$TEMP_DIR"
  exit 0
elif [ $EXIT_CODE -eq 139 ]; then
  echo "Memory limit exceeded"
  rm -rf "$TEMP_DIR" "$EXEC_DIR"
  exit 0
elif [ $EXIT_CODE -ne 0 ]; then
  echo "Runtime error"
  rm -rf "$TEMP_DIR"
  exit 0
fi
if grep -q -E "Wrong answer|Compilation failed|Time limit exceeded|Memory limit exceeded|Runtime error" "$TEMP_DIR/actual_output.txt"; then
    echo "Wrong answer"
    rm -rf "$TEMP_DIR"
    exit 0
fi

cat "$TEMP_DIR/actual_output.txt"
rm -rf "$TEMP_DIR"

# docker build -t go-code-runner ./go-runner 
