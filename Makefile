.PHONY: fresh

SHELL := /bin/bash

fresh:
	ENV=LOCAL \
	PORT=3001 \
	VERSION=VERSION \
	FIXTURES=fixtures.json \
	fresh
