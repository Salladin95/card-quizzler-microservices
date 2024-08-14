#!/bin/bash

# Get the directory where the script is located
script_dir=$(dirname "$(realpath "$0")")

# Define the Docker image name
image_name="khalid95/card-quizzler-api:1.0.0"

echo "Building image $image_name from Dockerfile located in $script_dir..."

# Build the Docker image using the Dockerfile in the script's directory
docker build -f "$script_dir/Dockerfile" -t $image_name "$script_dir"
if [ $? -eq 0 ]; then
    echo "Pushing image $image_name..."
    docker push $image_name
else
    echo "Failed to build image $image_name."
    exit 1
fi
