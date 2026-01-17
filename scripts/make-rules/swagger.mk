# Try to find swag binary: prefer GOROOT (if set), otherwise fall back to PATH.
# This avoids toolchain/stdlib mismatches when the environment sets GOROOT.
SWAG := $(shell \
	if [ -n "$$GOROOT" ] && [ -f "$$GOROOT/bin/swag" ]; then \
		echo $$GOROOT/bin/swag; \
	elif command -v swag >/dev/null 2>&1; then \
		echo swag; \
	else \
		echo swag; \
	fi)

.PHONY: swagger
swagger: tools.verify.swagger
	@echo "===========> Generating swagger API docs"
	@$(SWAG) init -g cmd/main.go -o api/swagger/docs
