SHELL := /bin/bash
BIN := dist/melon-status
HASH := $(shell git rev-parse --short HEAD)
COMMIT_DATE := $(shell git show -s --format=%ci ${HASH})
BUILD_DATE := $(shell date '+%Y-%m-%d %H:%M:%S')
VERSION := ${HASH}
LD_FLAGS := -s -w -X 'main.buildVersion=${VERSION}' -X 'main.buildDate=${BUILD_DATE}'
COMP_BIN := go

.PHONY: build dev clean

build:
	mkdir -p dist/
	${COMP_BIN} build -o "${BIN}" -ldflags="${LD_FLAGS}" ./cmd/melon-status

dev:
	${COMP_BIN} build -tags debug -o "${BIN}" -ldflags="${LD_FLAGS}" ./cmd/melon-status
	./${BIN}

clean:
	go clean
	rm ${BIN}
