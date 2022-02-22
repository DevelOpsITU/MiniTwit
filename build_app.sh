#!/bin/bash
echo "Running from ${0%/*}"

echo "Built and tested on 'go version go1.17.6 linux/amd64'"
go_version=$(go version)
echo "You have version:   '$go_version'"
echo

DEFAULT_FILENAME="group_d_go_app.out"
filename=$DEFAULT_FILENAME

if [[ -z $1 ]];
then
    echo "${0##*/}:  No parameter passed: using '$DEFAULT_FILENAME' as output name for go output file"
else
    echo "Using '$1' as output name for go output file"
    filename=$1
fi


echo "Running build command: 'build -o \"out/$filename\" src/minitwit.go'"
go build -o out/"$filename" src/minitwit.go
res=$?
echo

# Colors from https://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux

if [ $res -eq 0 ]
then
    GREEN='\033[0;32m'
    echo -e "${GREEN}Built file \"$filename\" in /out folder successfully"
    exit 0
else
    RED='\033[0;31m'
    echo -e "${RED}PANIC, something was wrong with the compilation!"
    exit 1
fi