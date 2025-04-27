#!/bin/bash

TEMP_DIR="/home/runner/tmp"

timeout -s KILL "$TIME_LIMIT" "$TEMP_DIR/code" < "$TEMP_DIR/input.txt" > "$TEMP_DIR/actual_output.txt" 2> /dev/null
EXIT_CODE=$?
if [ $EXIT_CODE -eq 137 ]; then
  echo "Time limit exceeded"
  rm -rf "$TEMP_DIR/"
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
