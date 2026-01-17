# Beehive IM 系统开发文档

## 文档目录

本目录包含 Beehive IM 系统的完整开发文档，帮助开发者理解系统架构、设计思路和实现细节。

### 文档列表

0. **[完整开发指南](./00-完整开发指南.md)** ⭐ **新手必读**
   - 从零开始的完整开发流程
   - 分阶段开发计划（30天）
   - 每日开发检查清单
   - 代码实现示例
   - 常见问题排查

1. **[用户登录与操作逻辑](./01-用户登录与操作逻辑.md)**
   - 用户认证流程（JWT Token）
   - WebSocket 连接建立
   - 发送单聊消息流程
   - 发送群聊消息流程
   - 查询消息历史
   - 离线消息处理
   - 心跳机制
   - 错误处理

2. **[WebSocket Gateway 设计](./02-WebSocket-Gateway设计.md)**
   - Gateway 整体架构
   - Connection Manager（连接管理器）
   - WebSocket Handler（连接处理器）
   - MQ Consumer（消息队列消费者）
   - 消息格式定义
   - 性能优化策略
   - 监控和日志

3. **[消息队列设计](./03-消息队列设计.md)**
   - RabbitMQ 架构概览
   - Exchange 设计（Direct、Topic、Fanout）
   - Queue 设计（用户队列、群组队列、状态队列）
   - 消息流转流程
   - 消息确认机制
   - 死信队列
   - 性能优化
   - 监控和告警

4. **[微服务架构与 Cobra 框架](./04-微服务架构与Cobra框架.md)** ⭐ **架构必读**
   - 微服务架构设计
   - 各服务职责划分
   - Cobra 框架使用指南
   - 命令结构设计
   - 配置管理
   - 服务间通信
   - 部署建议

## 快速开始

### 1. 阅读顺序建议

**对于新手开发者，强烈建议按以下顺序阅读：**

1. **首先阅读 [完整开发指南](./00-完整开发指南.md)** - 了解整体开发流程和计划
2. **然后阅读 [微服务架构与 Cobra 框架](./04-微服务架构与Cobra框架.md)** - 理解微服务架构和框架使用
3. 接着阅读 [用户登录与操作逻辑](./01-用户登录与操作逻辑.md)，了解系统的业务流程
4. 然后阅读 [WebSocket Gateway 设计](./02-WebSocket-Gateway设计.md)，理解 Gateway 的实现
5. 最后阅读 [消息队列设计](./03-消息队列设计.md)，深入理解消息路由机制

**对于有经验的开发者：**
- 可以直接参考各个设计文档进行开发
- 遇到问题时查阅完整开发指南中的常见问题部分

### 2. 核心概念

在阅读文档前，建议先了解以下核心概念：

- **WebSocket**：全双工通信协议，用于实时消息推送
- **gRPC**：高性能 RPC 框架，用于服务间通信
- **RabbitMQ**：消息队列中间件，用于消息路由和分发
- **JWT**：JSON Web Token，用于用户认证
- **Exchange**：RabbitMQ 中的消息路由组件
- **Queue**：RabbitMQ 中的消息存储组件
- **Routing Key**：消息路由的标识符

## 系统架构概览

```
┌─────────────┐
│   Client    │ (WebSocket)
└──────┬──────┘
       │
       ▼
┌─────────────────┐
│ WebSocket Gateway│ (处理长连接)
└────────┬────────┘
         │ gRPC
         ▼
┌─────────────────┐
│  gRPC Services  │ (业务逻辑层)
│  - User Service │
│  - Message Svc  │
│  - Presence Svc │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   RabbitMQ      │ (消息队列)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ MQ Consumers    │ (消息分发)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Database      │ (PostgreSQL/MySQL)
└─────────────────┘
```

## 技术栈

- **语言**：Go 1.25+
- **WebSocket**：gorilla/websocket
- **gRPC**：google.golang.org/grpc
- **消息队列**：RabbitMQ (amqp091-go)
- **数据库**：PostgreSQL/MySQL (GORM)
- **认证**：JWT (golang-jwt/jwt)

## 开发环境搭建

### 1. 依赖安装

```bash
# 安装 Go 依赖
go mod download

# 启动 RabbitMQ（使用 Docker）
docker-compose up -d rabbitmq

# 启动数据库（使用 Docker）
docker-compose up -d postgres
```

### 2. 配置环境变量

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=beehive

# RabbitMQ 配置
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

# JWT 密钥
JWT_SECRET=your-secret-key
```

### 3. 运行服务

```bash
# 启动 WebSocket Gateway
go run cmd/gateway/main.go

# 启动 Message Service
go run cmd/message-service/main.go

# 启动 User Service
go run cmd/user-service/main.go
```

## 代码结构

```
Beehive/
├── api/
│   └── proto/              # Protocol Buffers 定义
├── cmd/
│   ├── gateway/            # WebSocket 网关入口
│   ├── user-service/        # 用户服务入口
│   ├── message-service/     # 消息服务入口
│   └── presence-service/    # 在线状态服务入口
├── internal/
│   ├── gateway/
│   │   ├── websocket/       # WebSocket 处理
│   │   └── connection/       # 连接管理
│   ├── service/
│   │   ├── user/           # 用户服务实现
│   │   ├── message/         # 消息服务实现
│   │   └── presence/        # 在线状态服务实现
│   ├── mq/
│   │   ├── producer.go     # RabbitMQ 生产者
│   │   ├── consumer.go     # RabbitMQ 消费者
│   │   └── config.go       # MQ 配置
│   └── pkg/
│       ├── model/          # 数据模型
│       ├── db/             # 数据库连接
│       └── config/         # 配置管理
└── docs/                   # 开发文档
```

## 常见问题

### Q1: 为什么使用 RabbitMQ 而不是直接推送？

**A**: RabbitMQ 提供了以下优势：
- **解耦**：业务服务和消息推送分离
- **可靠性**：消息持久化，服务重启不丢失
- **扩展性**：支持多实例部署，水平扩展
- **削峰**：高并发时消息先进入队列，按能力消费

### Q2: 如何支持多设备登录？

**A**: Connection Manager 维护用户ID到连接列表的映射，一个用户可以有多个连接。消息推送时会推送给该用户的所有连接。

### Q3: 离线消息如何处理？

**A**: 当用户离线时，Consumer 将消息标记为未读存储在数据库。用户上线后，Gateway 查询未读消息并批量推送。

### Q4: 如何实现消息的可靠性？

**A**: 
- 消息持久化：Exchange 和 Queue 设置为 durable
- 手动确认：Consumer 处理成功后手动 Ack
- 死信队列：处理失败的消息发送到死信队列
- 消息重试：失败消息自动重试

## 贡献指南

1. 阅读相关文档，理解设计思路
2. 遵循代码规范，保持代码风格一致
3. 添加必要的注释和文档
4. 编写单元测试
5. 提交 Pull Request

## 更新日志

### 2024-01-XX
- 初始版本
- 完成用户登录与操作逻辑文档
- 完成 WebSocket Gateway 设计文档
- 完成消息队列设计文档

## 联系方式

如有问题或建议，请提交 Issue 或联系开发团队。