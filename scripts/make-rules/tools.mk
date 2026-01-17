TOOLS ?=$(BLOCKER_TOOLS) $(CRITICAL_TOOLS) $(TRIVIAL_TOOLS)

.PHONY: tools.install
tools.install: $(addprefix tools.install., $(TOOLS))

.PHONY: tools.install.%
tools.install.%:
	@echo "===========> Installing $*"
	@$(MAKE) install.$*

# Verify tool is installed (for tools in PATH)
.PHONY: tools.verify.%
tools.verify.%:
	@if ! which $* &>/dev/null; then $(MAKE) tools.install.$*; fi

# Verify Go tool is installed (for tools in GOPATH/bin)
.PHONY: tools.verify.go.%
tools.verify.go.%:
	@tool_path=$(GOPATH)/bin/$*; \
	if [ ! -f "$$tool_path" ]; then \
		$(MAKE) tools.install.$*; \
	fi

# Install rules for various tools
.PHONY: install.swagger
install.swagger:
	@$(GO) install github.com/swaggo/swag/cmd/swag@latest

.PHONY: install.protoc-gen-go
install.protoc-gen-go:
	@$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest

.PHONY: install.protoc-gen-go-grpc
install.protoc-gen-go-grpc:
	@$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest