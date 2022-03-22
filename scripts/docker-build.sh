#!/bin/bash
BINARY_NAME=minitwit-go-dev
VERSION=$(git rev-parse --short HEAD)
docker build --tag "$BINARY_NAME":latest .