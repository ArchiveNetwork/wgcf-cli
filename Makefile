NAME = wgcf-cli

VERSION = $(shell git fetch --tags 1>&2; git describe --tags --always --dirty)
export CGO_ENABLED ?= 0
GOFLAGS ?= -trimpath -ldflags "-X github.com/ArchiveNetwork/wgcf-cli/constant.Version=$(VERSION) -s -w -buildid=" -v
MAIN = ./cmd/wgcf-cli
PREFIX ?= $(shell go env GOPATH)

ifeq ($(GOOS),windows)
OUTPUT = $(NAME).exe
else
OUTPUT = $(NAME)
endif

.PHONY: build clean

build:
	go build -o $(OUTPUT) $(GOFLAGS) $(MAIN)

clean:
	go clean -v -i $(MAIN)
	rm -f $(OUTPUT)