SOURCEDIR="."
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=rpi_moteino_collector
BINARY_RELEASE=bin/${BINARY}_${VERSION}

VERSION=$(shell cat VERSION)

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES) deps bin_dir
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -a -installsuffix cgo -o ${BINARY_RELEASE}_linux_arm

.PHONY: deps
deps:
	go get ./...

.PHONY: update_deps
update_deps:
	go get -u ./...

.PHONY: bin_dir
bin_dir:
	mkdir -p bin

.PHONY: run
run: deps
	go run main.go $(filter-out $@, $(MAKECMDGOALS))

.PHONY: clean
clean:
	rm -f ${BINARY} ${BINARY}_*
