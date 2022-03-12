#!/bin/bash

docker-compose -f ./SQLiteToPostgres.docker-compose.yaml up
docker container rm -f sqlite