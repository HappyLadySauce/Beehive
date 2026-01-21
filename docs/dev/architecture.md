# Beehive IM 系统架构设计

## 一、系统概述

Beehive 是一个基于 go-zero 微服务框架的企业级即时通讯系统，采用微服务架构，支持单聊、群聊、文件传输、历史消息全文检索等功能。

### 1.1 技术栈

- **后端框架**: go-zero (微服务框架)
- **通信协议**: HTTP REST、gRPC、WebSocket
- **数据库**: PostgreSQL (主数据库)
- **缓存**: Redis (会话缓存、用户在线状态)
- **消息队列**: RabbitMQ (异步消息处理)
- **搜索引擎**: Elasticsearch (全文检索)
- **服务注册与发现**: etcd
- **认证**: JWT
- **数据库操作**: sqlx + goctl model
- **代码生成**: goctl

### 1.2 核心特性

- ✅ 用户注册/登录（支持用户名、邮箱）
- ✅ 邮箱验证码验证（网易邮箱 SMTP）
- ✅ 好友申请/处理/删除
- ✅ 单聊/群聊会话管理
- ✅ 文字/图片/语音消息
- ✅ WebSocket 实时消息推送
- ✅ 历史消息全文检索（Elasticsearch）
- ✅ 文件上传去重（内容哈希）
- ✅ 断点续传支持
- ✅ 用户在线状态管理

## 二、系统架构图

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│                         客户端层                                  │
│            Web Browser / Desktop (Electron) / Mobile             │
└─────────────────────────────────────────────────────────────────┘
                                 │
                                 │ HTTP/WebSocket
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API Gateway (8888)                          │
│  ├─ HTTP REST API Handler                                       │
│  ├─ WebSocket Manager                                           │
│  ├─ JWT Authentication                                          │
│  ├─ Rate Limiting & Circuit Breaker                             │
│  └─ Request Router                                              │
└─────────────────────────────────────────────────────────────────┘
                                 │
                    ┌────────────┴────────────┐
                    │                         │
                    ▼                         ▼
         ┌──────────────────┐      ┌──────────────────┐
         │   etcd Cluster   │      │   Redis Cluster  │
         │ (Service Discovery)│    │  (Session Cache) │
         └──────────────────┘      └──────────────────┘
                    │
      ┌─────────────┼─────────────────────────┐
      │             │                         │
      ▼             ▼                         ▼
┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐
│ User RPC │  │Friend RPC│  │ Chat RPC │  │Message   │
│  (8001)  │  │  (8002)  │  │  (8004)  │  │   RPC    │
└──────────┘  └──────────┘  └──────────┘  │  (8003)  │
      │             │             │        └──────────┘
      │             │             │             │
      └─────────────┴─────────────┴─────────────┘
                          │
                          ▼
              ┌───────────────────────┐
              │   PostgreSQL Cluster  │
              │   (Main Database)     │
              └───────────────────────┘

┌──────────┐                          ┌──────────┐
│ File RPC │  ┌──────────────────┐   │Search RPC│
│  (8005)  │─▶│  File Storage    │   │  (8006)  │
└──────────┘  │  (Local/OSS)     │   └──────────┘
              └──────────────────┘          │
                                            ▼
      ┌─────────────────────────────────────────┐
      │          RabbitMQ Exchange              │
      │  ┌──────────┬──────────┬──────────┐    │
      │  │ persist  │   push   │  index   │    │
      │  │  queue   │  queue   │  queue   │    │
      └──┴──────────┴──────────┴──────────┴────┘
                                      │
                                      ▼
                            ┌──────────────────┐
                            │  Elasticsearch   │
                            │   (Full-text     │
                            │    Search)       │
                            └──────────────────┘
