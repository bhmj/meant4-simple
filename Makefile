PROJECT ?=$(shell basename $(PWD))
SRC ?= ./cmd
BINARY ?= ./build

help:
	printf "Usage: make <command>\n\n<command>:\n\n  build-v1\t- build v1 (old and dumb one)\n  build-v2\t- build v2 (improved)\n  run-v1\n\n"

all: build-v1 build-v2

build-v1:
	mkdir -p $(BINARY)
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -trimpath -o $(BINARY) $(SRC)/version-1

build-v2:
	mkdir -p $(BINARY)
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -trimpath -o $(BINARY) $(SRC)/version-2

run-v1:
	CGO_ENABLED=0 go run -ldflags '$(LDFLAGS)' -trimpath $(SRC)/version-1

run-v2:
	CGO_ENABLED=0 go run -ldflags '$(LDFLAGS)' -trimpath $(SRC)/version-2

lint:
	golangci-lint run

test:
	go test ./...

.PHONY: build-v1 build-v2 run-v1 run-v2 lint test

$(V).SILENT:
