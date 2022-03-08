#!/bin/bash
#export DB_CONNECTION_STRING="file::memory:"
go mod tidy
go mod download
ls
go test -v ./...