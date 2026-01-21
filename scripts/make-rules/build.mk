# Makefile rules for building services

# ==================== 构建配置 ====================

# Go 构建参数
GO_BUILD_FLAGS ?= -v
GO_LDFLAGS ?= -s -w

# ==================== 本地构建目标 ====================

.PHONY: build build-all $(foreach svc,$(ALL_SERVICES),build-$(svc))

# 构建选中的服务（支持 SERVICES 变量）
build: $(foreach svc,$(SELECTED_SERVICES),build-$(svc))

# 构建所有服务
build-all: $(foreach svc,$(ALL_SERVICES),build-$(svc))

# ==================== Gateway 构建 ====================

build-beehive-gateway:
	@if [ ! -f "$(ROOT_DIR)/app/$(GATEWAY_NAME)/gateway.go" ]; then \
		echo "错误: Gateway 主文件不存在，请先运行 'make gen-gateway'"; \
		exit 1; \
	fi
	@mkdir -p $(BIN_DIR)
	$(call print_info,Building gateway...)
	@cd $(ROOT_DIR)/app/$(GATEWAY_NAME) && \
		go build $(GO_BUILD_FLAGS) -ldflags "$(GO_LDFLAGS)" -o $(BIN_DIR)/gateway gateway.go
	$(call print_success,Gateway built successfully -> $(BIN_DIR)/gateway)

# ==================== RPC 服务构建 ====================

# 定义 RPC 服务构建函数
define build_rpc_service
.PHONY: build-$(1)
build-$(1):
	@SERVICE_NAME=$(subst beehive-,,$(1)); \
	if [ ! -f "$(ROOT_DIR)/app/$(1)/$$$$SERVICE_NAME.go" ]; then \
		echo "错误: $(1) 主文件不存在，请先运行 'make gen-$(1)'"; \
		exit 1; \
	fi; \
	mkdir -p $(BIN_DIR); \
	echo "$(COLOR_BLUE)→ Building $(1)...$(COLOR_RESET)"; \
	cd $(ROOT_DIR)/app/$(1) && \
		go build $(GO_BUILD_FLAGS) -ldflags "$(GO_LDFLAGS)" -o $(BIN_DIR)/$$$$SERVICE_NAME $$$$SERVICE_NAME.go; \
	echo "$(COLOR_GREEN)✓ $(1) built successfully -> $(BIN_DIR)/$$$$SERVICE_NAME$(COLOR_RESET)"
endef

# 为每个 RPC 服务动态创建构建目标
$(foreach svc,$(RPC_SERVICES),$(eval $(call build_rpc_service,$(svc))))

# ==================== Docker 构建目标 ====================

.PHONY: docker-build docker-build-all $(foreach svc,$(ALL_SERVICES),docker-build-$(svc))

# 构建选中服务的 Docker 镜像（支持 SERVICES 变量）
docker-build: $(foreach svc,$(SELECTED_SERVICES),docker-build-$(svc))

# 构建所有服务的 Docker 镜像
docker-build-all: $(foreach svc,$(ALL_SERVICES),docker-build-$(svc))

# Gateway Docker 构建
docker-build-beehive-gateway:
	$(call print_info,Building Docker image for gateway...)
	@docker build -t $(DOCKER_IMAGE_PREFIX)/gateway:$(DOCKER_TAG) \
		--build-arg SERVICE=gateway \
		-f $(ROOT_DIR)/deploy/docker/Dockerfile.gateway \
		$(ROOT_DIR) 2>/dev/null || \
		docker build -t $(DOCKER_IMAGE_PREFIX)/gateway:$(DOCKER_TAG) \
		--build-arg SERVICE=$(GATEWAY_NAME) \
		-f $(ROOT_DIR)/Dockerfile \
		$(ROOT_DIR)
	$(call print_success,Gateway Docker image built: $(DOCKER_IMAGE_PREFIX)/gateway:$(DOCKER_TAG))

# 定义 RPC 服务 Docker 构建函数
define docker_build_rpc_service
.PHONY: docker-build-$(1)
docker-build-$(1):
	$$(call print_info,Building Docker image for $(1)...)
	@docker build -t $(DOCKER_IMAGE_PREFIX)/$(notdir $(1)):$(DOCKER_TAG) \
		--build-arg SERVICE=$(notdir $(1)) \
		-f $(ROOT_DIR)/deploy/docker/Dockerfile.$(notdir $(1)) \
		$(ROOT_DIR) 2>/dev/null || \
		docker build -t $(DOCKER_IMAGE_PREFIX)/$(notdir $(1)):$(DOCKER_TAG) \
		--build-arg SERVICE=$(1) \
		-f $(ROOT_DIR)/Dockerfile \
		$(ROOT_DIR)
	$$(call print_success,$(1) Docker image built: $(DOCKER_IMAGE_PREFIX)/$(notdir $(1)):$(DOCKER_TAG))
endef

# 为每个 RPC 服务动态创建 Docker 构建目标
$(foreach svc,$(RPC_SERVICES),$(eval $(call docker_build_rpc_service,$(svc))))

# ==================== 清理构建产物 ====================

.PHONY: clean clean-build

clean-build:
	$(call print_warning,Cleaning build artifacts...)
	@rm -rf $(BIN_DIR)
	$(call print_success,Build artifacts cleaned)

clean: clean-build

# ==================== 依赖管理 ====================

.PHONY: deps deps-download deps-tidy

# 下载依赖
deps-download:
	$(call print_info,Downloading dependencies...)
	@cd $(ROOT_DIR) && go mod download
	$(call print_success,Dependencies downloaded)

# 整理依赖
deps-tidy:
	$(call print_info,Tidying dependencies...)
	@cd $(ROOT_DIR) && go mod tidy
	$(call print_success,Dependencies tidied)

# 安装/更新依赖
deps: deps-download deps-tidy
