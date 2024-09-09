#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Apply server deployment, service, and route
echo "Deploying server components in OpenShift..."
oc apply -f server/server-deployment.yaml
oc apply -f server/server-service.yaml
oc apply -f server/server-route.yaml
echo "Server deployment, service, and route applied."

# Apply client deployment, service, and route
echo "Deploying client components in OpenShift..."
oc apply -f client/client-deployment.yaml
oc apply -f client/client-service.yaml
oc apply -f client/client-route.yaml
echo "Client deployment, service, and route applied."

# Check the status of pods, deployments, services, and routes
echo "Checking OpenShift resources..."
oc get pods
oc get deployments
oc get services
oc get routes

