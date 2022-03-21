#!/bin/bash
# Taken from https://github.com/beevelop/docker-scancode
echo $(pwd)
docker run -v $(pwd)/:/scan groupddevops/scancode
pip install jq
echo ''
echo '-----------------------------------'
cat tests/licenses.json | jq .