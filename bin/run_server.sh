#!/bin/sh

echo "Downloading dependencies..."
go mod download

echo "Running server..."
go run main.go

