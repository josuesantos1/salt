SHELL := /bin/bash

.PHONY: clean dev run build

run:
	./salt

dev:
	go run .

build:
	go build -o salt .

clean:
	go clean