```

### 2.2 服务端口分配

| 服务 | 端口 | 说明 |
|------|------|------|
| API Gateway | 8888 | HTTP/WebSocket 统一入口 |
| User RPC | 8001 | 用户服务 |
| Friend RPC | 8002 | 好友服务 |
| Message RPC | 8003 | 消息服务 |
| Chat RPC | 8004 | 会话服务 |
| File RPC | 8005 | 文件服务 |
| Search RPC | 8006 | 搜索服务 |
| PostgreSQL | 5432 | 主数据库 |
| Redis | 6379 | 缓存 |
| etcd | 2379 | 服务发现 |
| RabbitMQ | 5672 | 消息队列 |
| RabbitMQ Management | 15672 | 管理界面 |
| Elasticsearch | 9200 | 搜索引擎 |
| Kibana | 5601 | ES 可视化 |

## 三、微服务职责划分

### 3.1 API Gateway Service

**职责**:
- 统一的 HTTP REST API 入口
- WebSocket 长连接管理
- JWT token 认证和鉴权
- 请求路由到各个 RPC 服务
- API 限流、熔断、降级
- CORS 跨域处理
- 请求日志记录
- 在线用户消息推送

**技术实现**:
- 基于 go-zero rest 框架
- WebSocket 连接池管理
- JWT 中间件认证
- 自适应熔断器
- 自适应限流器

### 3.2 User RPC Service

**职责**:
- 用户注册（用户名/邮箱）
- 用户登录，生成 JWT token
- 邮箱验证码发送（网易邮箱 SMTP）
- 邮箱验证码校验
- 用户信息查询（单个/批量）
- 用户信息更新（昵称、头像）
- 用户在线状态管理
- 用户密码修改

**数据表**:
- users (用户表)
- email_verification_codes (邮箱验证码表)

**缓存策略**:
- 用户信息缓存: `user:info:{user_id}`, TTL=1小时
- 用户在线状态: `user:online:{user_id}`, TTL=5分钟

### 3.3 Friend RPC Service

**职责**:
- 发送好友申请
- 处理好友申请（同意/拒绝）
- 获取好友申请列表
- 获取好友列表
- 删除好友
- 验证好友关系

**数据表**:
- friends (好友关系表)
- friend_requests (好友申请表)

**业务规则**:
- 好友关系是双向的（A 添加 B，B 自动成为 A 的好友）
- 不能重复发送好友申请（24小时内）
- 好友申请 7 天后自动过期

### 3.4 Chat RPC Service

**职责**:
- 创建会话（单聊/群聊）
- 获取用户会话列表
- 获取会话详情
- 会话成员管理（添加/移除）
- 更新未读消息数
- 标记消息已读
- 会话信息更新（群名、群头像）

**数据表**:
- conversations (会话表)
- conversation_members (会话成员表)

**缓存策略**:
- 会话列表缓存: `conv:list:{user_id}`, TTL=30分钟
- 会话详情缓存: `conv:detail:{conv_id}`, TTL=1小时

**业务规则**:
- 单聊会话：两个用户之间只能有一个单聊会话
- 群聊会话：可以创建多个
- 退出群聊后不再接收消息

### 3.5 Message RPC Service

**职责**:
- 发送消息（文字/图片/语音）
- 获取历史消息（分页）
- 更新消息状态（已读/未读）
- 消息持久化到数据库
- 消息投递到 RabbitMQ

**数据表**:
- messages (消息表)

**消息类型**:
- 1: 文本消息
- 2: 图片消息
- 3: 语音消息
- 4: 文件消息

**消息流程**:
1. 客户端通过 WebSocket 发送消息
2. Gateway 调用 Message RPC 保存消息
3. Message RPC 保存到数据库
4. Message RPC 发送消息到 RabbitMQ (3个队列)
5. Gateway 从 MQ 接收推送消息，通过 WebSocket 推送给在线用户
6. Search RPC 从 MQ 接收消息，索引到 Elasticsearch

### 3.6 File RPC Service

**职责**:
- 文件上传（头像/图片/语音/文件）
- 文件下载
- 文件信息查询
- 内容哈希去重（SHA256）
- 断点续传支持
- 文件引用计数管理
- 文件删除（引用计数为0时）

**数据表**:
- files (文件表)

**去重策略**:
- 上传前计算文件 SHA256 哈希
- 检查数据库中是否已存在该哈希
- 存在则直接返回已有文件 URL，引用计数+1
- 不存在则保存新文件

**存储方案**:
- 开发环境：本地存储 `/data/files/`
- 生产环境：阿里云 OSS / AWS S3

### 3.7 Search RPC Service

**职责**:
- 监听 RabbitMQ 消息索引队列
- 将消息索引到 Elasticsearch
- 提供全文搜索接口
- 搜索结果高亮
- 支持按会话过滤
- 支持分页

**索引内容**:
- 仅索引文本消息内容
- 图片/语音消息只索引元数据

## 四、数据流转

### 4.1 用户注册流程

```
用户 -> Gateway: POST /api/v1/auth/send-code (发送验证码)
Gateway -> User RPC: SendVerificationCode
User RPC -> PostgreSQL: 保存验证码 (有效期5分钟)
User RPC -> Email SMTP: 发送邮件
User RPC -> Gateway: Success
Gateway -> 用户: 验证码已发送

