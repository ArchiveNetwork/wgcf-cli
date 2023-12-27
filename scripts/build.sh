#!/bin/bash
set -e
cd "$(realpath "${0%scripts/build.sh}")"
REVISION=$(git rev-parse HEAD)
VERSION=$(git describe --tags --always --dirty)

build() {
    cleanup() {
        sed -i "s/${VERSION}/__VERSION__/g" src/command.go
        sed -i "s/${REVISION}/__REVISION__/g" src/command.go
    }
    trap 'cleanup ;echo "Building wgcf-cli failed."' INT TERM
    trap 'cleanup ;echo "Building wgcf-cli done."' EXIT
    echo "Building wgcf-cli..."
    sed -i "s/__REVISION__/$REVISION/g" src/command.go
    sed -i "s/__VERSION__/$VERSION/g" src/command.go
    go build -trimpath -ldflags "-s -w -buildid=" -v -o wgcf-cli ./src/
}

clean() {
    echo "Cleaning..."
    go clean -v -i ./src/
    rm -f wgcf-cli
    echo "Cleaning done."
}

if [[ "$1" != "build" && "$1" != "clean" ]]; then
    echo "Usage: $0 [build|clean]"
    exit 1
fi
$1