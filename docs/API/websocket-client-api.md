## WebSocket 客户端 API（Gateway `/ws`）

本文件定义 IM Web 客户端（以及未来的移动/桌面客户端）通过 WebSocket 与 Beehive Gateway 交互时使用的统一 JSON 协议。

设计遵循 `docs/backend/gateway-design.md` 中的约定，仅补充更具体的字段和消息类型。

**ID 约定**：

- **用户 ID**：10 位数字字符串（如 `"1234567890"`），用于登录返回、会话成员、消息 from/to 等。
- **会话 ID（conversationId）**：单聊为 UUID 字符串；群聊为 11 位数字字符串（群号），用于创建群聊、加群、发消息等。

---

### 1. 顶层 Envelope 结构

所有 WebSocket 消息（请求/响应/推送）都使用统一的 JSON Envelope：

```json
{
  "type": "message.send",
  "tid": "client-generated-id",
  "payload": {},
  "error": null
}
```

- **`type`**：字符串枚举，标识消息类别与方向。约定命名规则：
  - 请求：`<domain>.<action>`，例如：`auth.login`, `message.send`, `conversation.list`
  - 成功响应：`<domain>.<action>.ok`，例如：`auth.login.ok`, `message.send.ok`
  - 失败响应：`<domain>.<action>.error`，例如：`auth.login.error`, `message.send.error`
  - 推送/事件：`<domain>.<event>`，例如：`message.push`, `conversation.updated`
- **`tid`**（Trace/Transaction ID，可选但推荐）：
  - 由客户端生成，在请求与对应响应中保持一致，用于前端匹配请求与响应、排查问题。
  - 服务端在响应/相关事件中应尽量原样回传 `tid`。
- **`payload`**：具体业务数据，结构根据 `type` 不同而不同，见下文各类消息定义。
- **`error`**：仅在响应/推送中填充错误信息，请求消息中应为空或省略。格式：

```json
{
  "code": "unauthorized",
  "message": "invalid token"
}
```

常见错误码（示例）：

- `bad_request`：请求格式错误 / 缺少字段
- `unauthorized`：未登录或 token 无效
- `forbidden`：被封禁 / 无操作权限
- `rate_limited`：触发限流
- `not_found`：目标资源不存在
- `internal_error`：服务器内部错误

---

### 2. 连接与认证相关消息

#### 2.1 使用 token 登录（推荐）

- **请求：`auth.tokenLogin`**

```json
{
  "type": "auth.tokenLogin",
  "tid": "login-1",
  "payload": {
    "accessToken": "jwt-token",
    "deviceId": "web-uuid-or-user-agent"
  }
}
```

- 字段说明：
  - `accessToken`：从 HTTP 登录（或其他入口）获取的 JWT 访问令牌。
  - `deviceId`：可选，用于区分不同设备/浏览器实例。

- **成功响应：`auth.tokenLogin.ok`**

```json
{
  "type": "auth.tokenLogin.ok",
  "tid": "login-1",
  "payload": {
    "userId": "u_123",
    "expiresIn": 3600
  },
  "error": null
}
```

- **失败响应示例**

```json
{
  "type": "auth.tokenLogin.error",
  "tid": "login-1",
  "payload": null,
  "error": {
    "code": "unauthorized",
    "message": "invalid or expired token"
  }
}
```

#### 2.2 账号密码登录（可选，不一定在 WebSocket 上暴露）

如需直接通过 WebSocket 完成登录，可以定义：

- **请求：`auth.login`**

```json
{
  "type": "auth.login",
  "tid": "login-2",
  "payload": {
    "username": "alice",
    "password": "password123",
    "deviceId": "web-uuid-or-user-agent"
  }
}
```

- **成功响应：`auth.login.ok`**

```json
{
  "type": "auth.login.ok",
  "tid": "login-2",
  "payload": {
    "userId": "u_123",
    "accessToken": "jwt-token",
    "refreshToken": "refresh-token",
    "expiresIn": 3600
  },
  "error": null
}
```

> 说明：实际实现中，可以选择「HTTP 登录 + WebSocket tokenLogin」的模式，或两者兼容。文档这里给出协议定义，不强制具体部署方式。

#### 2.3 注册

- **请求：`auth.register`**

```json
{
  "type": "auth.register",
  "tid": "reg-1",
  "payload": {
    "username": "alice",
    "password": "password123"
  }
}
```

- **成功响应：`auth.register.ok`**（与登录一致，返回 token，注册即登录）

```json
{
  "type": "auth.register.ok",
  "tid": "reg-1",
  "payload": {
    "userId": "u_123",
    "accessToken": "jwt-token",
    "refreshToken": "refresh-token",
    "expiresIn": 3600
  },
  "error": null
}
```

- **失败响应**：如用户名已存在，返回 `auth.register.error`，`error.code` 为 `bad_request`，`error.message` 如 "username already exists"。

#### 2.4 登出

- **请求：`auth.logout`**

```json
{
  "type": "auth.logout",
  "tid": "logout-1",
  "payload": {}
}
```

- **成功响应：`auth.logout.ok`**

