#!/bin/bash

docker-compose -f ./scripts/updateSqlitePw_hash.docker-compose.yaml up
docker container rm -f sqlite