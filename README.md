# Beehive IM - 企业级即时通讯系统

Beehive IM 是一个基于 Go 语言开发的企业级即时通讯系统，采用微服务架构，支持单聊、群聊、在线状态管理、消息全文搜索等功能。

## 核心特性

- ✅ **微服务架构**：独立的认证、用户、消息、在线状态、搜索服务
- ✅ **实时通信**：基于 WebSocket 的实时消息推送
- ✅ **消息队列**：使用 RabbitMQ 实现异步消息处理
- ✅ **全文搜索**：集成 Elasticsearch 实现历史消息快速检索
- ✅ **服务发现**：基于 etcd 的服务注册与发现
- ✅ **认证授权**：JWT Token 认证，支持 Token 刷新和撤销
- ✅ **数据持久化**：PostgreSQL 存储消息和用户数据
- ✅ **缓存支持**：Redis 缓存 Token 和在线状态
- ✅ **中文分词**：IK Analyzer 中文分词器

## 技术栈

### 后端技术

- **开发语言**: Go 1.21+
- **通信协议**: gRPC, WebSocket
- **消息队列**: RabbitMQ
- **搜索引擎**: Elasticsearch 8.11 + IK Analyzer
- **数据库**: PostgreSQL 15
- **缓存**: Redis 7
- **服务发现**: etcd 3.5
- **日志**: spdlog (Go 版本)
- **配置管理**: Viper
- **CLI**: Cobra

### 基础设施

- **容器化**: Docker, Docker Compose
- **数据可视化**: Kibana
- **消息队列管理**: RabbitMQ Management
- **构建工具**: Make, CMake (C++ 组件)

## 快速开始

### 1. 环境准备

**系统要求**：
- Go 1.21+
- Docker & Docker Compose
- Make

**安装依赖**：
```bash
# 安装 Go 依赖
go mod download

# 安装 Protocol Buffers 编译器
# macOS
brew install protobuf

# Linux
sudo apt-get install protobuf-compiler

# 安装 Go 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 2. 启动基础设施

```bash
# 启动所有基础设施服务
cd docker
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

基础设施服务包括：
- PostgreSQL (5432)
- Redis (6379)
- RabbitMQ (5672, 管理界面 15672)
- Elasticsearch (9200, 9300)
- Kibana (5601)
- etcd (2379, 2380)

### 3. 初始化 Elasticsearch

```bash
# 安装 IK 中文分词插件
docker exec -it beehive-elasticsearch bash
elasticsearch-plugin install https://github.com/medcl/elasticsearch-analysis-ik/releases/download/v8.11.0/elasticsearch-analysis-ik-8.11.0.zip
exit
docker restart beehive-elasticsearch

# 等待 Elasticsearch 启动完成（约30秒）
# 创建消息索引
curl -X PUT "http://localhost:9200/beehive-messages" -H 'Content-Type: application/json' -d @scripts/es-index-mapping.json

# 验证索引创建成功
curl http://localhost:9200/beehive-messages
```

### 4. 配置服务

```bash
# 复制配置文件模板
cp configs/beehive-auth-example.yaml configs/beehive-auth.yaml
cp configs/beehive-user-example.yaml configs/beehive-user.yaml
cp configs/beehive-message-example.yaml configs/beehive-message.yaml
cp configs/beehive-search-example.yaml configs/beehive-search.yaml

# 根据实际环境修改配置文件
vim configs/beehive-auth.yaml
```

### 5. 启动微服务

```bash
# 启动 Auth Service
go run cmd/beehive-auth/main.go --config configs/beehive-auth.yaml

# 启动 User Service
go run cmd/beehive-user/main.go --config configs/beehive-user.yaml

# 启动 Message Service
go run cmd/beehive-message/main.go --config configs/beehive-message.yaml

# 启动 Presence Service
go run cmd/beehive-presence/main.go --config configs/beehive-presence.yaml

# 启动 Search Service
go run cmd/beehive-search/main.go --config configs/beehive-search.yaml

# 启动 Gateway Service
go run cmd/beehive-gateway/main.go --config configs/beehive-gateway.yaml
```

或使用 Make 命令：

```bash
# 启动所有服务
make run-all

# 启动特定服务
make run-auth
make run-user
make run-message
make run-search
make run-gateway
```

### 6. 测试功能