用户 -> Gateway: POST /api/v1/auth/register (注册)
Gateway -> User RPC: Register
User RPC -> PostgreSQL: 验证验证码
User RPC -> PostgreSQL: 创建用户 (密码 bcrypt 加密)
User RPC -> Redis: 缓存用户信息
User RPC -> JWT: 生成 token
User RPC -> Gateway: UserID + Token
Gateway -> 用户: 注册成功
```

### 4.2 用户登录流程

```
用户 -> Gateway: POST /api/v1/auth/login
Gateway -> User RPC: Login
User RPC -> Redis: 查询缓存
User RPC -> PostgreSQL: 验证用户名密码
User RPC -> Redis: 缓存用户信息
User RPC -> JWT: 生成 token
User RPC -> Gateway: UserInfo + Token
Gateway -> 用户: 登录成功
```

### 4.3 消息发送流程

```
用户A -> Gateway: WS: send_message
Gateway -> JWT: 验证 token
Gateway -> Friend RPC: 验证好友关系 (单聊)
Gateway -> Chat RPC: 验证会话权限
Gateway -> Message RPC: SendMessage
Message RPC -> PostgreSQL: 保存消息
Message RPC -> RabbitMQ: Publish (3个队列)
  ├─ message.persist (持久化队列)
  ├─ message.push (推送队列)
  └─ message.index (索引队列)
Message RPC -> Chat RPC: UpdateUnreadCount (更新未读数)
Message RPC -> Gateway: MessageID + Timestamp
Gateway -> 用户A: WS: 消息发送成功

RabbitMQ (message.push) -> Gateway: 消息推送
Gateway -> Redis: 查询用户B在线状态
Gateway -> 用户B: WS: new_message

RabbitMQ (message.index) -> Search RPC: 消息索引
Search RPC -> Elasticsearch: 索引消息
```

### 4.4 好友申请流程

```
用户A -> Gateway: POST /api/v1/friends/request
Gateway -> Friend RPC: SendFriendRequest
Friend RPC -> User RPC: 验证用户B存在
Friend RPC -> PostgreSQL: 创建好友申请
Friend RPC -> Gateway: RequestID
Gateway -> 用户A: 申请已发送
Gateway -> 用户B: WS: friend_request (推送通知)

用户B -> Gateway: POST /api/v1/friends/request/handle
Gateway -> Friend RPC: HandleFriendRequest
Friend RPC -> PostgreSQL: 更新申请状态
Friend RPC -> PostgreSQL: 创建双向好友关系 (如果同意)
Friend RPC -> Redis: 删除好友列表缓存
Friend RPC -> Gateway: Success
Gateway -> 用户B: 处理成功
Gateway -> 用户A: WS: friend_accepted (推送通知)
```

### 4.5 文件上传流程

```
用户 -> Gateway: POST /api/v1/files/upload (带文件数据)
Gateway -> File RPC: UploadFile (含 file_hash)
File RPC -> PostgreSQL: 查询 file_hash 是否存在
  存在:
    File RPC -> PostgreSQL: 引用计数+1
    File RPC -> Gateway: 返回已有文件 URL (去重成功)
  不存在:
    File RPC -> File System/OSS: 保存文件
    File RPC -> PostgreSQL: 插入文件记录
    File RPC -> Gateway: 返回新文件 URL
