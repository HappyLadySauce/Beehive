# Beehive IM - Makefile
# 企业级即时通讯系统构建脚本

# ==================== 导入模块化规则 ====================

include scripts/make-rules/common.mk
include scripts/make-rules/gen.mk
include scripts/make-rules/build.mk
include scripts/make-rules/run.mk

# ==================== 默认目标 ====================

.DEFAULT_GOAL := help

# ==================== 帮助信息 ====================

.PHONY: help

help:
	@echo "=========================================="
	@echo "Beehive IM - 可用命令"
	@echo "=========================================="
	@echo ""
	@echo "📦 代码生成:"
	@echo "  make gen                           - 生成所有服务代码"
	@echo "  make gen-gateway                   - 生成 Gateway 代码"
	@echo "  make gen-<service>                 - 生成指定单个 RPC 服务"
	@echo "  make gen SERVICES=\"svc1 svc2\"      - 生成指定的多个服务"
	@echo ""
	@echo "🔨 构建:"
	@echo "  make build                         - 构建所有服务"
	@echo "  make build-<service>               - 构建指定单个服务"
	@echo "  make build SERVICES=\"svc1 svc2\"    - 构建指定的多个服务"
	@echo "  make docker-build                  - 构建所有 Docker 镜像"
	@echo "  make docker-build-<service>        - 构建指定服务的 Docker 镜像"
	@echo ""
	@echo "🚀 运行管理:"
	@echo "  make run-all                       - 后台运行所有服务"
	@echo "  make run                           - 后台运行所有服务（同 run-all）"
	@echo "  make run-<service>                 - 后台运行指定单个服务"
	@echo "  make run SERVICES=\"svc1 svc2\"      - 后台运行指定的多个服务"
	@echo ""
	@echo "  make stop-all                      - 停止所有服务"
	@echo "  make stop                          - 停止所有服务（同 stop-all）"
	@echo "  make stop-<service>                - 停止指定单个服务"
	@echo "  make stop SERVICES=\"svc1 svc2\"     - 停止指定的多个服务"
	@echo ""
	@echo "  make restart-all                   - 重启所有服务"
	@echo "  make restart                       - 重启所有服务（同 restart-all）"
	@echo "  make restart-<service>             - 重启指定单个服务"
	@echo ""
	@echo "  make status                        - 查看服务运行状态"
	@echo "  make logs-<service>                - 实时查看服务日志"
	@echo "  make logs-all                      - 查看所有服务日志（最近 20 行）"
	@echo ""
	@echo "🧹 清理:"
	@echo "  make clean                         - 清理构建产物"
	@echo "  make clean-build                   - 清理编译的二进制文件"
	@echo "  make clean-gen                     - 清理生成的代码"
	@echo "  make clean-run                     - 清理运行数据（PID、日志）"
	@echo "  make clean-all                     - 清理所有内容"
	@echo ""
	@echo "📦 依赖管理:"
	@echo "  make deps                          - 下载并整理依赖"
	@echo "  make deps-download                 - 下载依赖"
	@echo "  make deps-tidy                     - 整理依赖"
	@echo ""
	@echo "=========================================="
	@echo "🎯 可用服务列表:"
	@echo "  $(ALL_SERVICES)"
	@echo "=========================================="
	@echo ""
	@echo "💡 使用示例:"
	@echo "  # 完整工作流"
	@echo "  make gen && make build && make run-all"
	@echo ""
	@echo "  # 单服务开发"
	@echo "  make gen-beehive-user"
	@echo "  make build-beehive-user"
	@echo "  make run-beehive-user"
	@echo ""
	@echo "  # 部分服务开发（多选）"
	@echo "  make gen SERVICES=\"beehive-user beehive-friend\""
	@echo "  make build SERVICES=\"beehive-user beehive-friend\""
	@echo "  make run SERVICES=\"beehive-gateway beehive-user beehive-friend\""
	@echo ""
	@echo "  # 查看状态和日志"
	@echo "  make status"
	@echo "  make logs-beehive-user"
	@echo ""
	@echo "=========================================="

# ==================== 组合命令 ====================

.PHONY: all dev

# 生成、构建所有服务
all: gen build

# 开发模式：生成、构建、运行所有服务
dev: gen build run-all
	@echo ""
	$(call print_success,All services started in development mode)
	@echo "Run 'make status' to check service status"
	@echo "Run 'make logs-<service>' to view logs"

# ==================== 全局清理 ====================

.PHONY: clean-all

clean-all: clean-build clean-gen clean-run
	$(call print_success,All cleaned)
