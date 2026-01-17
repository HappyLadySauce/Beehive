# Makefile for Beehive

# Build all by default, even if it's not first
.DEFAULT_GOAL := all

# Silence "Entering/Leaving directory ..." output from recursive make invocations.
MAKEFLAGS += --no-print-directory

.PHONY: all
all: build

# ==============================================================================
# Build options

ROOT_PACKAGE=github.com/HappyLadySauce/Beehive

# ==============================================================================
# Includes

include scripts/make-rules/common.mk
include scripts/make-rules/golang.mk
include scripts/make-rules/service.mk
include scripts/make-rules/grpc.mk

# ==============================================================================
# Usage

USAGE_OPTIONS :=
USAGE_OPTIONS += "go.version: Show Go version"
USAGE_OPTIONS += "go.build: Build binary"
USAGE_OPTIONS += "go.run: Run binary"
USAGE_OPTIONS += "go.tidy: Tidy Go modules"
USAGE_OPTIONS += "grpc.gen: Generate gRPC code for all services"
USAGE_OPTIONS += "grpc.gen.SERVICE: Generate gRPC code for a specific service (e.g., grpc.gen.user)"
USAGE_OPTIONS += "grpc.clean: Clean generated gRPC code for all services"
USAGE_OPTIONS += "grpc.status: Show gRPC generation status"
USAGE_OPTIONS += "service.list: List all available services"
USAGE_OPTIONS += "help: Show this help info"

export USAGE_OPTIONS

# ==============================================================================
# Targets

.PHONY: env
env:
	@$(GO) version

.PHONY: tidy
tidy:
	@$(GO) mod tidy

## help: Show this help info.
.PHONY: help
help: Makefile
	@printf "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:\n"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"