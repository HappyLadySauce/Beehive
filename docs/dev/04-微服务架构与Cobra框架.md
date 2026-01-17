# 微服务架构与 Cobra 框架

## 概述

本文档详细说明 Beehive IM 系统的微服务架构设计，以及如何使用 Cobra 命令行框架构建和管理各个微服务。

## 微服务架构

### 架构图

```mermaid
graph TB
    subgraph Clients["客户端层"]
        Web[Web客户端]
        Mobile[移动端]
        API[API客户端]
    end
    
    subgraph Gateway["Gateway 服务<br/>无状态，可水平扩展"]
        WS[WebSocket Gateway<br/>:8080]
        HTTP[HTTP API<br/>:8080]
    end
    
    subgraph Services["业务服务层<br/>gRPC 服务"]
        US[User Service<br/>:50051]
        MS[Message Service<br/>:50052]
        PS[Presence Service<br/>:50053]
    end
    
    subgraph Infrastructure["基础设施层"]
        MQ[RabbitMQ<br/>:5672]
        DB[(PostgreSQL<br/>:5432)]
        Redis[(Redis<br/>:6379)]
    end
    
    Web -->|WebSocket| WS
    Mobile -->|WebSocket| WS
    API -->|HTTP| HTTP
    
    WS -->|gRPC| US
    WS -->|gRPC| MS
    WS -->|gRPC| PS
    
    MS -->|Producer| MQ
    WS -->|Consumer| MQ
    
    US --> DB
    MS --> DB
    PS --> Redis
    PS --> DB
```

### 微服务划分

#### 1. Gateway 服务

**职责**：
- WebSocket 连接管理
- 客户端消息接收和路由
- 消息推送（从 RabbitMQ 消费）
- HTTP API（健康检查、管理接口）

**技术栈**：
- gorilla/websocket
- gRPC Client
- RabbitMQ Consumer

**端口**：8080

**特点**：
- 无状态服务
- 支持水平扩展
- 每个实例独立管理连接

**数据**：无持久化数据，仅内存连接管理

#### 2. User Service

**职责**：
- 用户注册
- 用户登录和认证
- JWT Token 生成和验证
- 用户信息查询和更新

**技术栈**：
- gRPC Server
- GORM
- JWT

**端口**：50051

**数据表**：
- `users`

**特点**：
- 有状态服务（数据库）
- 支持数据库迁移

#### 3. Message Service

**职责**：
- 单聊消息发送
- 群聊消息发送
- 消息历史查询
- 未读消息管理
- 消息状态更新

**技术栈**：
- gRPC Server
- GORM
- RabbitMQ Producer

**端口**：50052

**数据表**：
- `messages`
- `conversations`
- `groups`
- `group_members`

**特点**：
- 有状态服务（数据库）
- 消息队列生产者
- 支持数据库迁移

#### 4. Presence Service

**职责**：
- 用户在线状态管理
- 用户上线/下线通知
- 在线用户查询

**技术栈**：
- gRPC Server
- Redis
- RabbitMQ Producer

**端口**：50053

**数据**：
- Redis（在线状态缓存）
- 数据库（可选，用于持久化）

**特点**：
- 轻量级服务
- 主要依赖 Redis

## Cobra 框架使用

### 为什么使用 Cobra？

1. **统一的命令行接口**：所有微服务使用相同的命令结构
2. **配置管理**：统一的配置加载和验证机制
3. **扩展性**：易于添加新的命令和子命令
4. **开发体验**：清晰的帮助信息和错误提示

### 命令结构

每个微服务都遵循以下命令结构：

```
{service-name}
├── serve          # 启动服务
├── migrate        # 数据库迁移（仅 User/Message Service）
├── version        # 显示版本信息
└── config         # 配置相关命令
    ├── validate   # 验证配置文件
    └── show       # 显示当前配置
```

### 项目结构

```
cmd/
├── gateway/
│   ├── main.go              # 入口文件
│   └── cmd/
│       ├── root.go          # 根命令
│       ├── serve.go         # serve 命令
│       └── version.go       # version 命令
├── user-service/
│   ├── main.go
│   └── cmd/
│       ├── root.go
│       ├── serve.go
│       ├── migrate.go       # migrate 命令
│       └── version.go
├── message-service/
│   └── ...
└── presence-service/
    └── ...
```

### 实现示例

#### 1. 根命令（Root Command）

