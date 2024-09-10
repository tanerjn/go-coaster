#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Load .env variables
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
else
  echo ".env file not found!"
  exit 1
fi

# Build Docker images for server and client
echo "Building Docker images..."

# Navigate to the server directory and build server image
cd server
docker build -t ${SERVER_IMAGE} .
echo "Server image built: ${SERVER_IMAGE}"

# Navigate to the client directory and build client image
cd ../client
docker build -t ${CLIENT_IMAGE} .
echo "Client image built: ${CLIENT_IMAGE}"

# Tag the images for Docker Hub
echo "Tagging Docker images for Docker Hub..."
docker tag ${SERVER_IMAGE} ${DOCKER_USERNAME}/${SERVER_IMAGE}:${SERVER_TAG}
docker tag ${CLIENT_IMAGE} ${DOCKER_USERNAME}/${CLIENT_IMAGE}:${CLIENT_TAG}

# Push the images to Docker Hub
echo "Pushing Docker images to Docker Hub..."
docker push ${DOCKER_USERNAME}/${SERVER_IMAGE}:${SERVER_TAG}
docker push ${DOCKER_USERNAME}/${CLIENT_IMAGE}:${CLIENT_TAG}

echo "Docker images pushed successfully!"
