# Beehive IM - 企业级即时通讯系统

基于 go-zero 微服务框架的企业级 IM 通讯系统，采用现代化的微服务架构，支持单聊、群聊、文件传输、历史消息全文检索等功能。

## 🚀 技术栈

### 后端

- **框架**: go-zero (微服务框架)
- **通信**: HTTP REST、gRPC、WebSocket
- **数据库**: PostgreSQL 15
- **缓存**: Redis 7
- **消息队列**: RabbitMQ 3.12
- **搜索引擎**: Elasticsearch 8.11
- **服务发现**: etcd 3.5
- **认证**: JWT
- **工具**: goctl (代码生成)

### 前端（规划中）

- **Web**: React + TypeScript
- **Desktop**: Electron
- **Mobile**: React Native / Flutter

## ✨ 核心功能

- ✅ 用户注册/登录（用户名、邮箱）
- ✅ 邮箱验证码验证
- ✅ 好友申请/处理/删除
- ✅ 单聊/群聊会话管理
- ✅ 文字/图片/语音消息
- ✅ WebSocket 实时消息推送
- ✅ 历史消息全文检索
- ✅ 文件上传去重（SHA256）
- ✅ 断点续传支持
- ✅ 用户在线状态管理

## 📖 项目结构

```
Beehive/
├── api/                       # API 定义文件（集中管理）
│   ├── beehive-gateway/       # Gateway API 定义
│   └── proto/                 # RPC Proto 定义
│       ├── beehive-user/
│       ├── beehive-friend/
│       ├── beehive-chat/
│       ├── beehive-message/
│       ├── beehive-file/
│       └── beehive-search/
├── app/                       # 应用实现代码
│   ├── beehive-gateway/       # API Gateway 实现
│   ├── beehive-user/          # User RPC 实现
│   ├── beehive-friend/        # Friend RPC 实现
│   ├── beehive-chat/          # Chat RPC 实现
│   ├── beehive-message/       # Message RPC 实现
│   ├── beehive-file/          # File RPC 实现
│   └── beehive-search/        # Search RPC 实现
├── common/                    # 公共代码
├── docker/                    # Docker 配置
├── docs/                      # 项目文档
├── scripts/                   # 脚本工具
└── Makefile                   # 常用命令
```

## 🎯 快速开始

### 1. 环境要求

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+
- RabbitMQ 3.12+
- Elasticsearch 8.11+

### 2. 安装 goctl

```bash
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 验证安装
goctl --version
```

### 3. 启动基础设施

```bash
# 启动所有基础设施（PostgreSQL, Redis, RabbitMQ, Elasticsearch, etcd）
make docker-up

# 等待服务就绪后，初始化数据库
make init-db

# 初始化 Elasticsearch
make init-es

# 初始化 RabbitMQ
make init-mq
```

### 4. 生成代码

```bash
# 生成所有 RPC 服务代码
make gen-rpc

# 或使用脚本
./scripts/gen_rpc_code.sh
```

### 5. 启动服务

**方式一：使用 Makefile（推荐）**

在不同的终端中运行：

```bash
# 终端 1: 启动 User RPC
make run-user

# 终端 2: 启动 Friend RPC
make run-friend

# 终端 3: 启动 Chat RPC
make run-chat

# 终端 4: 启动 Message RPC
make run-message

# 终端 5: 启动 File RPC
make run-file

# 终端 6: 启动 Search RPC
make run-search

# 终端 7: 启动 API Gateway
make run-gateway
```

**方式二：手动启动**

```bash
# User RPC
cd app/beehive-user && go run user.go -f etc/user.yaml

# API Gateway
cd app/beehive-gateway && go run gateway.go -f etc/gateway-api.yaml
```

### 6. 测试

```bash
# 健康检查
curl http://localhost:8888/ping

# 发送验证码
curl -X POST http://localhost:8888/api/v1/auth/send-code \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","purpose":"register"}'

# 用户注册
curl -X POST http://localhost:8888/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123","code":"123456"}'
```

## 📚 文档

完整文档位于 `docs/dev/` 目录：

- [架构设计](docs/dev/architecture.md) - 系统架构和设计思想
- [数据库设计](docs/dev/database.md) - 数据库表结构和设计
- [API 接口](docs/dev/api.md) - REST API 和 WebSocket 接口
- [RPC 服务](docs/dev/rpc.md) - gRPC 服务设计
- [消息队列](docs/dev/message-queue.md) - RabbitMQ 配置和使用
- [搜索引擎](docs/dev/elasticsearch.md) - Elasticsearch 配置和使用
- [部署文档](docs/dev/deployment.md) - 部署指南

