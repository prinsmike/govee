#!/usr/bin/env bash

docker run -it --rm -v $PWD:/usr/src/myapp \
	-e GO111MODULE=on \
	-w /usr/src/myapp \
	prinsmike/alpine-go-builder:latest \
	go test -v -coverprofile cover.out
