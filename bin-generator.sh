#!/bin/bash

# Build for Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -o anto-linux

# Build for macOS (64-bit)
GOOS=darwin GOARCH=amd64 go build -o anto-mac

# Build for Windows (64-bit)
GOOS=windows GOARCH=amd64 go build -o anto.exe

echo "Builds completed!"
