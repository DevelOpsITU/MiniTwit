#!/bin/bash
#export DB_CONNECTION_STRING="file::memory:"
go mod tidy
cd ..
go test -v ./...