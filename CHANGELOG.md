# Beehive IM 更新日志

## [0.1.0] - 2026-01-21

### 🎉 项目初始化

#### 新增

**项目结构**
- ✅ 创建项目目录结构（api/ + app/）
- ✅ 配置 Docker Compose 基础设施
- ✅ 创建所有服务的 Proto 文件

**API 定义**
- ✅ User RPC Proto (用户服务)
- ✅ Friend RPC Proto (好友服务)
- ✅ Chat RPC Proto (会话服务)
- ✅ Message RPC Proto (消息服务)
- ✅ File RPC Proto (文件服务)
- ✅ Search RPC Proto (搜索服务)
- ✅ Gateway API 定义

**脚本工具**
- ✅ `scripts/gen_rpc_code.sh` - RPC 代码生成脚本
- ✅ `scripts/init_db.sql` - 数据库初始化脚本
- ✅ `scripts/init_es.sh` - Elasticsearch 初始化脚本
- ✅ `scripts/init_rabbitmq.sh` - RabbitMQ 初始化脚本

**构建工具**
- ✅ Makefile - 常用命令集成
- ✅ 代码生成命令
- ✅ 服务启动命令
- ✅ 基础设施管理命令

**文档**
- ✅ [架构设计文档](docs/dev/architecture.md) - 系统架构和设计思想
- ✅ [数据库设计文档](docs/dev/database.md) - 9张表结构和 ER 图
- ✅ [API 接口文档](docs/dev/api.md) - REST API 和 WebSocket
- ✅ [RPC 服务文档](docs/dev/rpc.md) - 6个 RPC 服务设计
- ✅ [消息队列文档](docs/dev/message-queue.md) - RabbitMQ 配置
- ✅ [搜索引擎文档](docs/dev/elasticsearch.md) - Elasticsearch 配置
- ✅ [部署文档](docs/dev/deployment.md) - Docker Compose 和 Kubernetes
- ✅ [快速开始指南](QUICKSTART.md) - 10分钟上手
- ✅ [项目 README](README.md) - 项目介绍

**基础设施**
- ✅ PostgreSQL 15 (主数据库)
- ✅ Redis 7 (缓存)
- ✅ RabbitMQ 3.12 (消息队列)
- ✅ Elasticsearch 8.11 (搜索引擎)
- ✅ Kibana 5.6 (ES 可视化)
- ✅ etcd 3.5 (服务发现)

**数据库设计**
- ✅ users 表（用户表）
- ✅ email_verification_codes 表（邮箱验证码表）
- ✅ friends 表（好友关系表）
- ✅ friend_requests 表（好友申请表）
- ✅ conversations 表（会话表）
- ✅ conversation_members 表（会话成员表）
- ✅ messages 表（消息表，支持分区）
- ✅ files 表（文件表，支持去重）

#### 技术特性

- 🎯 微服务架构（6个 RPC + 1个 Gateway）
- 🎯 go-zero 框架（服务治理）
- 🎯 gRPC 高性能通信
- 🎯 etcd 服务发现
- 🎯 RabbitMQ 异步消息
- 🎯 Elasticsearch 全文检索
- 🎯 Redis 缓存策略
- 🎯 PostgreSQL 分区表
- 🎯 JWT 认证
- 🎯 文件去重（SHA256）

### 📊 统计

- **Proto 文件**: 6 个
- **RPC 服务**: 6 个
- **数据库表**: 8 张
- **文档页数**: 8 篇
- **脚本工具**: 4 个
- **代码行数**: ~10,000+ 行（文档 + 配置）

### 🚀 下一步计划

- [ ] 实现所有 RPC 服务业务逻辑
- [ ] 实现 WebSocket 连接管理
- [ ] 实现邮件发送服务（网易邮箱 SMTP）
- [ ] 实现 RabbitMQ 消费者
- [ ] 实现 Elasticsearch 搜索逻辑
- [ ] 实现文件上传服务
- [ ] 编写单元测试
- [ ] 编写集成测试
- [ ] 前端开发（Web、Desktop）

### 📝 备注

本版本完成了项目的基础架构设计和代码生成框架，为后续业务逻辑实现奠定了坚实基础。

---

**提交信息**: `feat: 初始化项目结构和完整架构设计`

**作者**: HappyLadySauce  
**日期**: 2026-01-21
