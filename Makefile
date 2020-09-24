GOLINT:=$(shell go list -f {{.Target}} golang.org/x/lint/golint)
TMCOMMIT := $(shell go list -m -u -f '{{.Version}}' github.com/tendermint/tendermint)
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')
LD_FLAGS = -X github.com/ovrclk/tmsigner/cmd.Version=$(VERSION) \
	-X github.com/ovrclk/tmsigner/cmd.Commit=$(COMMIT) \
	-X github.com/ovrclk/tmsigner/cmd.TMCommit=$(TMCOMMIT)
BUILD_FLAGS := -ldflags '$(LD_FLAGS)'


all: build

gomod:
	@go mod tidy

build: clean gomod
ifeq ($(OS),Windows_NT)
	@echo "building tmsigner binary..."
	@CGO_ENABLED=0 go build -mod=readonly $(BUILD_FLAGS) -o build/tmsigner.exe main.go
else
	@echo "building tmsigner binary..."
	@CGO_ENABLED=0 go build -mod=readonly $(BUILD_FLAGS) -o build/tmsigner main.go
endif

install: build
	@mv ./build/tmsigner $(GOBIN)
	# @go install -mod=readonly $(BUILD_FLAGS) main.go

lint: tools
	@$(GOLINT) -set_exit_status ./...

test:
	@go test -short ./...

race:
	@go test -race -short ./...

msan:
	@go test -msan -short ./...

tools:
	@go install golang.org/x/lint/golint

test-env: install
	@./scripts/build-simd.bash local skip
	@./scripts/config-signer.bash skip signerchain
	@./scripts/two-node-net.bash signerchain build
	@tmsigner start

clean:
	rm -rf build

.PHONY: all lint test race msan tools clean build
