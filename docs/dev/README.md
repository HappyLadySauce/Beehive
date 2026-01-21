# Beehive IM 开发文档

## 文档概览

本目录包含 Beehive IM 系统的完整开发文档，涵盖架构设计、数据库设计、API 接口、RPC 服务、消息队列、搜索引擎和部署方案。

## 文档列表

### 1. [架构设计](architecture.md)

**内容**:
- 系统概述和技术栈
- 微服务架构图和服务职责划分
- 数据流转和核心业务流程
- 高可用设计和安全设计
- 性能优化策略
- 监控和可观测性
- 项目目录结构

**适合人群**: 架构师、技术负责人、新加入的开发人员

**关键亮点**:
- 清晰的微服务拆分（6个RPC服务 + API Gateway）
- 完整的数据流转图（用户注册、消息发送、好友申请）
- 基于 go-zero 框架的微服务治理能力

---

### 2. [数据库设计](database.md)

**内容**:
- PostgreSQL 数据库选型理由
- 完整的表结构设计（9张表）
- 索引策略和性能优化
- 缓存策略（Redis）
- 数据备份和恢复方案
- ER 图和表关系

**适合人群**: 后端开发人员、DBA

**核心表**:
- `users`: 用户表
- `friends`: 好友关系表
- `friend_requests`: 好友申请表
- `conversations`: 会话表
- `conversation_members`: 会话成员表
- `messages`: 消息表（分区表）
- `files`: 文件表（支持去重）
- `email_verification_codes`: 邮箱验证码表

**特色设计**:
- 消息表按月分区
- 文件去重（SHA256哈希）
- 好友关系双向存储
- 完善的索引设计

---

### 3. [API 接口设计](api.md)

**内容**:
- RESTful API 规范
- 完整的接口定义（用户、好友、会话、消息、文件）
- WebSocket 接口定义
- 请求响应示例
- 错误码说明
- 接口测试指南

**适合人群**: 前端开发人员、接口测试人员

**接口分类**:
- **用户相关**: 注册、登录、获取/更新用户信息
- **好友相关**: 发送申请、处理申请、好友列表、删除好友
- **会话相关**: 创建会话、会话列表、会话详情、标记已读
- **消息相关**: 获取历史消息、搜索消息
- **文件相关**: 上传、下载、批量上传
- **WebSocket**: 实时消息推送、心跳检测

**特色功能**:
- JWT 认证
- WebSocket 长连接
- 全文检索
- 文件去重

---

### 4. [RPC 服务设计](rpc.md)

**内容**:
- gRPC 和 Protocol Buffers 介绍
- 6个 RPC 服务的 Proto 定义
- 服务实现要点
- RPC 调用示例
- 性能优化和监控

**适合人群**: 后端开发人员

**RPC 服务列表**:
1. **User RPC** (8001): 用户管理、认证、在线状态
2. **Friend RPC** (8002): 好友关系、好友申请
3. **Chat RPC** (8004): 会话管理、会话成员
4. **Message RPC** (8003): 消息发送、历史消息
5. **File RPC** (8005): 文件上传、下载、去重
6. **Search RPC** (8006): 消息全文检索

**技术特点**:
- 基于 gRPC 的高性能通信
- etcd 服务发现
- go-zero 内置负载均衡
- 自适应熔断和限流

---

### 5. [消息队列设计](message-queue.md)

**内容**:
- RabbitMQ 架构设计
- Exchange 和 Queue 配置
- 消息格式定义
- 生产者和消费者实现
- 消息可靠性保证
- 监控和告警

**适合人群**: 后端开发人员、运维人员

**队列设计**:
- **message.persist**: 消息持久化队列
- **message.push**: 消息推送队列（Gateway 消费）
- **message.index**: 消息索引队列（Search RPC 消费）

**可靠性保证**:
- 生产者确认
- 消费者手动 ACK
- 消息持久化
- 死信队列

---

