# Beehive IM 开发文档

本目录包含 Beehive IM 系统的详细开发文档。

## 文档目录

### 架构设计

#### [00-微服务架构设计.md](./00-微服务架构设计.md)
- 系统整体架构设计
- 微服务划分和职责
- 服务间通信方式
- 服务注册与发现
- 配置管理

**关键内容**：
- Gateway 服务（WebSocket、HTTP API）
- Auth Service（认证授权）
- User Service（用户管理）
- Message Service（消息服务）
- Presence Service（在线状态）
- Search Service（消息搜索）

#### [01-用户登录与操作逻辑.md](./01-用户登录与操作逻辑.md)
- 用户注册流程
- 用户登录流程
- Token 认证机制
- 消息发送和接收流程

#### [02-Auth认证架构设计.md](./02-Auth认证架构设计.md)
- JWT Token 认证机制
- Token 生成和验证
- Token 刷新策略
- Token 撤销机制
- Redis 缓存设计

#### [03-消息队列设计.md](./03-消息队列设计.md)
- RabbitMQ 架构设计
- Exchange 和 Queue 配置
- 消息路由策略
- 单聊消息路由
- 群聊消息路由
- 消息可靠性保证

#### [04-完整开发指南.md](./04-完整开发指南.md)
- 从零开始构建系统
- 开发阶段划分
- 每日开发检查清单
- 常见问题排查
- Docker 环境配置

**开发阶段**：
1. 基础搭建（第1-4天）
2. Protocol Buffers 定义（第5-6天）
3. gRPC 服务实现（第7-11天）
4. WebSocket Gateway（第12-16天）
5. RabbitMQ 集成（第17-21天）
6. Elasticsearch 集成（第22-27天）
7. 测试和优化（第28-36天）

#### [05-Elasticsearch搜索架构设计.md](./05-Elasticsearch搜索架构设计.md)
- Elasticsearch 集成方案
- 索引设计和映射
- 中文分词配置
- 搜索查询实现
- 数据同步策略
- 性能优化方案
- Search Service 实现

**核心功能**：
- 消息全文搜索
- 单聊消息搜索
- 群聊消息搜索
- 搜索结果高亮
- 时间范围筛选
- 批量数据同步

## 快速导航

### 新手入门
1. 阅读 [微服务架构设计](./00-微服务架构设计.md) 了解系统整体架构
2. 阅读 [完整开发指南](./04-完整开发指南.md) 开始开发
3. 参考 [用户登录与操作逻辑](./01-用户登录与操作逻辑.md) 理解业务流程

### 实现特定功能
- **认证功能**: [Auth认证架构设计](./02-Auth认证架构设计.md)
- **消息功能**: [消息队列设计](./03-消息队列设计.md)
- **搜索功能**: [Elasticsearch搜索架构设计](./05-Elasticsearch搜索架构设计.md)

### 部署和运维
- Docker 环境：参考 `/docker/README.md`
- 服务配置：参考 `/configs/` 目录下的示例文件

## 架构图

### 系统整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                        客户端层                              │
│        Web 客户端    |    移动端    |    API 客户端         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Gateway 服务层                          │
│                WebSocket Gateway (8080)                      │
│                    HTTP API (8080)                           │
└─────────────────────────────────────────────────────────────┘
                              │
            ┌─────────────────┼─────────────────┐
            ▼                 ▼                 ▼
┌───────────────────┐ ┌───────────────┐ ┌──────────────────┐
│   Auth Service    │ │ User Service  │ │ Message Service  │
│     (50050)       │ │   (50051)     │ │    (50052)       │
└───────────────────┘ └───────────────┘ └──────────────────┘
            │                 │                 │
            ▼                 ▼                 ▼
┌───────────────────┐ ┌───────────────┐ ┌──────────────────┐
│ Presence Service  │ │Search Service │ │                  │
│     (50053)       │ │   (50054)     │ │                  │
└───────────────────┘ └───────────────┘ └──────────────────┘
                              │
            ┌─────────────────┼─────────────────┐
            ▼                 ▼                 ▼