Gateway -> 用户: 上传成功
```

## 五、高可用设计

### 5.1 服务高可用

- **RPC 服务**: 每个 RPC 服务部署多个实例，通过 etcd 服务发现实现负载均衡
- **Gateway**: 部署多个实例，通过 Nginx/K8s LoadBalancer 负载均衡
- **数据库**: PostgreSQL 主从复制，读写分离
- **缓存**: Redis Cluster 或 Redis Sentinel
- **消息队列**: RabbitMQ 集群
- **搜索引擎**: Elasticsearch 集群

### 5.2 故障处理

- **熔断降级**: go-zero 自带自适应熔断器
- **超时控制**: 所有 RPC 调用设置合理超时时间
- **重试机制**: 对幂等接口实现自动重试
- **限流**: API Gateway 实现令牌桶限流
- **服务降级**: 非核心功能（如搜索）降级不影响主流程

### 5.3 数据一致性

- **最终一致性**: 消息队列保证消息最终被处理
- **分布式事务**: 不使用强一致性事务，采用最终一致性 + 补偿机制
- **缓存一致性**: 更新数据库后删除缓存，延迟双删策略

## 六、安全设计

### 6.1 认证鉴权

- **JWT Token**: 所有需要认证的接口使用 JWT
- **Token 过期**: 7天过期，支持 Refresh Token
- **密码加密**: bcrypt 加密存储
- **验证码**: 5分钟过期，单次使用

### 6.2 数据安全

- **SQL 注入防护**: 参数化查询
- **XSS 防护**: 输入校验和转义
- **CSRF 防护**: CSRF Token
- **敏感信息加密**: 数据库敏感字段加密存储

### 6.3 接口安全

- **限流**: 防止 DDoS 攻击
- **黑名单**: IP 黑名单机制
- **请求签名**: 关键接口添加签名验证

## 七、性能优化

### 7.1 缓存策略

| 缓存内容 | 缓存键 | TTL | 说明 |
|---------|--------|-----|------|
| 用户信息 | user:info:{user_id} | 1小时 | 减少数据库查询 |
| 用户在线状态 | user:online:{user_id} | 5分钟 | 实时性要求高 |
| 会话列表 | conv:list:{user_id} | 30分钟 | 会话变化不频繁 |
| 会话详情 | conv:detail:{conv_id} | 1小时 | 会话信息相对稳定 |
| 好友列表 | friend:list:{user_id} | 1小时 | 好友关系变化不频繁 |
| 热点消息 | msg:hot:{conv_id} | 10分钟 | 最近100条消息 |

### 7.2 数据库优化

- **索引优化**: 为高频查询字段添加索引
- **分表分区**: 消息表按月分区
- **读写分离**: 主库写，从库读
- **连接池**: 合理配置连接池大小
- **慢查询监控**: 监控并优化慢查询

### 7.3 消息队列优化

- **批量处理**: 批量消费消息，减少数据库写入次数
- **消息持久化**: 重要消息持久化到磁盘
- **死信队列**: 处理失败的消息进入死信队列

### 7.4 WebSocket 优化

- **连接池管理**: 使用高效的连接池
- **心跳检测**: 定时心跳检测，清理无效连接
- **消息压缩**: WebSocket 消息使用 gzip 压缩
- **分布式连接**: 多个 Gateway 实例共享连接信息（Redis）

## 八、监控告警

### 8.1 监控指标

- **服务监控**: QPS、响应时间、错误率
- **资源监控**: CPU、内存、磁盘、网络
- **业务监控**: 在线用户数、消息发送量、接口调用量
- **数据库监控**: 连接数、慢查询、锁等待

### 8.2 日志系统

- **日志级别**: DEBUG、INFO、WARN、ERROR
- **链路追踪**: go-zero 集成 OpenTelemetry
- **日志收集**: ELK (Elasticsearch + Logstash + Kibana)

### 8.3 告警规则

- **服务告警**: 服务不可用、响应时间 > 1s、错误率 > 1%
- **资源告警**: CPU > 80%、内存 > 80%、磁盘 > 90%
- **业务告警**: 消息积压 > 1000、用户投诉

## 九、技术亮点

1. **微服务架构**: 服务拆分合理，职责清晰，易于扩展
2. **自动代码生成**: 使用 goctl 自动生成代码，提高开发效率
3. **服务治理**: 集成限流、熔断、降级、链路追踪等微服务治理能力
4. **消息队列解耦**: 使用 RabbitMQ 实现服务解耦和异步处理
5. **全文检索**: Elasticsearch 实现历史消息全文检索
6. **文件去重**: SHA256 哈希去重，节省存储空间
7. **缓存设计**: 合理的缓存策略，提升系统性能
8. **WebSocket 长连接**: 实时消息推送
9. **高可用设计**: 服务高可用、数据高可用、故障自动恢复

## 十、项目目录结构

```
Beehive/
├── api/
│   ├── gateway/               # API Gateway
│   │   ├── etc/
│   │   │   └── gateway.yaml
│   │   ├── internal/
│   │   │   ├── config/       # 配置
│   │   │   ├── handler/      # HTTP Handler
│   │   │   ├── logic/        # 业务逻辑
│   │   │   ├── middleware/   # JWT、Auth 中间件
│   │   │   ├── svc/          # 服务依赖
│   │   │   ├── types/        # 请求响应结构体
│   │   │   └── ws/           # WebSocket 处理
│   │   └── gateway.go
│   └── proto/                # Proto 定义
│       ├── user/
│       ├── friend/
│       ├── chat/
│       ├── message/
│       ├── file/
│       └── search/
├── rpc/
│   ├── user/                 # User RPC Service
│   │   ├── etc/
│   │   ├── internal/
│   │   │   ├── config/
│   │   │   ├── logic/
│   │   │   ├── server/
│   │   │   └── svc/
│   │   ├── model/            # goctl model 生成
│   │   ├── user/             # pb.go
│   │   └── user.go
│   ├── friend/               # Friend RPC Service
│   ├── chat/                 # Chat RPC Service
│   ├── message/              # Message RPC Service
│   ├── file/                 # File RPC Service
│   └── search/               # Search RPC Service
├── common/
│   ├── utils/                # 工具函数
│   ├── errorx/               # 错误处理
│   └── middleware/           # 公共中间件
├── docker/
│   ├── docker-compose.yml    # 所有基础设施
│   └── etcd/
├── docs/
│   ├── dev/                  # 开发文档
│   │   ├── architecture.md   # 架构文档（本文件）
│   │   ├── database.md       # 数据库设计
│   │   ├── api.md            # API 文档
│   │   ├── rpc.md            # RPC 文档
│   │   ├── message-queue.md  # 消息队列文档
│   │   ├── elasticsearch.md  # Elasticsearch 文档
│   │   └── deployment.md     # 部署文档
│   └── api/                  # API 接口文档
├── scripts/
│   ├── init_db.sql           # 初始化数据库脚本
│   ├── init_es.sh            # 初始化 ES 索引
│   └── gen_code.sh           # 代码生成脚本
├── go.mod
└── README.md
```

## 十一、开发规范

### 11.1 代码规范

- 遵循 Go 官方代码规范
- 使用 gofmt 格式化代码
- 使用 golangci-lint 进行代码检查
- 统一使用 go-zero 框架规范

### 11.2 Git 提交规范

遵循 Angular 提交规范：

- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 重构
- `test`: 测试
- `chore`: 构建/工具链

### 11.3 接口命名规范

- RESTful API: 使用名词复数，如 `/api/v1/users`
- RPC 方法: 使用动词开头，如 `GetUserInfo`、`SendMessage`

## 十二、后续优化方向

1. **前端开发**: React + TypeScript + Electron 桌面端
2. **移动端**: React Native 或 Flutter
3. **音视频通话**: WebRTC 集成
4. **消息撤回**: 2分钟内消息撤回
5. **表情包**: 自定义表情包
6. **消息已读回执**: 显示消息已读状态
7. **群公告**: 群聊公告功能
8. **@提醒**: 群聊 @ 某人
9. **消息转发**: 消息转发到其他会话
10. **文件夹分组**: 好友分组管理
