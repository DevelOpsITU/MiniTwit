#!/bin/bash
# Taken from https://github.com/beevelop/docker-scancode
echo $(pwd)
docker run -it -v $(pwd)/tests/:/scan groupddevops/scancode
pip install jq
cat tests/licenses.json | jq .