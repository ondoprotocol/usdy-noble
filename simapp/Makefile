BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')

ifeq (,$(VERSION))
  VERSION := $(shell git describe --exact-match 2>/dev/null)
  ifeq (,$(VERSION))
    ifeq ($(shell git status --porcelain),)
    	VERSION := $(BRANCH)
    else
    	VERSION := $(BRANCH)-dirty
    endif
  endif
endif

ldflags := $(LDFLAGS)
ldflags += -X github.com/cosmos/cosmos-sdk/version.Name=simapp \
	-X github.com/cosmos/cosmos-sdk/version.AppName=simd \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
	-X github.com/cosmos/cosmos-sdk/version.BuildTags='ledger'
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -ldflags '$(ldflags)' -tags 'ledger'

build:
	@go build -mod=readonly $(BUILD_FLAGS) -o $(PWD)/build/ ./...

.PHONY: build