```json
{
  "type": "auth.logout.ok",
  "tid": "logout-1",
  "payload": {},
  "error": null
}
```

---

### 3. 心跳与在线状态

客户端需要定期发送心跳，维持在线状态。心跳可以简单设计为：

- **请求：`presence.ping`**

```json
{
  "type": "presence.ping",
  "tid": "ping-1",
  "payload": {
    "clientTime": 1710000000
  }
}
```

- **响应：`presence.ping.ok`**

```json
{
  "type": "presence.ping.ok",
  "tid": "ping-1",
  "payload": {
    "serverTime": 1710000001
  },
  "error": null
}
```

> Gateway 收到 `presence.ping` 后，会通过 PresenceService 刷新 Redis 中对应 session 的 TTL。

---

### 4. 消息发送与推送

#### 4.1 发送消息（单聊/群聊通用）

- **请求：`message.send`**

```json
{
  "type": "message.send",
  "tid": "msg-1",
  "payload": {
    "clientMsgId": "local-uuid-1",
    "conversationId": "conv_abc",      // 可选，若为空则视为新对话或从 toUserId 推断
    "toUserId": "u_456",               // 单聊时可用
    "body": {
      "type": "text",
      "text": "hello"
    }
  }
}
```

- 字段说明：
  - `clientMsgId`：客户端生成的消息 ID，用于幂等与本地状态更新。
  - `conversationId`：
    - 已知会话时必须传递；
    - 若为空且 `toUserId` 存在，后端可创建或查找对应单聊会话。
  - `toUserId`：点对点消息时的目标用户 ID；群聊/频道消息可以只依赖 `conversationId`。
  - `body`：
    - `type`: `"text" | "image" | "system" | ..."`
    - 其余字段根据 type 扩展，这里以文本为例。

- **成功响应：`message.send.ok`**

```json
{
  "type": "message.send.ok",
  "tid": "msg-1",
  "payload": {
    "serverMsgId": "msg_789",
    "serverTime": 1710000000,
    "conversationId": "conv_abc"
  },
  "error": null
}
```

> 说明：该响应仅表示「写库 + 事件发布」成功。真正的消息推送给自己/对端可通过 `message.push` 完成。

#### 4.2 消息推送（下行）

当有新消息需要投递给某个在线用户时，Gateway 会向其 WebSocket 连接推送：

- **推送：`message.push`**

```json
{
  "type": "message.push",
  "tid": "msg-1",   // 若来自某个发送请求，可沿用原 tid；否则可为空或新值
  "payload": {
    "serverMsgId": "msg_789",
    "clientMsgId": "local-uuid-1",   // 对于自己发的消息可回显
    "conversationId": "conv_abc",
    "fromUserId": "u_123",
    "toUserId": "u_456",
    "body": {
      "type": "text",
      "text": "hello"
    },
    "serverTime": 1710000000
  },
  "error": null
}
```

客户端处理：

- 将消息追加到对应会话的消息列表；
- 对非当前会话增加未读计数；
- 若 `clientMsgId` 匹配本地 pending 消息，更新其状态为「已发送」并替换为 `serverMsgId`。

#### 4.3 已读回执（可选）

- **请求：`message.read`**

```json
{
  "type": "message.read",
  "tid": "read-1",
  "payload": {
    "conversationId": "conv_abc",
    "serverMsgId": "msg_789"
  }
}
```

- **响应：`message.read.ok`**

```json
{
  "type": "message.read.ok",
  "tid": "read-1",
  "payload": {},
  "error": null
}
```

> 后端可异步产生 `message.read` 相关事件，用于同步已读状态给对端或统计。

---

### 5. 历史消息与会话列表

#### 5.1 拉取会话列表

- **请求：`conversation.list`**

```json
{
  "type": "conversation.list",
  "tid": "conv-list-1",
  "payload": {
    "cursor": null,
    "limit": 50
  }
}
```

- **响应：`conversation.list.ok`**

```json
{
  "type": "conversation.list.ok",
  "tid": "conv-list-1",
  "payload": {
    "items": [
      {
        "id": "conv_abc",
        "name": "Alice",
        "avatar": "",
        "type": "single",              // single | group | channel
        "lastMessage": {
          "serverMsgId": "msg_789",
          "preview": "hello",
          "serverTime": 1710000000
        },
        "unreadCount": 3,
        "lastActiveAt": 1710000000
      }
    ],
    "nextCursor": null
  },
  "error": null
}
```

> 该接口由 Gateway 调用 ConversationService + MessageService 聚合得到，前端用于会话侧边栏展示。

#### 5.2 拉取历史消息

- **请求：`message.history`**

```json
{
  "type": "message.history",
  "tid": "history-1",
  "payload": {
    "conversationId": "conv_abc",
    "before": 1710000000,   // 可选：某个 serverTime（int64 时间戳）之前
    "limit": 50
  }
}
```

- **响应：`message.history.ok`**

```json
{
  "type": "message.history.ok",
  "tid": "history-1",
  "payload": {
    "items": [
      {
        "serverMsgId": "msg_780",
        "clientMsgId": null,
        "conversationId": "conv_abc",
        "fromUserId": "u_123",
        "toUserId": "u_456",
        "body": {
          "type": "text",
          "text": "old message"
        },
        "serverTime": 1709999900
      }
    ],
    "hasMore": true
  },
  "error": null
}
```

