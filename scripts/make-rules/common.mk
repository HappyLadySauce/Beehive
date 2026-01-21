# Linux-only: always use bash
SHELL := /bin/bash

# Define COMMON_SELF_DIR if not already defined
# This gets the directory of the current Makefile (common.mk)
ifeq ($(origin COMMON_SELF_DIR),undefined)
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
endif

# Define ROOT_DIR if not already defined
# Go up two levels from scripts/make-rules/ to get project root
ifeq ($(origin ROOT_DIR),undefined)
# Linux: use cd and pwd
ROOT_DIR := $(shell cd $(COMMON_SELF_DIR)/../.. && pwd)
endif

# Define OUTPUT_DIR if not already defined
ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(ROOT_DIR)/_output
# Linux: ensure output dir exists
$(shell mkdir -p "$(OUTPUT_DIR)")
endif

# Define BIN_DIR for compiled binaries
ifeq ($(origin BIN_DIR),undefined)
BIN_DIR := $(OUTPUT_DIR)/bin
$(shell mkdir -p "$(BIN_DIR)")
endif

# Define PID_DIR for process IDs
ifeq ($(origin PID_DIR),undefined)
PID_DIR := $(OUTPUT_DIR)/pids
$(shell mkdir -p "$(PID_DIR)")
endif

# Define LOG_DIR for service logs
ifeq ($(origin LOG_DIR),undefined)
LOG_DIR := $(OUTPUT_DIR)/logs
$(shell mkdir -p "$(LOG_DIR)")
endif

# ==================== 服务自动发现 ====================

# Gateway 服务
GATEWAY_API := api/beehive-gateway/v1/gateway.api
GATEWAY_NAME := beehive-gateway

# 自动扫描所有 RPC 服务：直接扫描目录
# 从 api/proto/beehive-xxx/ 目录中提取服务名称
RPC_SERVICES := $(sort $(notdir $(wildcard $(ROOT_DIR)/api/proto/beehive-*)))

# 所有服务列表（Gateway + RPC Services）
ALL_SERVICES := $(GATEWAY_NAME) $(RPC_SERVICES)

# 灵活服务选择机制：优先使用 SERVICES 变量，否则使用全部服务
SELECTED_SERVICES ?= $(if $(SERVICES),$(filter $(SERVICES),$(ALL_SERVICES)),$(ALL_SERVICES))

# ==================== Docker 配置 ====================

DOCKER_IMAGE_PREFIX ?= beehive
DOCKER_TAG ?= latest
DOCKER_REGISTRY ?= 

# ==================== 辅助函数 ====================

COMMA := ,
SPACE :=
SPACE +=

# 定义颜色输出
COLOR_RESET := \033[0m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m
COLOR_BLUE := \033[34m

# 打印成功信息
define print_success
	@echo "$(COLOR_GREEN)✓ $(1)$(COLOR_RESET)"
endef

# 打印警告信息
define print_warning
	@echo "$(COLOR_YELLOW)⚠ $(1)$(COLOR_RESET)"
endef

# 打印信息
define print_info
	@echo "$(COLOR_BLUE)→ $(1)$(COLOR_RESET)"
endef