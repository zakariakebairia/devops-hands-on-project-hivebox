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

dockerize:
	@echo "Building image for $(PROJECT_NAME)"
	docker build -t $(PROJECT_NAME):$(VERSION) .
drun:
	docker run -p 8080:8080 $(PROJECT_NAME):$(VERSION)

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
