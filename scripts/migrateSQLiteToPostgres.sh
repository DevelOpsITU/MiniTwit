#!/bin/bash

docker-compose -f ./pgloader.docker-compose.yaml up
docker container rm -f pgloader