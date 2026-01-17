# ==============================================================================
# Makefile helper functions for grpc
#
# Note: service.mk should be included in the main Makefile, not here
# to avoid duplicate definitions

# Protobuf compiler
PROTOC ?= protoc

# Go protobuf plugin paths
PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(GOPATH)/bin/protoc-gen-go-grpc

# Verify protoc is installed (system tool, not in GOPATH)
.PHONY: grpc.verify.protoc
grpc.verify.protoc:
	@if ! command -v $(PROTOC) >/dev/null 2>&1; then \
		echo "Error: protoc is not installed. Please install it:"; \
		echo "  macOS: brew install protobuf"; \
		echo "  Linux: sudo apt-get install protobuf-compiler"; \
		exit 1; \
	fi

# Verify all grpc tools
# Uses tools.mk for Go tools (protoc-gen-go, protoc-gen-go-grpc)
# Special handling for protoc (system tool)
.PHONY: grpc.verify
grpc.verify: grpc.verify.protoc grpc.verify.go-tools

# Verify Go tools using tools.mk mechanism
.PHONY: grpc.verify.go-tools
grpc.verify.go-tools:
	@if [ ! -f "$(PROTOC_GEN_GO)" ]; then \
		$(MAKE) tools.install.protoc-gen-go; \
	fi
	@if [ ! -f "$(PROTOC_GEN_GO_GRPC)" ]; then \
		$(MAKE) tools.install.protoc-gen-go-grpc; \
	fi

# Generate gRPC code for a specific service
# Usage: make grpc.gen.user
.PHONY: grpc.gen.%
grpc.gen.%: grpc.verify service.validate.%
	@echo "===========> Generating gRPC code for service: $*"
	@proto_dir=$(call get-proto-dir,$*); \
	proto_file=$(call get-proto-file,$*); \
	output_dir=$$proto_dir; \
	if [ ! -f "$$proto_file" ]; then \
		echo "Error: Proto file not found: $$proto_file"; \
		exit 1; \
	fi; \
	echo "Proto file: $$proto_file"; \
	echo "Output directory: $$output_dir"; \
	$(PROTOC) \
		--proto_path=$$proto_dir \
		--proto_path=$(ROOT_DIR) \
		--go_out=$$output_dir \
		--go_opt=paths=source_relative \
		--go-grpc_out=$$output_dir \
		--go-grpc_opt=paths=source_relative \
		$$proto_file; \
	if [ $$? -eq 0 ]; then \
		echo "Successfully generated gRPC code for $*"; \
	else \
		echo "Error: Failed to generate gRPC code for $*"; \
		exit 1; \
	fi

# Generate gRPC code for all services
.PHONY: grpc.gen
grpc.gen: grpc.verify
	@echo "===========> Generating gRPC code for all services"
	@for svc in $(SERVICES); do \
		$(MAKE) grpc.gen.$$svc || exit 1; \
	done
	@echo "===========> All gRPC code generated successfully"

# Clean generated gRPC code for a specific service
.PHONY: grpc.clean.%
grpc.clean.%:
	@echo "===========> Cleaning generated gRPC code for service: $*"
	@generated_dir=$(call get-generated-dir,$*); \
	find $$generated_dir -name "*.pb.go" -type f -delete; \
	echo "Cleaned generated files in $$generated_dir"

# Clean generated gRPC code for all services
.PHONY: grpc.clean
grpc.clean:
	@echo "===========> Cleaning generated gRPC code for all services"
	@for svc in $(SERVICES); do \
		$(MAKE) grpc.clean.$$svc; \
	done

# Show gRPC generation status
.PHONY: grpc.status
grpc.status:
	@echo "gRPC Generation Status:"
	@echo "======================"
	@for svc in $(SERVICES); do \
		proto_file=$(call get-proto-file,$$svc); \
		generated_dir=$(call get-generated-dir,$$svc); \
		pb_go=$$generated_dir/$$svc.pb.go; \
		grpc_go=$$generated_dir/$${svc}_grpc.pb.go; \
		echo ""; \
		echo "Service: $$svc"; \
		if [ -f "$$proto_file" ]; then \
			echo "  Proto file: ✓ $$proto_file"; \
		else \
			echo "  Proto file: ✗ $$proto_file (NOT FOUND)"; \
		fi; \
		if [ -f "$$pb_go" ]; then \
			echo "  Generated .pb.go: ✓ $$pb_go"; \
		else \
			echo "  Generated .pb.go: ✗ $$pb_go (NOT FOUND)"; \
		fi; \
		if [ -f "$$grpc_go" ]; then \
			echo "  Generated _grpc.pb.go: ✓ $$grpc_go"; \
		else \
			echo "  Generated _grpc.pb.go: ✗ $$grpc_go (NOT FOUND)"; \
		fi; \
	done