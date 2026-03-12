## GatewayService 设计（WebSocket 网关）

本节详细说明 Gateway 的职责、连接模型、消息协议以及与后端服务（Auth/Presence/Message/Conversation）之间的交互方式。

> 约束：只使用 WebSocket/HTTP，不再暴露自定义 TCP 协议和自研二进制包格式；尽量依赖成熟三方库（如 gorilla/websocket、JWT 库等）。

---

### 1. 职责边界

- **负责的事情**
  - 统一客户端接入：`/ws` WebSocket 连接的建立、关闭。
  - 连接管理：保持每个连接的基本状态（关联 userId、设备信息、最后心跳时间）。
  - 消息编解码：将 WebSocket 文本/二进制帧解码为统一 JSON 消息对象，做基础校验。
  - 认证入口：
    - 首次登录：接受 `auth.login` / `auth.tokenLogin` 消息，调用 AuthService 完成账号密码/Token 校验。
    - 后续请求：校验连接上是否已绑定 userId + token 是否仍有效（可定期向 AuthService 校验）。
  - 路由与编排：
    - 根据消息 `type` 将请求路由到对应 handler。
    - handler 内部调用相应后端服务（Auth / Presence / Message / Conversation）。
  - 推送消息给在线客户端：
    - 自己实例上本地连接的推送。
    - 配合 DeliveryService/PresenceService 对其他实例在线的用户推送（可通过 MQ 或内部 RPC）。
  - 基础防护与限流：
    - 登录频率限制、消息速率限制（依赖 Redis）。
    - 简单黑名单/封禁控制（依赖 AdminAPI / UserService）。

- **不负责的事情**
  - 不负责业务数据的持久化（用户、会话、消息等都不写库）。
  - 不直接访问 PostgreSQL（只通过后端服务访问）。
  - 不直接操作 RabbitMQ（消息事件由 MessageService/DeliveryService 负责发布/消费）。

---

### 2. 连接与 Session 模型

- **连接（Connection）**
  - 抽象为：`ConnID`（本地唯一）、`UserID`（登录后绑定）、`DeviceID`（可选，来自客户端）、`GatewayID`（实例 ID，用于多实例部署）。
  - 连接生命周期：
    1. WebSocket 握手成功，分配 `ConnID`，尚未登录。
    2. 收到 `auth.login` / `auth.tokenLogin` 消息，调用 AuthService 验证。
    3. 登录成功后，绑定 `UserID`，调用 PresenceService 的 `RegisterSession`。
    4. 连接期间，通过 `presence.ping`/心跳更新在线状态（PresenceService 刷新 TTL）。
    5. 连接关闭时，调用 PresenceService 的 `UnregisterSession`。

- **PresenceService 交互**
  - 注册：`RegisterSession(userId, gatewayId, connId, deviceInfo)`。
  - 注销：`UnregisterSession(userId, connId)`。
  - 查询：`GetOnlineSessions(userId)` 给 DeliveryService/MessageService 使用。
  - 数据落地在 Redis 中，Gateway 只通过 RPC 访问。

---

### 3. 消息协议（JSON over WebSocket）

#### 3.1 顶层结构

采用统一的 JSON 消息 Envelope，兼容请求和响应：

```json
{
  "type": "message.send",
  "tid": "client-generated-id",
  "payload": { },
  "error": null
}
```

- **`type`**：字符串枚举，决定路由和 payload 结构，例如：
  - 认证相关：`auth.login`, `auth.tokenLogin`, `auth.logout`
  - 心跳相关：`presence.ping`
  - 消息相关：`message.send`, `message.ack`, `message.history`
  - 会话相关：`conversation.create`, `conversation.list`, `conversation.members`
- **`tid`**（Trace/Transaction ID，可选但推荐）：
  - 由客户端生成，用于在响应中对应请求。
  - 服务端在响应中回显，用于前端定位和重试幂等。
- **`payload`**：具体业务字段，结构依赖于 `type`。
- **`error`**：服务端响应时提供错误信息，格式示例：

  ```json
  "error": {
    "code": "unauthorized",
    "message": "invalid token"
  }
  ```

#### 3.2 示例：登录与消息发送

- 登录请求：

  ```json
  {
    "type": "auth.login",
    "tid": "login-1",
    "payload": {
      "username": "alice",
      "password": "password123",
      "deviceId": "ios-uuid"
    }
  }
  ```

- 登录成功响应：

  ```json
  {
    "type": "auth.login.ok",
    "tid": "login-1",
    "payload": {
      "userId": "u_123",
      "accessToken": "jwt-token",
      "refreshToken": "refresh-token",
      "expiresIn": 3600
    },
    "error": null
  }
  ```

- 发送消息请求：

  ```json
  {
    "type": "message.send",
    "tid": "msg-1",
    "payload": {
      "clientMsgId": "local-uuid",
      "conversationId": "conv_abc",
      "toUserId": "u_456",
      "body": {
        "type": "text",
        "text": "hello"
      }
    }
  }
  ```

- 发送消息响应（仅确认写库/受理）：

  ```json
  {
    "type": "message.send.ok",
    "tid": "msg-1",
    "payload": {
      "serverMsgId": "msg_789",
      "serverTime": 1710000000
    },
    "error": null
  }
  ```

---

### 4. Gateway 与后端服务的交互流程

#### 4.1 登录流程（Gateway ↔ AuthService ↔ PresenceService）

1. 客户端通过 `/ws` 建立连接。
2. 发送 `auth.login` 消息。
3. Gateway 解析消息，调用 AuthService 的 `Login` RPC：
   - 校验账号密码。
   - 生成 accessToken/refreshToken，写入 PostgreSQL/缓存。
4. AuthService 返回 `userId` 与 token 信息。
5. Gateway 将连接与 `userId` 绑定，并调用 PresenceService 的 `RegisterSession`：
   - PresenceService 在 Redis 中记录 `user:{userId}:sessions`。
6. Gateway 向客户端返回 `auth.login.ok` 响应。

#### 4.2 发送消息流程（Gateway ↔ MessageService ↔ RabbitMQ ↔ Delivery/Notification）

1. 客户端发送 `message.send`。
2. Gateway 校验连接是否已登录（有 `userId`）、速率是否超限（基于 Redis）。
3. Gateway 调用 MessageService 的 `PostMessage`：
   - 写入 PostgreSQL；
   - 成功后在 RabbitMQ 交换机 `im.events` 上发布 `message.created` 事件。
4. MessageService 返回 `serverMsgId` 等信息，Gateway 返回 `message.send.ok` 给客户端。
5. DeliveryService 消费 `message.created`：
   - 通过 ConversationService/PresenceService 获取目标用户在线 session。
   - 调用各 Gateway 实例内部接口/推送通道，将消息实时发给在线用户。
6. NotificationService 可同时消费该事件，对离线用户进行推送。

---

### 5. 水平扩展与实例间协同

- Gateway 为 **无状态服务**（仅本地连接状态）。
- 会话/在线状态由 PresenceService + Redis 统一管理。
- 消息的持久化与广播由 MessageService + RabbitMQ + DeliveryService 负责。
- Gateway 实例之间不直接通信，通过：
  - Redis（PresenceService）感知用户在哪个 Gateway 实例在线。
  - RabbitMQ（消息事件）进行跨实例的消息分发。

