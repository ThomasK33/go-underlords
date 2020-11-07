-include .env

# VERSION := $(shell git describe --tags 2> /dev/null || echo v0.0.0)
# BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
# GOFILES := $(wildcard *.go)
GOFILES=examples/main.go

# Use linker flags to provide version/build settings
# LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"
LDFLAGS=

# Redirect error output to a file, so we can show it in development mode.
STDERR := /tmp/.$(PROJECTNAME)-stderr.txt

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## install: Install missing dependencies. Runs `go get` internally. e.g; make install get=github.com/foo/bar
install: go-get

## run: Run the application
run:
	@-$(MAKE) go-run

## test: Run unit tests
test:
	@-$(MAKE) go-test

## compile: Compile the binary.
compile:
	@-touch $(STDERR)
	@-rm $(STDERR)
	@-$(MAKE) -s go-compile 2> $(STDERR)
	@cat $(STDERR) 1>&2

## clean: Clean build files. Runs `go clean` internally.
clean:
	@-rm ./bin/$(PROJECTNAME) 2> /dev/null
	@-$(MAKE) go-clean

go-compile: go-get go-build

go-build:
	@echo "  >  Building binary..."
	go build $(LDFLAGS) -o ./bin/$(PROJECTNAME) $(GOFILES)

go-generate:
	@echo "  >  Generating dependency files..."
	go generate $(generate)

go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	go get $(get) -v all

go-install:
	go install $(GOFILES)

go-run:
	go run $(LDFLAGS) $(GOFILES)

go-test:
	go test $(LDFLAGS) -race -covermode=atomic -coverprofile=coverage.out -v ./... 

go-clean:
	@echo "  >  Cleaning build cache"
	go clean

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
