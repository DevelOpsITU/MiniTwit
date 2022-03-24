#!/bin/bash

display_usage() {
	echo -e "'Usage: $0' or 'Usage: $0 dev' depending \nwhether you want to build 'minitwit-go' or 'minitwit-go-dev' image \n"
}

# check whether user had supplied -h or --help . If yes display usage
if [[ ( $1 == "--help") ||  ($1 == "-h" )]]
then
  display_usage
  exit 0
fi

if [[ -z $1 ]];
then
    echo 'Building minitwit-go'
    BINARY_NAME=minitwit-go
else
    echo 'Building minitwit-go-dev'
    BINARY_NAME=minitwit-go-dev
fi

VERSION=$(git rev-parse --short HEAD)
docker build --tag "$BINARY_NAME":latest .