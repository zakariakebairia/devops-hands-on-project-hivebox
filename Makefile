# Makefile for kvmcli project

# BINARY_NAME sets the name of the output executable.
BINARY_NAME = hivebox
PROJECT_NAME = hivebox
BINARY_PATH = /usr/local/go/bin/go
VERSION ?= $(shell git describe --tags --always)


LDFLAGS := -X main.Version=$(VERSION)

# The default target: when you run "make" without arguments, it will run the "build" target.
all: build

# build: Compiles the Go project into a binary executable.
build:
	@echo "Building $(BINARY_NAME)..."
	$(BINARY_PATH) build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .
	# cp $(BINARY_NAME) ~/.local/bin/
test:
	@echo "Testing $(BINARY_NAME)..."
	$(BINARY_PATH) run main.go

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
