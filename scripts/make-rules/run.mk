# Makefile rules for running services

# ==================== 运行管理目标 ====================

.PHONY: run run-all $(foreach svc,$(ALL_SERVICES),run-$(svc))
.PHONY: stop stop-all $(foreach svc,$(ALL_SERVICES),stop-$(svc))
.PHONY: restart restart-all $(foreach svc,$(ALL_SERVICES),restart-$(svc))
.PHONY: status logs-all $(foreach svc,$(ALL_SERVICES),logs-$(svc))

# 运行选中的服务（支持 SERVICES 变量）
run: $(foreach svc,$(SELECTED_SERVICES),run-$(svc))

# 运行所有服务
run-all: $(foreach svc,$(ALL_SERVICES),run-$(svc))

# 停止选中的服务（支持 SERVICES 变量）
stop: $(foreach svc,$(SELECTED_SERVICES),stop-$(svc))

# 停止所有服务
stop-all: $(foreach svc,$(ALL_SERVICES),stop-$(svc))

# 重启选中的服务（支持 SERVICES 变量）
restart: $(foreach svc,$(SELECTED_SERVICES),restart-$(svc))

# 重启所有服务
restart-all: $(foreach svc,$(ALL_SERVICES),restart-$(svc))

# ==================== Gateway 运行管理 ====================

run-beehive-gateway:
	@mkdir -p $(PID_DIR) $(LOG_DIR)
	@if [ ! -f "$(BIN_DIR)/gateway" ]; then \
		echo "错误: Gateway 可执行文件不存在，请先运行 'make build-gateway'"; \
		exit 1; \
	fi
	@if [ -f $(PID_DIR)/gateway.pid ] && kill -0 $$(cat $(PID_DIR)/gateway.pid) 2>/dev/null; then \
		echo "Gateway 已经在运行 (PID: $$(cat $(PID_DIR)/gateway.pid))"; \
	else \
		nohup $(BIN_DIR)/gateway -f $(ROOT_DIR)/app/$(GATEWAY_NAME)/etc/gateway-api.yaml \
			> $(LOG_DIR)/gateway.log 2>&1 & echo $$! > $(PID_DIR)/gateway.pid; \
		sleep 1; \
		if kill -0 $$(cat $(PID_DIR)/gateway.pid) 2>/dev/null; then \
			echo "$(COLOR_GREEN)✓ Gateway started (PID: $$(cat $(PID_DIR)/gateway.pid))$(COLOR_RESET)"; \
		else \
			echo "$(COLOR_YELLOW)⚠ Gateway failed to start, check $(LOG_DIR)/gateway.log$(COLOR_RESET)"; \
			rm -f $(PID_DIR)/gateway.pid; \
			exit 1; \
		fi; \
	fi

stop-beehive-gateway:
	@if [ -f $(PID_DIR)/gateway.pid ]; then \
		if kill -0 $$(cat $(PID_DIR)/gateway.pid) 2>/dev/null; then \
			kill $$(cat $(PID_DIR)/gateway.pid) && rm $(PID_DIR)/gateway.pid; \
			echo "$(COLOR_GREEN)✓ Gateway stopped$(COLOR_RESET)"; \
		else \
			rm $(PID_DIR)/gateway.pid; \
			echo "$(COLOR_YELLOW)⚠ Gateway was not running$(COLOR_RESET)"; \
		fi; \
	else \
		echo "$(COLOR_YELLOW)⚠ Gateway is not running$(COLOR_RESET)"; \
	fi

restart-beehive-gateway: stop-beehive-gateway run-beehive-gateway

logs-beehive-gateway:
	@if [ ! -f "$(LOG_DIR)/gateway.log" ]; then \
		echo "错误: 日志文件不存在: $(LOG_DIR)/gateway.log"; \
		exit 1; \
	fi
	@tail -f $(LOG_DIR)/gateway.log

# ==================== RPC 服务运行管理 ====================

# 定义 RPC 服务运行函数
define run_rpc_service
.PHONY: run-$(1) stop-$(1) restart-$(1) logs-$(1)

