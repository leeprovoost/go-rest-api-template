.PHONY: run build test

SHELL := /bin/bash

run:
	ENV=LOCAL \
	PORT=3001 \
	VERSION=VERSION \
	go run .

build:
	go build -o api-service .

test:
	go test ../../...
