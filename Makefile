.PHONY: all clean
.DEFAULT_GOAL := all

DEBUG ?= false

deps:
	go get ./...

test:
	go test

all: deps test build

build:
	go build -o rpi-moteino-collector

clean:
	rm -f rpi-moteino-collector

run:
	sudo DEBUG="$(DEBUG)" GOPATH=/home/ash/go /usr/local/go/bin/go run main.go $(filter-out $@, $(MAKECMDGOALS))
