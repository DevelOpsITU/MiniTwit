#!/bin/bash

display_usage() {
	echo -e "'Usage: $0 \$DOCKER_USERNAME \$DOCKER_PASSWORD' or 'Usage: $0 \$DOCKER_USERNAME \$DOCKER_PASSWORD dev' depending
	 \nwhether you want to release 'minitwit-go' or 'minitwit-go-dev' image to dockerhub."
}

# check whether user had supplied -h or --help . If yes display usage
if [[ ( $1 == "--help") ||  ($1 == "-h" )]]
then
  display_usage
  exit 0
fi

if [[ -z $1 ]];
then
    echo "${0##*/}:  No username passed"
    exit 1
else
    username=$1
fi
if [[ -z $2 ]];
then
    echo "${0##*/}:  No password passed"
    exit 2
else
    password=$2
fi

docker login -u "$username" -p "$password"

BINARY_NAME=minitwit-go-dev
DOCKER_REGISTRY=groupddevops/
VERSION=$(git rev-parse --short HEAD)

echo "tagging:" $BINARY_NAME:$VERSION $DOCKER_REGISTRY$BINARY_NAME:$VERSION
echo "pushing " $DOCKER_REGISTRY$BINARY_NAME:$VERSION "to dockerhub"

echo "docker tag"
docker tag $BINARY_NAME $DOCKER_REGISTRY$BINARY_NAME:$VERSION
# Push the docker images
echo "docker push"
docker push $DOCKER_REGISTRY$BINARY_NAME:$VERSION


if [[ -z $3 ]];
then
    echo "production release choosen, also tagging with latest"
    echo "tagging:" $BINARY_NAME:$VERSION $DOCKER_REGISTRY$BINARY_NAME:latest
    echo "pushing " $DOCKER_REGISTRY$BINARY_NAME:latest "to dockerhub"
    docker tag $BINARY_NAME $DOCKER_REGISTRY$BINARY_NAME:latest
    docker push $DOCKER_REGISTRY$BINARY_NAME:latest
    echo "Done"
else
    echo "Done"

fi