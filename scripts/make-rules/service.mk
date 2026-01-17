# ==============================================================================
# Makefile helper functions for services
#

# Define all services
SERVICES ?= user message presence

# Base directory for proto files
PROTO_BASE_DIR := $(ROOT_DIR)/internal/api/proto

# Proto file directory for each service: internal/api/proto/{service}/v1
define get-proto-dir
$(PROTO_BASE_DIR)/$(1)/v1
endef

# Proto file path for each service: internal/api/proto/{service}/v1/{service}.proto
define get-proto-file
$(call get-proto-dir,$(1))/$(1).proto
endef

# Generated Go code directory for each service: internal/api/proto/{service}/v1
define get-generated-dir
$(call get-proto-dir,$(1))
endef

# Example usage:
# $(call get-proto-dir,user) -> internal/api/proto/user/v1
# $(call get-proto-file,user) -> internal/api/proto/user/v1/user.proto
# $(call get-generated-dir,user) -> internal/api/proto/user/v1

# Validate that a service exists
.PHONY: service.validate.%
service.validate.%:
	@if [ ! -f "$(call get-proto-file,$*)" ]; then \
		echo "Error: Proto file not found for service '$*': $(call get-proto-file,$*)"; \
		exit 1; \
	fi

# List all proto files
.PHONY: service.list
service.list:
	@echo "Available services:"
	@for svc in $(SERVICES); do \
		proto_file=$(call get-proto-file,$$svc); \
		if [ -f "$$proto_file" ]; then \
			echo "  - $$svc: $$proto_file"; \
		else \
			echo "  - $$svc: $$proto_file (NOT FOUND)"; \
		fi; \
	done