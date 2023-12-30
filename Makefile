NAME = wgcf-cli

VERSION=$(shell git fetch --tags 1>&2; git describe --tags --always --dirty)

PARAMS = -trimpath -ldflags "-X github.com/ArchiveNetwork/wgcf-cli/constant.Version=$(VERSION) -s -w -buildid=" -v
MAIN = ./
PREFIX ?= $(shell go env GOPATH)

.PHONY: clean

build:
	go build -o $(NAME) $(PARAMS) $(MAIN)

install:
	go build -o $(PREFIX)/bin/$(NAME) $(PARAMS) $(MAIN)

clean:
	go clean -v -i $(PWD)
	rm -f wgcf-cli