#!/bin/bash

docker-compose -f ./scripts/pgloader.docker-compose.yaml up
docker container rm -f pgloader