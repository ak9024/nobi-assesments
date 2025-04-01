#!/bin/bash

docker compose up -d --build
echo "Waiting for services to be ready..."
# Use a more reliable approach to wait for services instead of a fixed sleep
timeout=60
elapsed=0
while [ $elapsed -lt $timeout ]; do
  if curl -s http://localhost:3000/health >/dev/null; then
    echo "Services are ready!"
    break
  fi
  sleep 2
  elapsed=$((elapsed + 2))
  echo "Still waiting... ($elapsed/$timeout seconds)"
done

if [ $elapsed -ge $timeout ]; then
  echo "Timeout waiting for services to be ready"
  docker compose down -v
  exit 1
fi

# Run tests with proper error handling
echo "Running tests..."
if go test ./tests/api_test.go; then
  echo "Tests completed successfully"
else
  echo "Tests failed with exit code $?"
  docker compose down -v
  exit 1
fi

# Clean up
echo "Cleaning up..."
docker compose down -v
echo "Done!"

docker compose up -d --build
echo "Waiting for services to be ready..."
# Use a more reliable approach to wait for services instead of a fixed sleep
timeout=60
elapsed=0
while [ $elapsed -lt $timeout ]; do
  if curl -s http://localhost:3000/health >/dev/null; then
    echo "Services are ready!"
    break
  fi
  sleep 2
  elapsed=$((elapsed + 2))
  echo "Still waiting... ($elapsed/$timeout seconds)"
done

if [ $elapsed -ge $timeout ]; then
  echo "Timeout waiting for services to be ready"
  docker compose down -v
  exit 1
fi

# Run tests with proper error handling
echo "Running tests..."
if go test ./tests/api_test.go; then
  echo "Tests completed successfully"
else
  echo "Tests failed with exit code $?"
  docker compose down -v
  exit 1
fi

# Clean up
echo "Cleaning up..."
docker compose down -v
echo "Done!"
