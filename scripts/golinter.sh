#!/bin/bash

# Install with 'go install golang.org/x/lint/golint' or 'go get -u golang.org/x/lint/golint'

# Test by adding this example go code to a file:
#	s := struct {
#		__a string
#	}{}
#
#	print(s)
cd scripts
cd ..

#golint -set_exit_status ./...
golint ./...