┌───────────────────┐ ┌───────────────┐ ┌──────────────────┐
│   PostgreSQL      │ │     Redis     │ │   RabbitMQ       │
│     (5432)        │ │    (6379)     │ │   (5672)         │
└───────────────────┘ └───────────────┘ └──────────────────┘
            │                                   │
            ▼                                   ▼
┌───────────────────┐                 ┌──────────────────┐
│  Elasticsearch    │                 │      etcd        │
│  (9200, 9300)     │                 │     (2379)       │
└───────────────────┘                 └──────────────────┘
            │
            ▼
┌───────────────────┐
│      Kibana       │
│      (5601)       │
└───────────────────┘
```

## 技术栈说明

### 后端技术
- **Go 1.21+**: 主要开发语言
- **gRPC**: 微服务间通信
- **WebSocket**: 客户端实时通信
- **Protocol Buffers**: 接口定义和序列化

### 数据存储
- **PostgreSQL 15**: 主数据库，存储用户、消息、会话等数据
- **Redis 7**: 缓存 Token、在线状态等
- **Elasticsearch 8.11**: 消息全文搜索引擎

### 消息队列
- **RabbitMQ 3**: 异步消息处理和推送

### 服务发现
- **etcd 3.5**: 服务注册与发现、配置管理

### 可视化工具
- **Kibana 8.11**: Elasticsearch 数据可视化和管理

## 开发规范

### 代码规范
- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 使用 `golint` 进行代码检查

### Git 提交规范
使用 Angular 提交规范：
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建工具或辅助工具的变动

示例：
```
feat(search): 添加 Elasticsearch 全文搜索功能

- 实现 Search Service
- 集成 IK 中文分词
- 添加搜索结果高亮
```

### API 设计规范
- RESTful API 设计
- 使用 HTTP 状态码表示请求结果
- 统一的错误响应格式
- API 版本控制（v1, v2）

### 文档规范
- 每个微服务提供 README
- API 文档使用 Swagger/OpenAPI
- gRPC 服务使用 Proto 注释

## 常见问题

### Q1: 如何添加新的微服务？

1. 在 `cmd/` 目录创建服务入口
2. 在 `internal/` 目录实现服务逻辑
3. 在 `pkg/api/proto/` 定义 Proto 文件
4. 更新 `docker-compose.yml` 添加必要的基础设施
5. 更新文档

### Q2: 如何修改数据库表结构？

1. 使用数据库迁移工具（如 golang-migrate）
2. 创建迁移文件
3. 更新数据模型代码
4. 运行迁移

### Q3: 如何测试 gRPC 服务？

使用 `grpcurl` 工具：
```bash
grpcurl -plaintext localhost:50050 auth.v1.AuthService/Login
```

### Q4: 如何查看 Elasticsearch 索引？

访问 Kibana：http://localhost:5601

或使用 curl：
```bash
curl http://localhost:9200/_cat/indices?v
```

## 参考资源

### 官方文档
- [Go 语言](https://golang.org/doc/)
- [gRPC](https://grpc.io/docs/)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [Elasticsearch](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html)
- [RabbitMQ](https://www.rabbitmq.com/documentation.html)
- [etcd](https://etcd.io/docs/)

### 第三方库
- [Cobra](https://github.com/spf13/cobra) - CLI 框架
- [Viper](https://github.com/spf13/viper) - 配置管理
- [GORM](https://gorm.io/) - ORM 框架
- [go-elasticsearch](https://github.com/elastic/go-elasticsearch) - ES 客户端

### 工具
- [grpcurl](https://github.com/fullstorydev/grpcurl) - gRPC 命令行工具
- [Evans](https://github.com/ktr0731/evans) - gRPC 客户端
- [Postman](https://www.postman.com/) - API 测试工具

## 贡献指南

欢迎贡献代码和文档！

### 贡献流程
1. Fork 项目
2. 创建特性分支
3. 提交代码
4. 创建 Pull Request

### 文档贡献
- 修复文档错误
- 补充缺失的文档
- 优化文档结构
- 添加示例代码

## 联系方式

- 项目地址: https://gitee.com/wang-guangke/chat_code.git
- 问题反馈: 提交 Issue