### 6. [Elasticsearch 设计](elasticsearch.md)

**内容**:
- Elasticsearch 索引设计
- IK 中文分词器配置
- 索引和搜索操作
- Go 客户端实现
- 性能优化
- 备份和恢复

**适合人群**: 后端开发人员、搜索工程师

**核心功能**:
- 历史消息全文检索
- 搜索结果高亮
- 按会话过滤
- 分页查询

**技术实现**:
- IK 分词器（ik_smart + ik_max_word）
- 按月创建索引
- 别名机制
- 批量索引（Bulk API）

---

### 7. [部署文档](deployment.md)

**内容**:
- 开发环境部署（Docker Compose）
- 生产环境部署（Kubernetes）
- Docker 镜像构建
- 监控和日志
- 备份和恢复
- 运维操作

**适合人群**: 运维人员、DevOps 工程师

**部署方案**:
- **开发环境**: Docker Compose（单机）
- **生产环境**: Kubernetes（集群）

**基础设施**:
- PostgreSQL（主数据库）
- Redis（缓存）
- RabbitMQ（消息队列）
- Elasticsearch（搜索引擎）
- etcd（服务发现）

---

## 快速开始

### 1. 环境准备

```bash
# 安装 Go
brew install go  # macOS
apt install golang  # Ubuntu

# 安装 Docker
brew install docker docker-compose  # macOS
apt install docker.io docker-compose  # Ubuntu

# 安装 goctl
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 验证安装
goctl --version
```

### 2. 启动基础设施

```bash
cd /opt/Beehive/docker
docker-compose up -d

# 等待服务启动
docker-compose ps
```

### 3. 初始化数据库

```bash
# 创建数据库
docker exec -it beehive-postgres psql -U postgres -c "CREATE DATABASE beehive;"

# 执行初始化脚本
docker exec -i beehive-postgres psql -U postgres -d beehive < scripts/init_db.sql
```

### 4. 初始化 Elasticsearch

```bash
./scripts/init_es.sh
```

### 5. 初始化 RabbitMQ

```bash
./scripts/init_rabbitmq.sh
```

### 6. 生成代码

```bash
./scripts/gen_code.sh
```

### 7. 启动服务

```bash
# 终端1: User RPC
cd rpc/user && go run user.go -f etc/user.yaml

# 终端2: Friend RPC
cd rpc/friend && go run friend.go -f etc/friend.yaml

# 终端3: Chat RPC
cd rpc/chat && go run chat.go -f etc/chat.yaml

# 终端4: Message RPC
cd rpc/message && go run message.go -f etc/message.yaml

# 终端5: File RPC
cd rpc/file && go run file.go -f etc/file.yaml

# 终端6: Search RPC
cd rpc/search && go run search.go -f etc/search.yaml

# 终端7: Gateway
cd api/gateway && go run gateway.go -f etc/gateway.yaml
```

### 8. 测试

```bash
# 健康检查
curl http://localhost:8888/ping

# 发送验证码
curl -X POST http://localhost:8888/api/v1/auth/send-code \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","purpose":"register"}'
```

---

## 技术栈总结

### 后端

- **语言**: Go 1.21+
- **框架**: go-zero
- **通信协议**: HTTP REST、gRPC、WebSocket
- **数据库**: PostgreSQL 15
- **缓存**: Redis 7
- **消息队列**: RabbitMQ 3.12
- **搜索引擎**: Elasticsearch 8.11
- **服务发现**: etcd 3.5
- **认证**: JWT

### 前端（待开发）

- **Web**: React + TypeScript
- **Desktop**: Electron
- **Mobile**: React Native / Flutter

### 基础设施

- **容器化**: Docker
- **编排**: Kubernetes
- **监控**: Prometheus + Grafana
- **日志**: ELK (Elasticsearch + Logstash + Kibana)
- **链路追踪**: OpenTelemetry + Jaeger

---

## 架构特点

### 1. 微服务架构