run-$(1):
	@SERVICE_NAME=$$(echo "$(1)" | sed 's/beehive-//'); \
	mkdir -p $(PID_DIR) $(LOG_DIR); \
	if [ ! -f "$(BIN_DIR)/$$SERVICE_NAME" ]; then \
		echo "错误: $(1) 可执行文件不存在，请先运行 'make build-$(1)'"; \
		exit 1; \
	fi; \
	if [ -f $(PID_DIR)/$$SERVICE_NAME.pid ] && kill -0 $$$$(cat $(PID_DIR)/$$SERVICE_NAME.pid) 2>/dev/null; then \
		echo "$(1) 已经在运行 (PID: $$$$(cat $(PID_DIR)/$$SERVICE_NAME.pid))"; \
	else \
		nohup $(BIN_DIR)/$$SERVICE_NAME -f $(ROOT_DIR)/app/$(1)/etc/$$SERVICE_NAME.yaml \
			> $(LOG_DIR)/$$SERVICE_NAME.log 2>&1 & echo $$$$! > $(PID_DIR)/$$SERVICE_NAME.pid; \
		sleep 1; \
		if kill -0 $$$$(cat $(PID_DIR)/$$SERVICE_NAME.pid) 2>/dev/null; then \
			echo "$(COLOR_GREEN)✓ $(1) started (PID: $$$$(cat $(PID_DIR)/$$SERVICE_NAME.pid))$(COLOR_RESET)"; \
		else \
			echo "$(COLOR_YELLOW)⚠ $(1) failed to start, check $(LOG_DIR)/$$SERVICE_NAME.log$(COLOR_RESET)"; \
			rm -f $(PID_DIR)/$$SERVICE_NAME.pid; \
			exit 1; \
		fi; \
	fi

stop-$(1):
	@SERVICE_NAME=$$(echo "$(1)" | sed 's/beehive-//'); \
	if [ -f $(PID_DIR)/$$SERVICE_NAME.pid ]; then \
		if kill -0 $$$$(cat $(PID_DIR)/$$SERVICE_NAME.pid) 2>/dev/null; then \
			kill $$$$(cat $(PID_DIR)/$$SERVICE_NAME.pid) && rm $(PID_DIR)/$$SERVICE_NAME.pid; \
			echo "$(COLOR_GREEN)✓ $(1) stopped$(COLOR_RESET)"; \
		else \
			rm $(PID_DIR)/$$SERVICE_NAME.pid; \
			echo "$(COLOR_YELLOW)⚠ $(1) was not running$(COLOR_RESET)"; \
		fi; \
	else \
		echo "$(COLOR_YELLOW)⚠ $(1) is not running$(COLOR_RESET)"; \
	fi

restart-$(1): stop-$(1) run-$(1)

logs-$(1):
	@SERVICE_NAME=$$(echo "$(1)" | sed 's/beehive-//'); \
	if [ ! -f "$(LOG_DIR)/$$SERVICE_NAME.log" ]; then \
		echo "错误: 日志文件不存在: $(LOG_DIR)/$$SERVICE_NAME.log"; \
		exit 1; \
	fi; \
	tail -f $(LOG_DIR)/$$SERVICE_NAME.log

endef

# 为每个 RPC 服务动态创建运行管理目标
$(foreach svc,$(RPC_SERVICES),$(eval $(call run_rpc_service,$(svc))))

# ==================== 服务状态查看 ====================

status:
	@echo "=========================================="
	@echo "Beehive 服务运行状态"
	@echo "=========================================="
	@for svc in $(ALL_SERVICES); do \
		if [ "$$svc" = "$(GATEWAY_NAME)" ]; then \
			bin_name="gateway"; \
		else \
			bin_name=$$(echo "$$svc" | sed 's/beehive-//'); \
		fi; \
		if [ -f $(PID_DIR)/$$bin_name.pid ]; then \
			pid=$$(cat $(PID_DIR)/$$bin_name.pid); \
			if kill -0 $$pid 2>/dev/null; then \
				echo "✓ $$svc: 运行中 (PID: $$pid)"; \
			else \
				echo "✗ $$svc: 已停止 (PID 文件存在但进程不存在)"; \
			fi; \
		else \
			echo "✗ $$svc: 未运行"; \
		fi; \
	done
	@echo "=========================================="
	@echo "日志目录: $(LOG_DIR)"
	@echo "PID 目录: $(PID_DIR)"
	@echo "=========================================="

# 查看所有服务日志
logs-all:
	@for svc in $(ALL_SERVICES); do \
		if [ "$$svc" = "$(GATEWAY_NAME)" ]; then \
			bin_name="gateway"; \
		else \
			bin_name=$$(echo "$$svc" | sed 's/beehive-//'); \
		fi; \
		if [ -f "$(LOG_DIR)/$$bin_name.log" ]; then \
			echo "======== $$svc ========"; \
			tail -n 20 $(LOG_DIR)/$$bin_name.log; \
			echo ""; \
		fi; \
	done

# ==================== 清理运行数据 ====================

.PHONY: clean-run

clean-run: stop-all
	$(call print_warning,Cleaning runtime data...)
	@rm -rf $(PID_DIR) $(LOG_DIR)
	$(call print_success,Runtime data cleaned)
