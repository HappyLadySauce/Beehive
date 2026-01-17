# ==============================================================================
# Makefile helper functions for golang
#

# Try to find go binary: prefer GOROOT (if set), otherwise fall back to PATH.
# This avoids toolchain/stdlib mismatches when the environment sets GOROOT.
GO := $(shell \
	if [ -n "$$GOROOT" ] && [ -f "$$GOROOT/bin/go" ]; then \
		echo $$GOROOT/bin/go; \
	elif command -v go >/dev/null 2>&1; then \
		echo go; \
	else \
		echo go; \
	fi)

# Get GOPATH from go env, fallback to environment variable or default
GOPATH := $(shell $(GO) env GOPATH 2>/dev/null || echo $$GOPATH || echo $(HOME)/go)

# Get GOROOT from go env, fallback to environment variable
GOROOT := $(shell $(GO) env GOROOT 2>/dev/null || echo $$GOROOT)

# ==============================================================================
# Version flags
#
# BUILD_VERSION 用于控制编译注入到二进制中的版本号:
# - 不设置时，默认使用 pkg/environment/version.go 中的默认值（dev）
# - release 构建时，通过 Makefile 显式传入 BUILD_VERSION 覆盖
# 版本常量定义在 pkg/environment/version.go 中:
#   const (
#     release = "1.0.0"
#     dev     = "1.0.0-dev"
#   )
#
BUILD_VERSION ?=

LDFLAGS_VERSION :=
ifneq ($(strip $(BUILD_VERSION)),)
LDFLAGS_VERSION += -X github.com/HappyLadySauce/NexusPointWG/pkg/environment.Version=$(BUILD_VERSION)
endif

LDFLAGS := $(strip $(LDFLAGS_VERSION))

.PHONY: go.build
go.build:
	@echo "===========> Building binary (static, no CGO)"
	@if [ -z "$(BINARY_NAME)" ]; then \
		echo "Error: BINARY_NAME is not set. Please set it in Makefile."; \
		exit 1; \
	fi
	@CGO_ENABLED=0 GOOS=linux $(GO) build -a -installsuffix cgo -ldflags '$(LDFLAGS)' -o $(OUTPUT_DIR)/$(BINARY_NAME) cmd/main.go

.PHONY: go.run
go.run:
	@echo "===========> Running NexusPointWG"
	@mkdir -p "$(OUTPUT_DIR)" "$(dir $(NEXUS_POINT_WG_LOGS_LOG_FILE))"
	@$(GO) run cmd/main.go -c $(CONFIG_FILE)

