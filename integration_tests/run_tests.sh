#!/usr/bin/env bash

# Exit on error
set -e

# Kill any existing process on port 8088
if command -v lsof >/dev/null 2>&1; then
    # For Unix-like systems with lsof
    kill $(lsof -t -i:8088) 2>/dev/null || true
elif command -v netstat >/dev/null 2>&1; then
    # For systems with netstat (including Git Bash on Windows)
    pid=$(netstat -ano | grep "8088" | awk '{print $5}') && kill $pid 2>/dev/null || true
fi

# Function to cleanup
cleanup() {
    echo "Cleaning up..."
    if [ -n "$SERVER_PID" ]; then
        kill $SERVER_PID 2>/dev/null || true
    fi
    cd "$ORIGINAL_DIR"
}

# Store original directory
ORIGINAL_DIR=$(pwd)

# Setup cleanup on script exit
trap cleanup EXIT

# Start the server in background
cd examples/basic
go run main.go &
SERVER_PID=$!

# Wait for server to be ready
echo "Waiting for server to start..."
for i in {1..30}; do
    if curl -s http://localhost:8088/std > /dev/null; then
        echo "Server is ready!"
        break
    fi
    sleep 1
    if [ $i -eq 30 ]; then
        echo "Server failed to start within 30 seconds"
        exit 1
    fi
done

# Run the tests
cd ../../integration_tests
hurl --test basic.hurl 