# Beehive IM 开发文档

## 文档索引

本文档目录包含 Beehive IM 系统的完整开发文档。

### 核心文档

1. **[微服务架构设计](./00-微服务架构设计.md)**
   - 系统整体架构
   - 微服务划分和职责
   - 服务间通信
   - 配置管理

2. **[用户登录与操作逻辑](./01-用户登录与操作逻辑.md)**
   - 用户注册流程
   - 用户登录流程
   - Token 验证流程
   - 消息发送和接收
   - 会话管理

3. **[Auth 认证架构设计](./02-Auth认证架构设计.md)**
   - Auth Service 架构
   - 认证流程详解
   - Token 管理策略
   - 安全策略
   - 性能优化

4. **[消息队列设计](./03-消息队列设计.md)**
   - RabbitMQ 架构
   - Exchange 和 Queue 配置
   - 消息路由策略
   - Producer 和 Consumer 实现
   - 消息可靠性保证

5. **[完整开发指南](./04-完整开发指南.md)**
   - 从零开始的开发步骤
   - 项目结构
   - 开发阶段概览
   - 检查清单
   - 常见问题排查

## 快速开始

### 1. 阅读顺序建议

**新手开发者**：
1. 先阅读 [微服务架构设计](./00-微服务架构设计.md) 了解整体架构
2. 阅读 [用户登录与操作逻辑](./01-用户登录与操作逻辑.md) 了解业务流程
3. 阅读 [完整开发指南](./04-完整开发指南.md) 开始开发

**架构师/技术负责人**：
1. [微服务架构设计](./00-微服务架构设计.md)
2. [Auth 认证架构设计](./02-Auth认证架构设计.md)
3. [消息队列设计](./03-消息队列设计.md)

**后端开发者**：
1. [微服务架构设计](./00-微服务架构设计.md)
2. [用户登录与操作逻辑](./01-用户登录与操作逻辑.md)
3. [Auth 认证架构设计](./02-Auth认证架构设计.md)
4. [消息队列设计](./03-消息队列设计.md)

### 2. 架构概览

```
客户端层
    ↓
Gateway 服务（WebSocket/HTTP）
    ↓
业务服务层
    ├── Auth Service（认证授权）
    ├── User Service（用户数据）
    ├── Message Service（消息服务）
    └── Presence Service（在线状态）
    ↓
基础设施层
    ├── PostgreSQL（数据库）
    ├── Redis（缓存）
    └── RabbitMQ（消息队列）
```

### 3. 关键设计决策

1. **独立的 Auth Service**：认证和用户管理完全分离
2. **JWT + Refresh Token**：支持 Token 刷新和撤销
3. **Redis 缓存**：Token 验证结果缓存，提升性能
4. **RabbitMQ 异步处理**：消息发送异步处理，提高响应速度
5. **WebSocket 长连接**：实时消息推送

## 技术栈

- **语言**：Go 1.21+
- **框架**：gRPC, Cobra
- **数据库**：PostgreSQL
- **缓存**：Redis
- **消息队列**：RabbitMQ
- **WebSocket**：gorilla/websocket
- **认证**：JWT

## 开发环境

### 前置要求

- Go 1.21 或更高版本
- Docker 和 Docker Compose
- Protocol Buffers 编译器
- Make（可选）

### 快速启动

```bash
# 1. 启动基础设施
cd docker
docker-compose up -d

# 2. 生成 Proto 代码
./scripts/generate-proto.sh

# 3. 启动服务
go run cmd/gateway/main.go serve
go run cmd/auth-service/main.go serve
go run cmd/user-service/main.go serve
go run cmd/message-service/main.go serve
go run cmd/presence-service/main.go serve
```

## 贡献指南

1. 阅读相关文档
2. 遵循代码规范
3. 编写单元测试
4. 更新文档

## 参考

- [Protocol Buffers 官方文档](https://protobuf.dev/)
- [gRPC 官方文档](https://grpc.io/)
- [RabbitMQ 官方文档](https://www.rabbitmq.com/)
- [JWT 官方文档](https://jwt.io/)
