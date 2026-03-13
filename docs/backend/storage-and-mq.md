## PostgreSQL / Redis / RabbitMQ / etcd 的职责与使用方式

本节明确关键基础设施组件在 Beehive 新架构中的定位和边界：PostgreSQL（主存储）、Redis（在线状态与缓存）、RabbitMQ（事件总线与异步处理）、etcd（服务发现与配置管理）。

---

### 1. PostgreSQL：强一致业务数据的主存储

**定位**

- 作为系统的 **system-of-record**，存放需要长期保存、支持复杂查询和事务语义的数据。
- 所有关键业务数据最终都必须持久化到 PostgreSQL 中。

**主要存放内容**

- 用户与认证：
  - `users`：账号、密码哈希、状态（正常/封禁）、创建时间等。
  - `user_profiles`：昵称、头像、简介、扩展字段等。
- 会话与群组：
  - `conversations`：会话/群组/频道信息（类型、标题、创建者等）。
  - `conversation_members`：成员关系表（会话 ID + 用户 ID + 角色/状态）。
- 消息与相关索引：
  - `messages`：点对点/群聊消息体，关键字段如 `conversation_id`、`from_user_id`、`server_time`、`body` 等。

**使用建议**

- 使用 go-zero 自带 model 或 GORM 等成熟 ORM，不自写通用 ORM 框架。
- 针对典型查询场景设计索引，例如：
  - 按 `conversation_id + server_time` 分页拉取历史消息。
  - 按 `user_id` 获取最近会话列表（可配合汇总表或物化视图）。
- 初期以单库单实例为主，后续根据规模再考虑分库分表与读写分离。

---

### 2. Redis：在线状态、缓存与限流

**定位**

- 高性能内存存储，用于短期状态、高频读写和限流控制。
- 非关键持久化存储，不应单独作为业务真相来源。

**主要使用场景**

- 在线状态 & Session：
  - PresenceService 维护：
    - `user:{userId}:sessions` → set/list of `{gatewayId, connId, deviceInfo}`。
  - 每次心跳或业务消息可刷新 TTL，避免僵尸连接。
- 缓存：
  - 用户信息缓存：`user:profile:{userId}`。
  - 会话列表缓存：`user:{userId}:conversations:recent`。
- 限流与防刷：
  - 以 IP 或 userId 为维度，使用 Redis 计数器/滑动窗口实现请求限流：
    - 如 `rate:login:{userId}`、`rate:msg:{userId}`。

**不做的事情**

- 不将消息内容只存 Redis，不依赖 Redis 作为唯一消息存储。
- 不在 Redis 中存过多历史数据，避免内存压力。

---

### 3. RabbitMQ：事件总线与异步消息处理

**定位**

- 可靠消息队列 / 事件总线，用于解耦服务、实现异步处理和大规模 fan-out。
- 在消息写入数据库后，通过事件流将变化推送到 Delivery、Notification、Analytics 等服务。

**主要使用场景**

- 消息生命周期事件：
  - MessageService 在成功写消息到 PostgreSQL 后，向 RabbitMQ 交换机（如 `im.events`）发送 `message.created` 事件。
  - 下游服务（DeliveryService、NotificationService、AnalyticsService）订阅相应 routing key。
- 在线消息投递（推荐方案）：
  - Gateway 接收 `message.send` 后调用 MessageService → 写库成功。
  - MessageService 发送 `message.created` 事件。
  - DeliveryService 消费事件，根据 PresenceService 提供的在线列表决定需要推送到哪些 Gateway 实例。
  - DeliveryService 调用 Gateway 内部接口/推送通道，将消息扇出到各在线用户的 WebSocket 连接。
- 离线通知、审计与统计：
  - NotificationService 可订阅 `message.created`、`message.read` 等事件，进行推送或写入审计表。
  - AnalyticsService 消费事件，生成统计报表与监控指标。

**队列与路由设计示例**

- 交换机：`im.events`（类型：`topic`）
- 路由键示例：
  - `message.created.direct`：即时消息投递。
  - `message.created.analytics`：统计与日志分析。
  - `message.read.direct`：已读回执相关逻辑。
- 不同下游服务可以绑定不同的 routing key，互不影响。

---

### 4. etcd：服务发现与配置管理

**定位**

- 作为集群级别的 **服务注册/发现中心**，用于各服务之间的地址发现与负载均衡。
- 作为统一的 **配置管理** 存储，用于存放跨环境的动态配置（开关、灰度、限流阈值等），而不是业务数据。

**主要使用场景**

- 服务注册与发现：
  - 各个服务在启动时向 etcd 注册自己的实例信息（服务名、IP、端口、权重等）。
  - Gateway、DeliveryService 等消费者通过 etcd 获取 MessageService/AuthService 等下游服务的地址列表，实现客户端侧负载均衡。
- 配置中心：
  - 存储系统级配置，如：
    - 某些功能的开关（feature flag）。
    - 动态限流阈值。
    - 特定环境的灰度规则等。
  - 服务在启动时加载一次，并可通过 Watch 机制感知配置变更。

**不做的事情**

- 不在 etcd 中存储业务数据（用户、消息、会话等）。
- 不将 etcd 用作高 QPS 缓存或队列。

---

### 5. 各服务与基础设施的关系汇总

- **AuthService**
  - PostgreSQL：`users` 等认证数据，以及系统级 RBAC 相关表（如 `roles` / `permissions` / `role_permissions` / `user_roles`）。
  - Redis：可缓存用户登录信息、黑名单等（可选）。
  - etcd：参与服务注册发现，获取其他依赖组件的地址（如配置中心自身也可使用 etcd 存储）。
- **UserService**
  - PostgreSQL：用户资料。
  - Redis：用户资料缓存。
  - etcd：服务注册发现。
- **PresenceService**
  - Redis：在线状态和 session 映射（核心依赖）。
  - etcd：服务注册发现。
- **MessageService**
  - PostgreSQL：消息与会话相关持久化。
  - RabbitMQ：写库成功后发布消息事件。
  - etcd：服务注册发现。
- **ConversationService**
  - PostgreSQL：会话/群组/成员关系。
  - Redis：会话列表与成员列表缓存（可选）。
  - etcd：服务注册发现。
- **GatewayService（WebSocket）**
  - Redis：读取/更新在线状态（通过 PresenceService），实现限流计数等。
  - RabbitMQ：通常不直接访问，由 Delivery/MessageService 连接；Gateway 通过内部调度获得推送任务。
  - etcd：发现各后端服务实例地址，必要时读取部分动态配置。
- **NotificationService / AnalyticsService**
  - RabbitMQ：消费事件流。
  - PostgreSQL/其他存储：落地统计或审计结果。
  - etcd：服务注册发现。

