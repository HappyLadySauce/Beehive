# Beehive IM - Makefile 使用指南

## 概述

本项目使用模块化的 Makefile 系统，支持自动服务发现、灵活的服务选择（单服务、部分服务、全部服务），以及代码生成、构建、运行管理等功能。

## 目录结构

```
Makefile                          # 主 Makefile 入口
scripts/make-rules/
  ├── common.mk                   # 通用变量、函数、服务发现
  ├── gen.mk                      # 代码生成规则
  ├── build.mk                    # 构建规则（本地 + Docker）
  └── run.mk                      # 运行管理规则
```

## 核心特性

### 1. 自动服务发现

系统自动扫描 `api/` 目录，发现所有服务：
- **Gateway**: `api/beehive-gateway/v1/gateway.api`
- **RPC 服务**: `api/proto/beehive-*/v1/*.proto`

当前自动识别的服务：
- `beehive-gateway`
- `beehive-user`
- `beehive-friend`
- `beehive-chat`
- `beehive-message`
- `beehive-file`
- `beehive-search`

### 2. 灵活的服务选择

支持三种操作模式：

#### a) 全部服务（默认）
```bash
make gen          # 生成所有服务代码
make build        # 构建所有服务
make run-all      # 运行所有服务
```

#### b) 单个服务
```bash
make gen-beehive-user
make build-beehive-user
make run-beehive-user
```

#### c) 部分服务（多选）
```bash
# 只开发用户和好友模块
make gen SERVICES="beehive-user beehive-friend"
make build SERVICES="beehive-user beehive-friend"
make run SERVICES="beehive-gateway beehive-user beehive-friend"
```

## 常用命令

### 代码生成

```bash
# 生成所有服务代码
make gen

# 生成 Gateway 代码
make gen-gateway

# 生成指定 RPC 服务代码
make gen-beehive-user
make gen-beehive-friend

# 生成部分服务代码
make gen SERVICES="beehive-user beehive-friend beehive-chat"
```

### 构建

```bash
# 构建所有服务
make build

# 构建单个服务
make build-beehive-user

# 构建部分服务
make build SERVICES="beehive-user beehive-friend"

# 构建 Docker 镜像
make docker-build
make docker-build-beehive-user
```

### 运行管理

```bash
# 运行所有服务（后台）
make run-all

# 运行单个服务
make run-beehive-gateway
make run-beehive-user

# 运行部分服务
make run SERVICES="beehive-gateway beehive-user"

# 停止服务
make stop-all                    # 停止所有服务
make stop-beehive-user          # 停止单个服务
make stop SERVICES="..."        # 停止部分服务

# 重启服务
make restart-all
make restart-beehive-user

# 查看服务状态
make status

# 查看日志
make logs-beehive-user          # 实时查看单个服务日志
make logs-all                   # 查看所有服务最近日志
```

### 清理

```bash
make clean              # 清理构建产物
make clean-build        # 只清理二进制文件
make clean-gen          # 只清理生成的代码
make clean-run          # 清理运行数据（PID、日志）
make clean-all          # 清理所有内容
```

### 依赖管理

```bash
make deps               # 下载并整理依赖
make deps-download      # 只下载依赖
make deps-tidy          # 只整理依赖
```

## 完整工作流示例

### 场景 1: 完整开发所有服务

```bash
# 1. 生成代码
make gen

# 2. 构建服务
make build

# 3. 运行所有服务
make run-all

# 4. 查看状态
make status

# 5. 查看日志
make logs-beehive-gateway

# 6. 停止服务
make stop-all
```

### 场景 2: 单个服务开发

```bash
# 只开发用户服务
make gen-beehive-user
make build-beehive-user
make run-beehive-user

# 查看日志
make logs-beehive-user

# 停止服务
make stop-beehive-user
```

### 场景 3: 模块化开发

```bash
# 只开发用户和好友模块
make gen SERVICES="beehive-user beehive-friend"
make build SERVICES="beehive-user beehive-friend"

# 运行核心服务
make run SERVICES="beehive-gateway beehive-user beehive-friend"

# 查看状态
make status

# 重启部分服务
make restart SERVICES="beehive-user beehive-friend"
```

### 场景 4: 快速启动开发环境

```bash
# 一条命令：生成、构建、运行
make dev

# 等价于
make gen && make build && make run-all
```

## 输出目录结构

所有生成的文件都存放在 `_output/` 目录：

```
_output/
  ├── bin/                      # 编译的二进制文件
  │   ├── gateway
  │   ├── beehive-user
  │   ├── beehive-friend
  │   └── ...
  ├── logs/                     # 服务日志
  │   ├── gateway.log
  │   ├── beehive-user.log
  │   └── ...
  └── pids/                     # 进程 PID 文件
      ├── gateway.pid
      ├── beehive-user.pid
      └── ...
```

## 添加新服务

要添加新服务，只需：

1. **添加 API 定义**：
   - Gateway: 在 `api/beehive-gateway/v1/gateway.api` 中添加路由
   - RPC: 在 `api/proto/beehive-newservice/v1/newservice.proto` 创建 proto 文件

2. **生成代码**：
   ```bash
   make gen-beehive-newservice
   ```

3. **构建并运行**：
   ```bash
   make build-beehive-newservice
   make run-beehive-newservice
   ```

无需修改 Makefile，系统会自动识别新服务！

## 环境变量

可以通过环境变量自定义配置：

```bash
# 指定操作的服务列表
export SERVICES="beehive-user beehive-friend"
make build

# 自定义 Docker 镜像前缀
export DOCKER_IMAGE_PREFIX="myregistry/beehive"
export DOCKER_TAG="v1.0.0"
make docker-build

# 自定义输出目录
export OUTPUT_DIR="/tmp/beehive-output"
make build
```

## 故障排查

### 服务无法启动

```bash
# 1. 检查服务状态
make status

# 2. 查看日志
make logs-beehive-user

# 3. 重新构建
make build-beehive-user

# 4. 重新运行
make run-beehive-user
```

### 清理并重新开始

```bash
# 清理所有内容
make clean-all

# 重新生成和构建
make gen && make build
```

## 最佳实践

1. **开发单个服务时**：使用 `make <action>-<service>` 命令，提高效率
2. **开发相关模块时**：使用 `SERVICES` 变量指定多个服务
3. **完整测试时**：使用 `make run-all` 启动所有服务
4. **查看日志**：使用 `make logs-<service>` 实时监控服务状态
5. **定期清理**：使用 `make clean` 清理构建产物

## 获取帮助

```bash
# 查看所有可用命令
make help

# 或直接运行
make
```

---

**提示**: 本 Makefile 系统支持 Tab 补全（如果你的 shell 配置了 make 补全）。