```go
// cmd/{service}/cmd/root.go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/HappyLadySauce/Beehive/pkg/options"
    "github.com/HappyLadySauce/Beehive/internal/pkg/config"
)

var (
    cfg *config.Config
)

func New{Service}Command() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "{service-name}",
        Short: "{Service} for Beehive IM",
        Long:  `{Service} handles {service description}.`,
        PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
            var err error
            cfg, err = config.LoadConfig(".")
            if err != nil {
                return fmt.Errorf("failed to load config: %w", err)
            }
            return nil
        },
    }
    
    // 添加全局配置标志
    options.AddConfigFlag("{service-name}", cmd.PersistentFlags())
    
    // 添加子命令
    cmd.AddCommand(NewServeCommand())
    cmd.AddCommand(NewVersionCommand())
    
    // 如果有数据库，添加 migrate 命令
    // cmd.AddCommand(NewMigrateCommand())
    
    return cmd
}
```

#### 2. Serve 命令

```go
// cmd/{service}/cmd/serve.go
package cmd

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "github.com/spf13/cobra"
    // 导入服务相关包
)

func NewServeCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "serve",
        Short: "Start the {service}",
        Long:  `Start the {service} server.`,
        RunE: func(cmd *cobra.Command, args []string) error {
            return runServe()
        },
    }
    
    // 可以添加服务特定的标志
    // cmd.Flags().IntVarP(&port, "port", "p", 0, "Server port")
    
    return cmd
}

func runServe() error {
    // 初始化服务依赖
    // 1. 数据库连接（如果需要）
    // 2. Redis 连接（如果需要）
    // 3. RabbitMQ 连接（如果需要）
    // 4. 创建服务实例
    // 5. 启动服务器
    
    // 优雅关闭
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        <-sigChan
        log.Println("Shutting down server...")
        // 执行清理操作
        cancel()
    }()
    
    // 启动服务
    // ...
    
    <-ctx.Done()
    return nil
}
```

#### 3. Migrate 命令（仅 User/Message Service）

```go
// cmd/{service}/cmd/migrate.go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/HappyLadySauce/Beehive/internal/pkg/db"
)

func NewMigrateCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "migrate",
        Short: "Run database migrations",
        Long:  `Run database migrations for {service}.`,
        RunE: func(cmd *cobra.Command, args []string) error {
            if err := db.InitDB(&cfg.Database); err != nil {
                return fmt.Errorf("failed to init database: %w", err)
            }
            
            if err := db.AutoMigrate(); err != nil {
                return fmt.Errorf("failed to migrate: %w", err)
            }
            
            fmt.Println("Database migration completed successfully")
            return nil
        },
    }
    
    return cmd
}
```

#### 4. Version 命令

```go
// cmd/{service}/cmd/version.go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var (
    Version   = "dev"
    BuildTime = "unknown"
    GitCommit = "unknown"
)

func NewVersionCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "version",
        Short: "Print version information",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Printf("{Service}\n")
            fmt.Printf("Version: %s\n", Version)
            fmt.Printf("Build Time: %s\n", BuildTime)
            fmt.Printf("Git Commit: %s\n", GitCommit)
        },
    }
    
    return cmd
}
```

#### 5. 主入口文件

```go
// cmd/{service}/main.go
package main

import (
    "os"
    "github.com/HappyLadySauce/Beehive/cmd/{service}/cmd"
)

func main() {
    command := cmd.New{Service}Command()
    if err := command.Execute(); err != nil {
        os.Exit(1)
    }
}
```

### 配置管理

#### 配置文件优先级

1. 命令行参数 `--config` 指定的文件
2. 环境变量 `{SERVICE}_CONFIG`
3. 当前目录的 `{service-name}.yaml`
4. `~/.beehive/{service-name}.yaml`
5. `/etc/beehive/{service-name}.yaml`

#### 环境变量支持

所有配置项都支持环境变量，格式：`{SERVICE}_{SECTION}_{KEY}`

**示例**：

```bash
# User Service
export USER_SERVICE_SERVER_PORT=50051
export USER_SERVICE_DATABASE_HOST=localhost
export USER_SERVICE_DATABASE_PORT=5432
export USER_SERVICE_DATABASE_USER=postgres
export USER_SERVICE_DATABASE_PASSWORD=postgres
export USER_SERVICE_DATABASE_DBNAME=beehive

# Gateway
export GATEWAY_SERVER_GATEWAY_PORT=8080
export GATEWAY_RABBITMQ_URL=amqp://guest:guest@localhost:5672/
```

#### 配置文件示例

```yaml
# configs/gateway.yaml
server:
  gateway_port: 8080

rabbitmq:
  url: amqp://guest:guest@localhost:5672/

grpc:
  user_service_addr: localhost:50051
  message_service_addr: localhost:50052
  presence_service_addr: localhost:50053

log:
  level: info
  format: json
```

```yaml
# configs/user-service.yaml
server:
  user_service_port: 50051

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: beehive
  sslmode: disable

jwt:
  secret: your-secret-key
  expire_hours: 24

log:
  level: info
  format: json
```

### 使用示例

