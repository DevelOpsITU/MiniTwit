#!/bin/bash
BINARY_NAME=minitwit-go
DOCKER_REGISTRY=groupddevops/
VERSION=$(git rev-parse --short HEAD)

docker tag $BINARY_NAME $DOCKER_REGISTRY$BINARY_NAME:latest
docker tag $BINARY_NAME $DOCKER_REGISTRY$BINARY_NAME:$VERSION
# Push the docker images
docker push $DOCKER_REGISTRY$BINARY_NAME:latest
docker push $DOCKER_REGISTRY$BINARY_NAME:$VERSION