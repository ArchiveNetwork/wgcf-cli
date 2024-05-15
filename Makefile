NAME = wgcf-cli

VERSION ?= $(shell git describe --tags --long | sed 's/^v//;s/\([^-]*-g\)/r\1/;s/-\([^-]*\)-\([^-]*\)$$/.\1.\2/;s/-//')
export CGO_ENABLED ?= 0
ifeq ($(CGO_ENABLED),1)
LDFLAG_LINKMODE = -linkmode=external
else
LDFLAG_LINKMODE = 
endif
LDFLAGS = -X github.com/ArchiveNetwork/wgcf-cli/constant.Version=$(VERSION) -s -w -buildid= $(LDFLAG_LINKMODE)
GOFLAGS ?= -trimpath -mod=readonly -modcacherw -v -ldflags "$(LDFLAGS)"
MAIN = ./cmd/wgcf-cli
PREFIX ?= $(shell go env GOPATH)

ifeq ($(GOOS),windows)
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