#### 开发环境

```bash
# 启动 Gateway
go run cmd/gateway/main.go serve

# 启动 User Service
go run cmd/user-service/main.go serve

# 启动 Message Service
go run cmd/message-service/main.go serve

# 启动 Presence Service
go run cmd/presence-service/main.go serve

# 使用配置文件
go run cmd/gateway/main.go serve --config configs/gateway.yaml
```

#### 生产环境

```bash
# 构建
go build -o bin/gateway cmd/gateway/main.go
go build -o bin/user-service cmd/user-service/main.go
go build -o bin/message-service cmd/message-service/main.go
go build -o bin/presence-service cmd/presence-service/main.go

# 运行
./bin/gateway serve --config /etc/beehive/gateway.yaml
./bin/user-service serve --config /etc/beehive/user-service.yaml
./bin/user-service migrate --config /etc/beehive/user-service.yaml
```

#### 查看帮助

```bash
# 查看所有命令
./bin/gateway --help

# 查看 serve 命令帮助
./bin/gateway serve --help

# 查看版本
./bin/gateway version
```

### Makefile 集成

```makefile
# Makefile
.PHONY: build run migrate test clean

SERVICES = gateway user-service message-service presence-service

# 构建所有服务
build:
	@echo "Building services..."
	@for service in $(SERVICES); do \
		go build -o bin/$$service cmd/$$service/main.go; \
	done

# 运行所有服务
run:
	@echo "Starting services..."
	@docker-compose -f docker/docker-compose.yml up -d
	@./bin/gateway serve --config configs/gateway.yaml &
	@./bin/user-service serve --config configs/user-service.yaml &
	@./bin/message-service serve --config configs/message-service.yaml &
	@./bin/presence-service serve --config configs/presence-service.yaml &

# 数据库迁移
migrate:
	@./bin/user-service migrate --config configs/user-service.yaml
	@./bin/message-service migrate --config configs/message-service.yaml

# 运行测试
test:
	@go test ./...

# 清理
clean:
	@rm -rf bin/
	@docker-compose -f docker/docker-compose.yml down
```

## 服务间通信

### Gateway 到业务服务

Gateway 通过 gRPC 调用业务服务：

```go
// internal/gateway/grpc/client.go
package grpc

import (
    "google.golang.org/grpc"
    pb "github.com/HappyLadySauce/Beehive/internal/proto/user"
)

type Client struct {
    userService    pb.UserServiceClient
    messageService pb.MessageServiceClient
    presenceService pb.PresenceServiceClient
}

func NewClient(cfg *config.Config) (*Client, error) {
    // 连接 User Service
    userConn, err := grpc.Dial(cfg.GRPC.UserServiceAddr, grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    
    // 连接 Message Service
    msgConn, err := grpc.Dial(cfg.GRPC.MessageServiceAddr, grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    
    // 连接 Presence Service
    presConn, err := grpc.Dial(cfg.GRPC.PresenceServiceAddr, grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    
    return &Client{
        userService:    pb.NewUserServiceClient(userConn),
        messageService: pb.NewMessageServiceClient(msgConn),
        presenceService: pb.NewPresenceServiceClient(presConn),
    }, nil
}
```

### 消息队列通信

Message Service 作为 Producer，Gateway 作为 Consumer：

```go
// Message Service 发送消息
producer.PublishSingleChatMessage(userID, message)

// Gateway 消费消息
consumer.ConsumeSingleChatMessages()
```

## 部署建议

### 开发环境

- 所有服务运行在同一台机器
- 使用 Docker Compose 管理基础设施
- 使用 `go run` 启动服务

### 生产环境

- 每个服务独立部署
- 使用 Kubernetes 或 Docker Swarm 编排
- 使用配置中心管理配置
- 实现服务发现和负载均衡

### 监控和日志

- 每个服务暴露健康检查接口
- 统一日志格式（JSON）
- 集成 Prometheus 指标
- 使用分布式追踪（Jaeger）

## 最佳实践

1. **配置管理**：使用环境变量覆盖配置文件
2. **错误处理**：统一的错误码和错误信息
3. **日志记录**：结构化日志，包含请求ID
4. **优雅关闭**：正确处理信号，等待请求完成
5. **健康检查**：实现 `/health` 和 `/ready` 接口
6. **版本管理**：使用 Git 标签管理版本号
7. **文档**：为每个命令添加清晰的帮助信息

## 参考

- [Cobra 官方文档](https://github.com/spf13/cobra)
- [完整开发指南](./00-完整开发指南.md)
- [用户登录与操作逻辑](./01-用户登录与操作逻辑.md)
- [WebSocket Gateway 设计](./02-WebSocket-Gateway设计.md)
- [消息队列设计](./03-消息队列设计.md)