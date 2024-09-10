#!/bin/bash

# Load .env variables
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
else
  echo ".env file not found!"
  exit 1
fi


# Check if Docker image exists on Docker Hub
#
check_image_exists() {
  repo=$1

  # Use Docker Hub API to check if the image repository exists
  response=$(curl -s -o /dev/null -w "%{http_code}" \
    -H "Authorization: Bearer ${DOCKER_TOKEN}" \
    https://hub.docker.com/v2/repositories/${DOCKER_USERNAME}/${repo}/)

  if [ "$response" -eq 200 ]; then
    return 0  # Image exists
  else
    return 1  # Image does not exist
  fi
}

# Function to delete all resources in OpenShift
delete_openshift_resources() {
  echo "...Checking OpenShift..."

  # Count the number of resources
  count_pods=$(oc get pods --no-headers 2>/dev/null | wc -l)
  count_routes=$(oc get routes --no-headers 2>/dev/null | wc -l)
  count_deployments=$(oc get deployments --no-headers 2>/dev/null | wc -l)

  # Print the number of each type of resource
  echo "Number of pods: ${count_pods}"
  echo "Number of routes: ${count_routes}"
  echo "Number of deployments: ${count_deployments}"

  # Check if any resources exist and delete them if they do
  if [ "$count_pods" -gt 0 ] || [ "$count_routes" -gt 0 ] || [ "$count_deployments" -gt 0 ]; then
    echo "Deleting all resources in OpenShift project..."
    
    oc delete pods --all
    oc delete routes --all
    oc delete deployments --all

    echo "All resources deleted in OpenShift."
  else
    echo "No resources found to delete in OpenShift."
  fi
}


# Function to delete Docker images from Docker Hub
delete_docker_images() {
  echo "...Checking Docker..."
  # Client image deletion
  if check_image_exists "${DOCKER_REPO_CLIENT}"; then
    echo "Deleting client image from Docker Hub..."
    curl -X DELETE \
      -H "Authorization: Bearer ${DOCKER_TOKEN}" \
      https://hub.docker.com/v2/repositories/${DOCKER_USERNAME}/${DOCKER_REPO_CLIENT}/
    echo "Client image ${DOCKER_REPO_CLIENT} deleted from Docker Hub."
  else
    echo "No client image found for ${DOCKER_REPO_CLIENT}."
  fi

  # Server image deletion
  if check_image_exists "${DOCKER_REPO_SERVER}"; then
    echo "Deleting server image from Docker Hub..."
    curl -X DELETE \
      -H "Authorization: Bearer ${DOCKER_TOKEN}" \
      https://hub.docker.com/v2/repositories/${DOCKER_USERNAME}/${DOCKER_REPO_SERVER}/
    echo "Server image ${DOCKER_REPO_SERVER} deleted from Docker Hub."
  else
    echo "No server image found for ${DOCKER_REPO_SERVER}."
  fi
}

# Main script
delete_openshift_resources
delete_docker_images

