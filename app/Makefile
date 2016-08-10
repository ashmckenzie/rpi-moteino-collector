SOURCEDIR="."
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=moteino_collector
BINARY_RELEASE=bin/${BINARY}_${VERSION}

VERSION=$(shell cat VERSION)

.DEFAULT_GOAL: $(BINARY)

$(BINARY): bin_dir deps $(SOURCES)
	go build -o bin/${BINARY}

releases: deps release_linux_amd64 release_linux_arm
#releases: deps release_linux_amd64 release_linux_arm release_darwin_amd64

release_linux_amd64: bin_dir
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${BINARY_RELEASE}_linux_amd64

release_linux_arm: bin_dir
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6  go build -a -installsuffix cgo -o ${BINARY_RELEASE}_arm_amd64

release_darwin_amd64: bin_dir
	CGO_ENABLED=0 GOOS=darwin go build -a -installsuffix cgo -o ${BINARY_RELEASE}_darwin_amd64

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