```bash
# 测试用户注册
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "nickname": "测试用户",
    "email": "test@example.com",
    "password": "password123"
  }'

# 测试用户登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "id": "user_id",
    "password": "password123"
  }'

# 测试消息搜索
curl -X POST http://localhost:8080/api/v1/messages/search \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "keyword": "你好",
    "limit": 10,
    "offset": 0
  }'
```

## 系统架构

### 微服务架构图

```
客户端 (Web/Mobile)
    |
    v
Gateway (WebSocket/HTTP)
    |
    +-- Auth Service (认证授权)
    +-- User Service (用户管理)
    +-- Message Service (消息服务)
    +-- Presence Service (在线状态)
    +-- Search Service (消息搜索)
    |
    v
基础设施层
    +-- PostgreSQL (数据持久化)
    +-- Redis (缓存)
    +-- RabbitMQ (消息队列)
    +-- Elasticsearch (全文搜索)
    +-- etcd (服务注册)
```

### 核心微服务

#### 1. Auth Service (50050)
- 用户登录认证
- JWT Token 生成和验证
- Token 刷新和撤销
- Token 黑名单管理

#### 2. User Service (50051)
- 用户注册
- 用户信息查询和更新
- 用户资料管理

#### 3. Message Service (50052)
- 单聊/群聊消息发送
- 消息历史查询
- 未读消息管理
- 消息状态更新
- **自动同步消息到 Elasticsearch**

#### 4. Presence Service (50053)
- 用户在线状态管理
- 上线/下线通知
- 在线用户查询

#### 5. Search Service (50054)
- 消息全文搜索
- 单聊消息搜索
- 群聊消息搜索
- 搜索结果高亮
- 时间范围筛选

#### 6. Gateway Service (8080)
- WebSocket 连接管理
- HTTP API
- 消息路由和推送

## 消息搜索功能

### 搜索特性

- **全文搜索**：支持消息内容的全文检索
- **中文分词**：使用 IK Analyzer 实现精准的中文分词
- **高亮显示**：搜索结果自动高亮关键词
- **多维度筛选**：支持按用户、群组、时间范围、消息类型筛选
- **高性能**：毫秒级搜索响应，支持海量消息检索

### 搜索示例

**搜索用户所有相关消息**：
```bash
curl -X POST http://localhost:50054/search.v1.SearchService/SearchMessages \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_001",
    "keyword": "项目进展",
    "limit": 20,
    "offset": 0
  }'
```

**搜索两个用户之间的消息**：
```bash
curl -X POST http://localhost:50054/search.v1.SearchService/SearchUserMessages \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_001",
    "target_user_id": "user_002",
    "keyword": "会议",
    "limit": 20
  }'
```

**搜索群组消息**：
```bash
curl -X POST http://localhost:50054/search.v1.SearchService/SearchGroupMessages \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": "group_001",
    "keyword": "方案",
    "limit": 20
  }'
```

### Kibana 数据可视化

访问 http://localhost:5601 使用 Kibana 可视化消息数据：

1. **查看索引**: Management → Index Management
2. **搜索消息**: Discover → 选择 beehive-messages 索引
3. **分析统计**: Dashboard → 创建自定义图表

## 项目结构

```
Beehive/
├── cmd/                          # 微服务入口
│   ├── beehive-auth/            # Auth Service
│   ├── beehive-user/            # User Service
│   ├── beehive-message/         # Message Service
│   ├── beehive-presence/        # Presence Service
│   ├── beehive-search/          # Search Service
│   └── beehive-gateway/         # Gateway Service
├── internal/                     # 内部实现
│   ├── beehive-auth/            # Auth 服务实现
│   ├── beehive-user/            # User 服务实现
│   ├── beehive-message/         # Message 服务实现
│   ├── beehive-search/          # Search 服务实现
│   └── pkg/                     # 内部共享包
│       ├── elasticsearch/       # ES 客户端
│       ├── registry/            # 服务注册
│       └── middleware/          # 中间件
├── pkg/                         # 公共包
│   ├── api/                     # API 定义
│   │   └── proto/              # Proto 文件
│   │       ├── auth/v1/        # Auth Service Proto
│   │       ├── user/v1/        # User Service Proto
│   │       ├── message/v1/     # Message Service Proto
│   │       ├── presence/v1/    # Presence Service Proto
│   │       └── search/v1/      # Search Service Proto
│   └── utils/                  # 工具函数
├── configs/                     # 配置文件
│   ├── beehive-auth-example.yaml
│   ├── beehive-user-example.yaml
│   ├── beehive-message-example.yaml
│   └── beehive-search-example.yaml
├── docker/                      # Docker 配置
│   ├── docker-compose.yml      # 完整基础设施
│   └── README.md               # Docker 使用文档
├── docs/                        # 文档
│   └── dev/                    # 开发文档
│       ├── 00-微服务架构设计.md
│       ├── 01-用户登录与操作逻辑.md
│       ├── 02-Auth认证架构设计.md
│       ├── 03-消息队列设计.md
│       ├── 04-完整开发指南.md
│       └── 05-Elasticsearch搜索架构设计.md
├── scripts/                     # 脚本文件
│   └── sync-to-es/             # ES 数据同步工具
├── Makefile                     # Make 构建文件
└── README.md                    # 项目说明
```

