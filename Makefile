GOLINT:=$(shell go list -f {{.Target}} golang.org/x/lint/golint)
TMCOMMIT := $(shell go list -m -u -f '{{.Version}}' github.com/tendermint/tendermint)
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')
LD_FLAGS = -X github.com/ovrclk/relayer/cmd.Version=$(VERSION) \
	-X github.com/ovrclk/relayer/cmd.Commit=$(COMMIT) \
	-X github.com/ovrclk/relayer/cmd.TMCommit=$(TMCOMMIT)
BUILD_FLAGS := -ldflags '$(LD_FLAGS)'



all: build

gomod:
	@go mod tidy

build: gomod
ifeq ($(OS),Windows_NT)
	@echo "building tmsigner binary..."
	@CGO_ENABLED=0 go build -mod=readonly $(BUILD_FLAGS) -o build/tmsigner.exe main.go
else
	@echo "building tmsigner binary..."
	@CGO_ENABLED=0 go build -mod=readonly $(BUILD_FLAGS) -o build/tmsigner main.go
endif

build-zip: gomod
	@echo "building tmsigner binaries for windows, mac and linux"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/linux-amd64-tmsigner main.go
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/darwin-amd64-tmsigner main.go
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/windows-amd64-tmsigner.exe main.go
	@tar -czvf release.tar.gz ./build

install: gomod
	@echo "installing tmsigner binary..."
	@go install -mod=readonly $(BUILD_FLAGS) main.go

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

clean:
	rm -rf build

.PHONY: all lint test race msan tools clean build
