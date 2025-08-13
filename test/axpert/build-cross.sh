#!/bin/bash

# Cross-compilation script for Raspberry Pi 4 (linux/arm64) using Docker

echo "Building axpert-test for Raspberry Pi 4 (linux/arm64) using Docker..."

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed or not in PATH"
    exit 1
fi

if ! docker info &> /dev/null; then
    echo "Error: Docker is not running"
    exit 1
fi

# Build with Docker
echo "Building with Docker..."
cd ../../
docker buildx build --platform linux/arm64 -f test/axpert/Dockerfile.cross -t energia-cross --load .
if [ $? -ne 0 ]; then
    echo "Error: Docker build failed"
    exit 1
fi

docker create --name temp-container energia-cross
if [ $? -ne 0 ]; then
    echo "Error: Failed to create container"
    exit 1
fi

docker cp temp-container:/app/axpert-test ./test/axpert/axpert-test-arm64
if [ $? -ne 0 ]; then
    echo "Error: Failed to copy binary from container"
    docker rm temp-container 2>/dev/null
    exit 1
fi

docker rm temp-container
echo "Binary created: axpert-test-arm64"