## 🏗️ 微服务架构

### 服务列表

| 服务 | 端口 | 说明 |
|------|------|------|
| API Gateway | 8888 | HTTP/WebSocket 统一入口 |
| User RPC | 8001 | 用户服务 |
| Friend RPC | 8002 | 好友服务 |
| Chat RPC | 8004 | 会话服务 |
| Message RPC | 8003 | 消息服务 |
| File RPC | 8005 | 文件服务 |
| Search RPC | 8006 | 搜索服务 |

### 基础设施

| 服务 | 端口 | 说明 |
|------|------|------|
| PostgreSQL | 5432 | 主数据库 |
| Redis | 6379 | 缓存 |
| RabbitMQ | 5672, 15672 | 消息队列 |
| Elasticsearch | 9200 | 搜索引擎 |
| Kibana | 5601 | ES 可视化 |
| etcd | 2379 | 服务发现 |

## 🛠️ 常用命令

```bash
# 代码生成
make gen-api          # 生成 API Gateway 代码
make gen-rpc          # 生成所有 RPC 服务代码
make gen-all          # 生成所有代码

# 基础设施
make docker-up        # 启动基础设施
make docker-down      # 停止基础设施
make init-db          # 初始化数据库
make init-es          # 初始化 Elasticsearch
make init-mq          # 初始化 RabbitMQ

# 启动服务
make run-gateway      # 启动 API Gateway
make run-user         # 启动 User RPC
make run-friend       # 启动 Friend RPC
make run-chat         # 启动 Chat RPC
make run-message      # 启动 Message RPC
make run-file         # 启动 File RPC
make run-search       # 启动 Search RPC

# 帮助
make help             # 查看所有可用命令
```

## 🌟 技术亮点

### 1. 微服务架构

- 服务拆分合理，职责清晰
- 基于 go-zero 框架，自带服务治理能力
- etcd 服务发现，动态负载均衡
- 自适应熔断、限流、降级

### 2. 消息队列解耦

- RabbitMQ 实现服务异步通信
- 消息持久化、推送、索引三个队列
- 死信队列兜底，保证消息可靠性

### 3. 全文检索

- Elasticsearch 实现历史消息搜索
- IK 中文分词器
- 搜索结果高亮
- 按月创建索引，易于归档

### 4. 文件去重

- SHA256 哈希去重
- 引用计数管理
- 节省存储空间
- 支持断点续传

### 5. 高可用设计

- RPC 服务多实例部署
- 数据库主从复制
- Redis 集群
- RabbitMQ 集群
- Elasticsearch 集群

## 🔒 安全性

- JWT Token 认证
- bcrypt 密码加密
- SQL 注入防护
- XSS 防护
- 接口限流
- IP 黑名单

## 📈 性能优化

- Redis 缓存热点数据
- 消息队列异步处理
- 数据库索引优化
- 消息表分区
- WebSocket 长连接
- gRPC 高性能通信

## 🚧 开发计划

- [ ] 实现所有 RPC 服务业务逻辑
- [ ] 实现 WebSocket 连接管理
- [ ] 实现邮件发送服务
- [ ] 实现 RabbitMQ 消费者
- [ ] 实现 Elasticsearch 搜索
- [ ] 实现文件上传服务
- [ ] 编写单元测试
- [ ] 编写集成测试
- [ ] 前端开发（Web、Desktop）
- [ ] 音视频通话（WebRTC）
- [ ] 消息撤回功能
- [ ] 群公告功能
- [ ] @提醒功能

## 📝 开发规范

### Git 提交规范（Angular）

- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 重构
- `test`: 测试
- `chore`: 构建/工具链

### 代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 使用 `golangci-lint` 检查代码
- 注释使用中文

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可

MIT License

## 👨‍💻 作者

- **HappyLadySauce**
- Email: 13452552349@163.com
- GitHub: https://github.com/HappyLadySauce/Beehive

## 🙏 致谢

- [go-zero](https://github.com/zeromicro/go-zero) - 优秀的微服务框架
- [go-zero 书店示例](https://github.com/zeromicro/zero-examples/tree/main/bookstore) - 参考示例

---

⭐ 如果这个项目对你有帮助，欢迎 Star！
