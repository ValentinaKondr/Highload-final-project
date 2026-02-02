#!/bin/bash

set -e

echo "Building Go service Docker image..."

docker build -t metrics-analyzer:latest .

echo "Docker image built successfully!"
echo "To load into Minikube, run: minikube image load metrics-analyzer:latest"

