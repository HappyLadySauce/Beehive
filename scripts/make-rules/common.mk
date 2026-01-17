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

# Specify tools severity, include: BLOCKER_TOOLS, CRITICAL_TOOLS, TRIVIAL_TOOLS.
# Missing BLOCKER_TOOLS can cause the CI flow execution failed, i.e. `make all` failed.
# Missing CRITICAL_TOOLS can lead to some necessary operations failed. i.e. `make release` failed.
# TRIVIAL_TOOLS are Optional tools, missing these tool have no affect.
BLOCKER_TOOLS ?= gsemver golines go-junit-report golangci-lint addlicense goimports codegen
CRITICAL_TOOLS ?= swagger mockgen gotests git-chglog github-release coscmd go-mod-outdated protoc-gen-go protoc-gen-go-grpc cfssl go-gitlint
TRIVIAL_TOOLS ?= depth go-callvis gothanks richgo rts kube-score

COMMA := ,
SPACE :=
SPACE +=