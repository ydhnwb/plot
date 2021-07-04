#!/bin/sh

echo "Running test case 1"
go run ./parking_lot/parking_lot.go -url=http://localhost:8080 -case=1


echo "Running test case 2"
go run ./parking_lot/parking_lot.go -url=http://localhost:8080 -case=2