- 服务拆分合理，职责清晰
- 通过 etcd 实现服务发现
- go-zero 内置负载均衡

### 2. 高可用设计

- RPC 服务多实例部署
- 数据库主从复制
- Redis Cluster
- RabbitMQ 集群
- Elasticsearch 集群

### 3. 高性能设计

- Redis 缓存热点数据
- 消息队列异步处理
- 数据库索引优化
- 消息表分区

### 4. 可扩展性

- 水平扩展：增加服务实例
- 垂直扩展：增加机器资源
- 数据库分库分表

### 5. 安全性

- JWT 认证
- 密码 bcrypt 加密
- SQL 注入防护
- XSS 防护

---

## 开发规范

### 1. Git 提交规范

遵循 Angular 提交规范：

- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 重构
- `test`: 测试
- `chore`: 构建/工具链

示例：

```
feat(user): 添加用户注册功能

- 实现邮箱验证码发送
- 实现用户注册接口
- 添加单元测试
```

### 2. 代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 使用 `golangci-lint` 检查代码
- 注释使用中文

### 3. 接口命名规范

- RESTful API: 使用名词复数，如 `/api/v1/users`
- RPC 方法: 使用动词开头，如 `GetUserInfo`、`SendMessage`

---

## 常见问题

### Q1: 如何添加新的 RPC 服务？

A:
1. 在 `api/proto/` 下创建新的 proto 文件
2. 使用 `goctl rpc protoc` 生成代码
3. 实现业务逻辑
4. 在 Gateway 中注入新的 RPC Client

### Q2: 如何添加新的 API 接口？

A:
1. 在 `api/gateway/v1/gateway.api` 中定义接口
2. 使用 `goctl api go` 重新生成代码
3. 实现 Logic 层业务逻辑
4. 调用 RPC 服务

### Q3: 如何调试？

A:
1. 查看日志：`logs/` 目录
2. 使用 `pprof` 性能分析
3. 使用 Postman 测试接口
4. 使用 grpcurl 测试 RPC

### Q4: 如何部署到生产环境？

A: 参考 [部署文档](deployment.md)

---

## 项目进度

### 已完成

- [x] 架构设计
- [x] 数据库设计
- [x] API 接口设计
- [x] RPC 服务设计
- [x] 消息队列设计
- [x] Elasticsearch 设计
- [x] 部署方案设计

### 待完成

- [ ] 生成所有 RPC 服务代码
- [ ] 实现业务逻辑
- [ ] 实现 WebSocket 连接管理
- [ ] 实现邮件发送服务
- [ ] 实现 RabbitMQ 消费者
- [ ] 实现 Elasticsearch 搜索
- [ ] 实现文件上传服务
- [ ] 编写单元测试
- [ ] 编写集成测试
- [ ] 前端开发（Web、Desktop）

---

## 联系方式

- **项目地址**: https://github.com/HappyLadySauce/Beehive
- **文档地址**: /opt/Beehive/docs/
- **作者**: HappyLadySauce
- **邮箱**: 13452552349@163.com

---

## 参考资源

- [go-zero 官方文档](https://go-zero.dev/)
- [go-zero GitHub](https://github.com/zeromicro/go-zero)
- [go-zero 书店示例](https://github.com/zeromicro/zero-examples/tree/main/bookstore)
- [Protocol Buffers 文档](https://protobuf.dev/)
- [gRPC 文档](https://grpc.io/docs/)
- [PostgreSQL 文档](https://www.postgresql.org/docs/)
- [Redis 文档](https://redis.io/docs/)
- [RabbitMQ 文档](https://www.rabbitmq.com/docs)
- [Elasticsearch 文档](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html)

---

## 更新日志

### 2026-01-21

- 创建完整的开发文档
- 完成架构设计
- 完成数据库设计
- 完成 API 接口设计
- 完成 RPC 服务设计
- 完成消息队列设计
- 完成 Elasticsearch 设计
- 完成部署方案设计
