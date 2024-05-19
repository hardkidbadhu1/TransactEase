working_directory := $(shell pwd)

# Path: Makefile
PROJECT := $(shell basename $(working_directory))

Run:
	@echo "Running $(PROJECT)..."
	@go run cmd/main.go

install_deps:
	@echo "Installing dependencies..."
	@go mod tidy