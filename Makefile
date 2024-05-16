NAME = wgcf-cli
MAIN = ./cmd/wgcf-cli

VERSION ?= $(shell git describe --tags --long | sed 's/^v//;s/\([^-]*-g\)/r\1/;s/-\([^-]*\)-\([^-]*\)$$/.\1.\2/;s/-//')

export GOOS ?= $(shell go env GOOS)
export GOARCH ?= $(shell go env GOARCH)
export GOARM ?= $(shell go env GOARM)
export CGO_ENABLED ?= 0
BUILD_MODE = -buildmode=pie
ifeq ($(CGO_ENABLED),1)
LDFLAG_LINKMODE = -linkmode=external
else
LDFLAG_LINKMODE = 
# armv7, riscv64, s390x, x86
ifeq ($(shell if [ "$(GOOS)" = "linux" ] && ([ "$(GOARCH)" = "386" ] || [ "$(GOARCH)" = "riscv64" ] || [ "$(GOARCH)" = "s390x" ] || ([ "$(GOARCH)" = "arm" ] && [ "$(GOARM)" = "7" ])); then echo true; fi) ,true)
BUILD_MODE = -buildmode=exe
endif
endif

TAGS ?= 
LDFLAGS = -X github.com/ArchiveNetwork/wgcf-cli/constant.Version=$(VERSION) -s -w -buildid= $(LDFLAG_LINKMODE)
GOFLAGS ?= -trimpath $(BUILD_MODE) -tags=$(TAGS) -mod=readonly -modcacherw -v -ldflags "$(LDFLAGS)"

ifeq ($(shell if [ "$(GOOS)" = "windows" ]; then echo true; fi),true)
OUTPUT = $(NAME).exe
else
OUTPUT = $(NAME)
endif
ifdef completion
all: completion
else
all: $(NAME)
endif

.PHONY: clean completion

$(NAME): $(MAIN)
	go build -o $(OUTPUT) \
	 	$(GOFLAGS) \
	 	$(MAIN)

clean:
	go clean -v -i $(MAIN)
	rm -f $(OUTPUT)

completion: 
	@$(MAKE) $(NAME) 2>&1 >/dev/null
	@PATH=$(PATH):$(shell realpath ./)
	@$(NAME) completion $(completion)