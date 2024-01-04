NAME := wgcf-cli

VERSION := $(shell git fetch --tags 1>&2; git describe --tags --always --dirty)

PARAMS := -trimpath -ldflags "-X github.com/ArchiveNetwork/wgcf-cli/constant.Version=$(VERSION) -s -w -buildid=" -v
MAIN := ./
PREFIX ?= $(shell go env GOPATH)
ifeq ($(GOOS),windows)
OUTPUT := $(NAME).exe
else
OUTPUT := $(NAME)
endif
.PHONY: clean

build:
	go build -o $(OUTPUT) $(PARAMS) $(MAIN)

install:
	go build -o $(PREFIX)/bin/$(OUTPUT) $(PARAMS) $(MAIN)

clean:
	go clean -v -i $(PWD)
	rm -f wgcf-cli wgcf-cli.exe