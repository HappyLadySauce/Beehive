## Admin HTTP API 设计（AdminAPIService）

本文件定义 Admin 管理后台使用的 HTTP API 接口，主要由 `AdminAPIService` 提供，前端实现参考 `docs/frontend/admin-app.md`。

路径示例统一以 `/admin` 为前缀，可根据实际部署调整（例如通过网关或反向代理统一前缀）。

---

### 1. 通用约定

- **协议**：HTTPS + JSON。
- **认证**：推荐使用基于 JWT 的 Bearer Token 或管理端专用的会话机制。
- **返回结构**（示例）：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

- 字段说明：
  - `code`：0 表示成功，非 0 为错误码（可与 `pkg/errors` 中的定义对齐）。
  - `message`：简短的错误或提示信息。
  - `data`：业务数据。

常见错误码（示意）：

- `0`：成功
- `1001`：未认证（未登录或 token 失效）
- `1003`：权限不足
- `2001`：参数错误
- `3001`：资源不存在
- `5000`：内部错误

---

### 2. 用户管理接口

#### 2.1 查询用户列表

- **方法与路径**
  - `GET /admin/users`

- **请求参数（Query）**

| 名称        | 类型     | 必填 | 说明                         |
| ----------- | -------- | ---- | ---------------------------- |
| `q`         | string   | 否   | 关键字（ID/昵称/邮箱模糊匹配） |
| `status`    | string   | 否   | 用户状态：`active` / `banned` 等 |
| `online`    | bool     | 否   | 是否在线                     |
| `page`      | int      | 否   | 页码，从 1 开始             |
| `pageSize`  | int      | 否   | 每页数量，默认 20           |

- **响应示例**

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "items": [
      {
        "id": "u_123",
        "nickname": "Alice",
        "email": "alice@example.com",
        "status": "active",
        "isOnline": true,
        "lastLoginAt": "2024-01-01T12:00:00Z"
      }
    ],
    "page": 1,
    "pageSize": 20,
    "total": 1
  }
}
```

#### 2.2 获取用户详情

- **方法与路径**
  - `GET /admin/users/{id}`

- **响应示例**

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": "u_123",
    "nickname": "Alice",
    "email": "alice@example.com",
    "status": "active",
    "createdAt": "2024-01-01T00:00:00Z",
    "lastLoginAt": "2024-01-01T12:00:00Z",
    "profile": {
      "avatarUrl": "",
      "bio": ""
    }
  }
}
```

#### 2.3 获取用户在线会话/设备

- **方法与路径**
  - `GET /admin/users/{id}/sessions`

- **响应示例**

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "sessions": [
      {
        "gatewayId": "gw-1",
        "connId": "c-123",
        "deviceId": "web-uuid",
        "deviceType": "web",
        "ip": "1.2.3.4",
        "loginAt": "2024-01-01T12:00:00Z",
        "lastPingAt": "2024-01-01T12:10:00Z"
      }
    ]
  }
}
```

#### 2.4 强制下线

- **方法与路径**
  - `POST /admin/users/{id}/kick`

- **请求体**

```json
{
  "reason": "manual_kick",
  "sessionIds": ["c-123"]   // 可选，空时表示该用户所有会话
}
```

- **响应**

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

#### 2.5 封禁 / 解封用户

- **方法与路径**
  - 封禁：`POST /admin/users/{id}/ban`
  - 解封：`POST /admin/users/{id}/unban`

- **请求体**

```json
{
  "reason": "abuse_spam",
  "until": "2024-02-01T00:00:00Z"  // 可选，空表示永久
}
```

- **响应**

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

---

### 3. 会话与消息相关接口

#### 3.1 查询会话列表

- **方法与路径**
  - `GET /admin/conversations`

- **请求参数（Query）**

| 名称        | 类型   | 必填 | 说明                     |
| ----------- | ------ | ---- | ------------------------ |
| `type`      | string | 否   | `single` / `group` 等   |
| `memberId`  | string | 否   | 包含指定用户的会话       |
| `page`      | int    | 否   | 页码                     |
| `pageSize`  | int    | 否   | 每页数量                 |

- **响应示例**

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "items": [
      {
        "id": "conv_abc",
        "type": "single",
        "name": "Alice & Bob",
        "memberCount": 2,
        "createdAt": "2024-01-01T00:00:00Z",
        "lastActiveAt": "2024-01-01T12:00:00Z"
      }
    ],
    "page": 1,
    "pageSize": 20,
    "total": 1
  }
}
```

#### 3.2 获取会话详情与成员

- **方法与路径**
  - 会话详情：`GET /admin/conversations/{id}`
  - 成员列表：`GET /admin/conversations/{id}/members`

- **响应示例（成员列表）**

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "items": [
      {
        "userId": "u_123",
        "role": "owner",     // owner | admin | member
        "joinedAt": "2024-01-01T00:00:00Z",
        "status": "active"
      }
    ]
  }
}
```

#### 3.3 查询会话消息

- **方法与路径**
  - `GET /admin/conversations/{id}/messages`

- **请求参数（Query）**

| 名称        | 类型   | 必填 | 说明                       |
| ----------- | ------ | ---- | -------------------------- |
| `fromUserId`| string | 否   | 发送方用户 ID             |
| `start`     | string | 否   | 起始时间（ISO8601）       |
| `end`       | string | 否   | 截止时间（ISO8601）       |
| `keyword`   | string | 否   | 文本关键字（仅文本消息）  |
| `page`      | int    | 否   | 页码                       |
| `pageSize`  | int    | 否   | 每页数量                   |

- **响应示例**

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "items": [
      {
        "serverMsgId": "msg_789",
        "conversationId": "conv_abc",
        "fromUserId": "u_123",
        "body": {
          "type": "text",
          "text": "hello"
        },
        "serverTime": "2024-01-01T12:00:00Z"
      }
    ],
    "page": 1,
    "pageSize": 50,
    "total": 1
  }
}
```

---

### 4. 配置管理接口（etcd 配置）

用于在后台 UI 中管理存储于 etcd 的运行时配置（feature flags、限流参数等）。

#### 4.1 查询配置列表

- **方法与路径**
  - `GET /admin/config`

- **请求参数（Query）**

| 名称        | 类型   | 必填 | 说明         |
| ----------- | ------ | ---- | ------------ |
| `namespace` | string | 否   | 配置命名空间 |

- **响应示例**

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "items": [
      {
        "key": "feature.enable_register",
        "value": "true",
        "description": "是否允许用户注册",
        "updatedBy": "admin",
        "updatedAt": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

#### 4.2 更新配置

- **方法与路径**
  - `PUT /admin/config/{key}`

- **请求体**

```json
{
  "value": "false",
  "description": "关闭注册以进行维护"
}
```

- **响应**

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

---

### 5. 运维工具相关接口（占位）

以下接口为未来预留，具体字段可在实现前细化：

- 队列状态查询：
  - `GET /admin/ops/queues`
- 消息重放（测试环境使用）：
  - `POST /admin/ops/replay`
- 服务健康状态汇总：
  - `GET /admin/ops/health`