## 开发指南

### 生成 Proto 代码

```bash
make proto
```

### 运行测试

```bash
# 运行所有测试
make test

# 运行特定服务测试
make test-auth
make test-user
make test-message
make test-search
```

### 代码检查

```bash
# 代码格式化
make fmt

# 代码检查
make lint
```

### 构建

```bash
# 构建所有服务
make build

# 构建特定服务
make build-auth
make build-user
make build-message
make build-search
```

## 性能优化

### Elasticsearch 优化

1. **批量索引**：Message Service 使用批量索引提高写入性能
2. **异步索引**：消息同步到 ES 采用异步方式，不阻塞消息发送
3. **分片配置**：根据数据量合理配置分片数（默认3个主分片）
4. **ILM 策略**：配置索引生命周期管理，自动归档历史数据

### 搜索性能

- 单关键词搜索：< 50ms
- 复杂查询（多条件）：< 100ms
- 支持百万级消息量检索

## 监控和运维

### 查看服务状态

```bash
# 查看所有容器状态
docker-compose ps

# 查看特定服务日志
docker-compose logs -f elasticsearch
docker-compose logs -f postgres
```

### Elasticsearch 集群健康

```bash
# 查看集群健康状态
curl http://localhost:9200/_cluster/health?pretty

# 查看索引信息
curl http://localhost:9200/_cat/indices?v

# 查看索引统计
curl http://localhost:9200/beehive-messages/_stats?pretty
```

### 数据备份

```bash
# 备份 PostgreSQL
docker exec beehive-postgres pg_dump -U postgres beehive > backup.sql

# 备份 Elasticsearch（使用快照）
curl -X PUT "http://localhost:9200/_snapshot/my_backup" -H 'Content-Type: application/json' -d'{
  "type": "fs",
  "settings": {
    "location": "/usr/share/elasticsearch/data/backup"
  }
}'
```

## 常见问题

### Q1: Elasticsearch 启动失败？

**A**: 检查 vm.max_map_count 配置：
```bash
sudo sysctl -w vm.max_map_count=262144
```

### Q2: 搜索不到中文内容？

**A**: 确认已安装 IK 分词插件：
```bash
docker exec -it beehive-elasticsearch elasticsearch-plugin list
```

### Q3: 消息没有同步到 Elasticsearch？

**A**: 
1. 检查 Message Service 配置中 `elasticsearch.sync-enabled` 是否为 true
2. 查看 Message Service 日志是否有错误
3. 检查 Elasticsearch 服务是否正常运行

### Q4: 如何重建 Elasticsearch 索引？

**A**: 使用数据同步工具：
```bash
go run cmd/tools/sync-to-es/main.go --batch-size 1000
```

## 文档

详细文档请参考 `docs/dev/` 目录：

- [微服务架构设计](docs/dev/00-微服务架构设计.md)
- [用户登录与操作逻辑](docs/dev/01-用户登录与操作逻辑.md)
- [Auth认证架构设计](docs/dev/02-Auth认证架构设计.md)
- [消息队列设计](docs/dev/03-消息队列设计.md)
- [完整开发指南](docs/dev/04-完整开发指南.md)
- [Elasticsearch搜索架构设计](docs/dev/05-Elasticsearch搜索架构设计.md)

## 路线图

- [x] 微服务架构设计
- [x] 用户认证和授权
- [x] 单聊和群聊功能
- [x] 在线状态管理
- [x] 消息队列集成
- [x] Elasticsearch 全文搜索
- [ ] 文件上传和存储
- [ ] 消息撤回功能
- [ ] 消息加密
- [ ] 管理后台
- [ ] 分布式部署方案

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 联系方式

项目地址: https://gitee.com/wang-guangke/chat_code.git
