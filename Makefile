GOLINT:=$(shell go list -f {{.Target}} golang.org/x/lint/golint)

all: build

build: 
	@CGO_ENABLED=0 go build -o ./build/tmsigner main.go

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
