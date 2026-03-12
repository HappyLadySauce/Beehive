## Beehive 新架构总览（仅 WebSocket/HTTP）

### 目标

- **统一协议**：客户端只通过 WebSocket/HTTP 访问系统，不再暴露原来的自定义 TCP 协议和自研二进制包格式。
- **减少自研轮子**：尽量使用 go-zero + 成熟第三方库（如 gorilla/websocket、标准 JWT 库、官方 Redis/PostgreSQL/RabbitMQ 客户端）。
- **职责清晰**：以 Gateway 为入口，将认证、在线状态、消息、会话、通知等拆分成独立服务。
- **易扩展与高可用**：各服务独立部署、独立扩缩容，通过 gRPC 与消息队列解耦。

### 对外入口层

- **GatewayService（WebSocket 网关）**
  - 暴露统一入口：`/ws` WebSocket，以及必要的 HTTP 健康检查接口。
  - 只处理：连接管理（握手、心跳、关闭）、基础认证（调用 AuthService 校验 token）、限流与简单路由。
  - 使用 **JSON over WebSocket** 作为消息格式，后续可按需升级为 Proto over WebSocket。

- **AdminAPIService（管理 HTTP 接口，可后期补充）**
  - 使用 go-zero `rest` 实现，提供运营/诊断能力（查询用户、连接、会话、消息等）。

### 内部核心业务服务

- **AuthService**
  - 用户登录、注册（可选）、密码校验。
  - 使用成熟 JWT 库签发和校验访问令牌。
  - 持久化在 PostgreSQL，支持后续扩展为多种认证方式。

- **UserService**
  - 管理用户基础资料（昵称、头像、扩展信息）。
  - 主要通过 PostgreSQL 存储，Redis 做缓存。

- **PresenceService**
  - 负责在线状态与多端 Session 管理。
  - 使用 Redis 存储 userId → 连接/设备 列表，定期基于心跳刷新 TTL。

- **MessageService**
  - 处理点对点和群聊消息的写入与查询（历史消息拉取）。
  - 使用 PostgreSQL 持久化消息与会话相关数据。
  - 在写入成功后，通过 RabbitMQ 发布消息事件，供投递和统计等下游使用。

- **ConversationService**
  - 管理会话/群组的元数据与成员关系（单聊、群聊、频道等）。
  - 使用 PostgreSQL 存储，Redis 可缓存用户常用会话列表。

- **NotificationService / AnalyticsService（可选）**
  - 订阅 RabbitMQ 上的事件流，实现离线推送、统计与审计等功能。

### 通信模式

- **客户端 ↔ Gateway**
  - WebSocket：`wss://host/ws`。
  - 统一 JSON 消息格式，例如：

    ```json
    {
      "type": "message.send",
      "tid": "client-msg-id",
      "payload": {
        "conversationId": "conv_123",
        "toUserId": "u_456",
        "body": {
          "type": "text",
          "text": "hello"
        }
      }
    }
    ```

- **服务间同步调用（RPC）**
  - 使用 go-zero `zrpc` + gRPC。
  - 为 Auth / Presence / Message / Conversation 分别设计独立的 proto 文件。

- **服务间异步事件（MQ）**
  - 使用 RabbitMQ 作为事件总线。
  - MessageService 在消息写入成功后向交换机（如 `im.events`）发布事件，Delivery/Notification/Analytics 等服务订阅。