---

### 6. 用户资料与会话、联系人

#### 6.1 获取当前用户资料

- **请求：`user.me`**（需已登录）

```json
{
  "type": "user.me",
  "tid": "me-1",
  "payload": {}
}
```

- **成功响应：`user.me.ok`**

```json
{
  "type": "user.me.ok",
  "tid": "me-1",
  "payload": {
    "id": "1234567890",
    "nickname": "Alice",
    "avatarUrl": "https://...",
    "bio": "hello",
    "status": "normal"
  },
  "error": null
}
```

> Gateway 根据连接上绑定的 `userId` 调用 **UserService.GetUser** 返回当前用户资料，供客户端展示昵称、头像等。

#### 6.2 会话与联系人管理

- **创建会话：`conversation.create` / `conversation.create.ok`**

请求（单聊可传 `memberIds` 或 `toUsername`/`toAccount` 之一）：

```json
{
  "type": "conversation.create",
  "tid": "create-1",
  "payload": {
    "type": "single",
    "name": "",
    "memberIds": ["1234567890", "0987654321"],
    "toUsername": "alice",
    "toAccount": "1234567890"
  }
}
```

- `type`：会话类型，`single` | `group` | `channel`
- `name`：可选，群名/频道名（群聊时使用）
- `memberIds`：成员用户 ID（10 位）数组；单聊时可与 `toUsername`/`toAccount` 二选一
- `toUsername`：单聊时按用户名解析对方，与 `toAccount`、`memberIds` 互斥
- `toAccount`：单聊时按 10 位账号解析对方，与 `toUsername`、`memberIds` 互斥

成功响应：单聊返回 UUID 格式 `conversationId`，群聊返回 11 位群号。

```json
{
  "type": "conversation.create.ok",
  "tid": "create-1",
  "payload": {
    "conversationId": "conv_uuid_or_11_digit_group_id"
  },
  "error": null
}
```

- **添加成员：`conversation.addMember` / `conversation.addMember.ok`**

请求：

```json
{
  "type": "conversation.addMember",
  "tid": "add-1",
  "payload": {
    "conversationId": "conv_abc",
    "userId": "u_789",
    "role": "member"
  }
}
```

- `role`：可选，默认 `member`，可为 `owner` | `admin` | `member`

成功响应：

```json
{
  "type": "conversation.addMember.ok",
  "tid": "add-1",
  "payload": {},
  "error": null
}
```

- **移除成员：`conversation.removeMember` / `conversation.removeMember.ok`**

请求：

```json
{
  "type": "conversation.removeMember",
  "tid": "rm-1",
  "payload": {
    "conversationId": "conv_abc",
    "userId": "u_789"
  }
}
```

成功响应：

```json
{
  "type": "conversation.removeMember.ok",
  "tid": "rm-1",
  "payload": {},
  "error": null
}
```

- **联系人：`contact.list` / `contact.add` / `contact.remove`**

  - **请求：`contact.list`**

    ```json
    { "type": "contact.list", "tid": "cl-1", "payload": {} }
    ```

  - **成功响应：`contact.list.ok`**

    ```json
    {
      "type": "contact.list.ok",
      "tid": "cl-1",
      "payload": { "contactUserIds": ["1234567890", "0987654321"] },
      "error": null
    }
    ```

  - **请求：`contact.add`**（传 `toUserId`、`toUsername`、`toAccount` 之一）

    ```json
    {
      "type": "contact.add",
      "tid": "ca-1",
      "payload": {
        "toUserId": "1234567890",
        "toUsername": "alice",
        "toAccount": "1234567890"
      }
    }
    ```

  - **成功响应：`contact.add.ok`**（payload 可为空）

  - **请求：`contact.remove`**

    ```json
    {
      "type": "contact.remove",
      "tid": "cr-1",
      "payload": { "contactUserId": "1234567890" }
    }
    ```

  - **成功响应：`contact.remove.ok`**

- 后续可扩展：`contact.request` / `contact.accept`（好友申请/通过流程）。

---

### 7. 限流与错误处理约定

- 当触发限流时，服务端返回：

```json
{
  "type": "message.send.error",
  "tid": "msg-1",
  "payload": null,
  "error": {
    "code": "rate_limited",
    "message": "too many requests"
  }
}
```

- 客户端建议行为：
  - 在发送按钮上显示节流提示；
  - 对于持续触发的用户，可考虑 UI 级别冷却倒计时。

---

---

### 8. 与后端服务的对应关系（参考）

> 本节仅说明协议与后端服务的责任划分，具体 gRPC 接口见 `rpc-auth-presence-message-conversation.md`。

- `auth.*` 消息 → Gateway 调用 **AuthService**
- `presence.*` 消息 → Gateway 调用 **PresenceService**
- `user.me` → Gateway 调用 **UserService**
- `message.send` / `message.history` / `message.read` → Gateway 调用 **MessageService**
- `conversation.*` 消息 → Gateway 调用 **ConversationService**

