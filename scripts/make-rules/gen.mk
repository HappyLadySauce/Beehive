# Makefile rules for code generation using goctl

# ==================== 代码生成目标 ====================

.PHONY: gen gen-gateway gen-all $(foreach svc,$(RPC_SERVICES),gen-$(svc))

# 生成选中的服务代码（支持 SERVICES 变量）
gen: $(foreach svc,$(SELECTED_SERVICES),gen-$(svc))

# 生成所有服务代码
gen-all: $(foreach svc,$(ALL_SERVICES),gen-$(svc))

# ==================== Gateway 代码生成 ====================

gen-beehive-gateway:
	@if [ ! -f "$(ROOT_DIR)/$(GATEWAY_API)" ]; then \
		echo "错误: Gateway API 文件不存在: $(GATEWAY_API)"; \
		exit 1; \
	fi
	$(call print_info,Generating API Gateway code...)
	@goctl api go --api $(ROOT_DIR)/$(GATEWAY_API) --dir $(ROOT_DIR)/app/$(GATEWAY_NAME)/
	$(call print_success,Gateway code generated)

# ==================== RPC 服务代码生成 ====================

# 定义 RPC 服务生成函数
define gen_rpc_service
.PHONY: gen-$(1)
gen-$(1):
	@if [ ! -f "$(ROOT_DIR)/api/proto/$(1)/v1/"*.proto ]; then \
		echo "错误: RPC proto 文件不存在: api/proto/$(1)/v1/*.proto"; \
		exit 1; \
	fi
	$$(call print_info,Generating RPC service: $(1)...)
	@cd $(ROOT_DIR) && goctl rpc protoc api/proto/$(1)/v1/*.proto \
		--go_out=app/$(1)/ \
		--go-grpc_out=app/$(1)/ \
		--zrpc_out=app/$(1)/
	$$(call print_success,$(1) code generated)
endef

# 为每个 RPC 服务动态创建生成目标
$(foreach svc,$(RPC_SERVICES),$(eval $(call gen_rpc_service,$(svc))))

# ==================== 清理生成的代码 ====================

.PHONY: clean-gen

clean-gen:
	$(call print_warning,Cleaning generated code...)
	@rm -rf $(ROOT_DIR)/app/*/internal/pb
	@rm -rf $(ROOT_DIR)/app/*/types
	@find $(ROOT_DIR)/app -name "*_grpc.pb.go" -delete
	@find $(ROOT_DIR)/app -name "*.pb.go" -delete
	$(call print_success,Generated code cleaned)
