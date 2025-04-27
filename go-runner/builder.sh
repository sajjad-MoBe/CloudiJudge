#!/bin/bash

TEMP_DIR="/home/runner/tmp"
if [ -d "$TEMP_DIR" ]; then
    runner
    exit 0
else
    mkdir "$TEMP_DIR"
    cp "/mnt/problem/code.go" "$TEMP_DIR/code.go"
    cp "/mnt/problem/input.txt" "$TEMP_DIR/input.txt"

    go build -o "$TEMP_DIR/code" "$TEMP_DIR/code.go"
    if [ $? -ne 0 ]; then
        echo "Compilation failed"
        rm -rf "$TEMP_DIR"
        exit 1
    fi
    echo "Compilation success"
    exit 0
fi