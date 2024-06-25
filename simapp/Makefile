BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true

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

build_tags := $(BUILD_TAGS)
build_tags += netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

build_tags := $(strip $(build_tags))
whitespace :=
whitespace += $(whitespace)
comma := ,

ldflags := $(LDFLAGS)
ldflags += -X github.com/cosmos/cosmos-sdk/version.Name=simapp \
	-X github.com/cosmos/cosmos-sdk/version.AppName=simd \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
	-X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(subst $(whitespace),$(comma),$(build_tags))"
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -ldflags '$(ldflags)' -tags '$(build_tags)' -trimpath

build:
	@go build -mod=readonly $(BUILD_FLAGS) -o $(PWD)/build/ ./...

.PHONY: build
