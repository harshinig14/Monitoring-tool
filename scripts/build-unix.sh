#!/bin/bash

echo "Building Linux Binary..."

GOOS=linux GOARCH=amd64 go build -o dist/linux/MONITORING-TOOL ./cmd/agent

echo "Building macOS Intel Binary..."

GOOS=darwin GOARCH=amd64 go build -o dist/mac/MONITORING-TOOL ./cmd/agent

echo "Building macOS ARM Binary..."

GOOS=darwin GOARCH=arm64 go build -o dist/mac-arm/MONITORING-TOOL ./cmd/agent

echo "Unix Builds Completed"