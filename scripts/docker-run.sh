#!/bin/bash
BINARY_NAME=minitwit-go
CONTAINER_NAME=minitwit-go

docker run -d -p 8080:8080 --name=$CONTAINER_NAME $BINARY_NAME

