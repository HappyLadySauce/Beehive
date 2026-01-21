# ==============================================================================
# Swagger configuration

# Swag 工具路径
SWAG := $(GOPATH)/bin/swag

# 定义需要生成 swagger 的服务
SWAGGER_SERVICES := auth user

# 获取服务的 main.go 路径
define get-swagger-main
cmd/beehive-$(1)/main.go
endef

# 获取服务的 swagger 输出目录
define get-swagger-output-dir
internal/beehive-$(1)/api/swagger/docs
endef

# ==============================================================================
# Swagger generation targets

# 验证 swag 工具是否安装
.PHONY: swagger.verify
swagger.verify:
	@if [ ! -f "$(SWAG)" ]; then \
		echo "===========> Installing swag"; \
		$(MAKE) install.swagger; \
	fi

# 为每个服务生成 swagger 文档
.PHONY: swagger.auth swagger.user
swagger.auth swagger.user: swagger.verify
	@echo "===========> Generating swagger API docs for $(@:swagger.%=%)"
	@mkdir -p $(call get-swagger-output-dir,$(@:swagger.%=%))
	@$(SWAG) init -g $(call get-swagger-main,$(@:swagger.%=%)) -o $(call get-swagger-output-dir,$(@:swagger.%=%)) --exclude internal/beehive-gateway

# 生成所有服务的 swagger 文档
.PHONY: swagger
swagger: swagger.verify
	@echo "===========> Generating swagger API docs for all services"
	@for svc in $(SWAGGER_SERVICES); do \
		$(MAKE) swagger.$$svc || exit 1; \
	done
	@echo "===========> All swagger docs generated"

# ==============================================================================
# Swagger cleanup targets

# 清理单个服务的 swagger 文档
.PHONY: swagger.clean.auth swagger.clean.user
swagger.clean.auth swagger.clean.user:
	@echo "===========> Cleaning swagger docs for $(@:swagger.clean.%=%)"
	@rm -rf $(call get-swagger-output-dir,$(@:swagger.clean.%=%))

# 清理所有服务的 swagger 文档
.PHONY: swagger.clean
swagger.clean:
	@echo "===========> Cleaning swagger docs for all services"
	@for svc in $(SWAGGER_SERVICES); do \
		$(MAKE) swagger.clean.$$svc; \
	done
	@echo "===========> All swagger docs cleaned"
