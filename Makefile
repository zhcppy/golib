#!/usr/bin/make -f

.PHONY: build-linux gomod install format

name = golib
VERSION = $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT = $(shell git log -1 --format='%H')
ldflags = '-X github.com/zhcppy/golib/version.Name=$(name) \
		   -X github.com/zhcppy/golib/version.Version=$(VERSION) \
           -X github.com/zhcppy/golib/version.Commit=$(COMMIT) \
           -w -s'

export GO111MODULE = on
export GOPROXY = https://goproxy.io,direct

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags $(ldflags) -o build/$(name) .

gomod:
	@go mod tidy && go mod download && go mod verify

install:
	@go install -ldflags $(ldflags) -v -o $ .

format:
	@find . -name '*.go' -type f -not -path "*.git*" | xargs goimports -w -e -l
