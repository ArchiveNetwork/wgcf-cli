NAME = wgcf-cli

VERSION = $(shell git fetch --tags 1>&2; git describe --tags --always --dirty)

GOFLAGS ?= -trimpath -ldflags "-X github.com/ArchiveNetwork/wgcf-cli/constant.Version=$(VERSION) -s -w -buildid=" -v
MAIN = ./
PREFIX ?= $(shell go env GOPATH)

ifeq ($(GOOS),windows)
OUTPUT = $(NAME).exe
else
OUTPUT = $(NAME)
endif

.PHONY: build install clean

build:
	go build -o $(OUTPUT) $(GOFLAGS) $(MAIN)

install: build
	install -Dm 755 $(OUTPUT) $(PREFIX)/bin/$(OUTPUT)

clean:
	go clean -v -i $(MAIN)
	rm -f $(OUTPUT)