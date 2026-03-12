## 业务服务设计与职责边界

本节定义 Auth / User / Presence / Message / Conversation / Admin 等服务的职责与 gRPC 接口边界（使用 go-zero `zrpc` 实现），后续可以据此编写对应的 `.proto` 文件并通过 goctl 生成代码骨架。

---

### 1. AuthService

**职责**

- 用户登录/登出、注册（可选）。
- 访问令牌（accessToken）与刷新令牌（refreshToken）的签发与校验。
- 提供用户认证相关的统一入口给 Gateway 和其他内部服务。

**核心接口（示意）**

- `rpc Login(LoginRequest) returns (LoginResponse)`
- `rpc TokenLogin(TokenLoginRequest) returns (LoginResponse)`  
- `rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse)`
- `rpc Logout(LogoutRequest) returns (LogoutResponse)`

**数据存储**

- PostgreSQL：`users` 表（账号、密码哈希、状态等）。
- 可使用 Redis 做登录状态或黑名单缓存。

---

### 2. UserService

**职责**

- 用户基础资料管理（昵称、头像、简介、扩展字段）。
- 提供用户信息读接口给 Gateway/Message/Conversation 等服务。

**核心接口（示意）**

- `rpc GetUser(GetUserRequest) returns (GetUserResponse)`
- `rpc BatchGetUsers(BatchGetUsersRequest) returns (BatchGetUsersResponse)`
- `rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse)`

**数据存储**

- PostgreSQL：`user_profiles` 等表。
- Redis：用户资料缓存，减轻数据库压力。

---

### 3. PresenceService

**职责**

- 管理用户在线状态与多端 session。
- 记录每个用户当前在哪些 Gateway 实例、哪些连接上在线，以及设备信息。
- 为消息投递和在线状态展示提供数据支持。

**核心接口（示意）**

- `rpc RegisterSession(RegisterSessionRequest) returns (RegisterSessionResponse)`
- `rpc UnregisterSession(UnregisterSessionRequest) returns (UnregisterSessionResponse)`
- `rpc RefreshSession(RefreshSessionRequest) returns (RefreshSessionResponse)`  // 心跳续期
- `rpc GetOnlineSessions(GetOnlineSessionsRequest) returns (GetOnlineSessionsResponse)`
- `rpc GetUserPresence(GetUserPresenceRequest) returns (GetUserPresenceResponse)`

**数据存储**

- Redis：
  - `user:{userId}:sessions` → set/list of `{gatewayId, connId, device}`。
  - TTL 用于自动过期，结合心跳刷新。

---

### 4. MessageService

**职责**

- 处理点对点/群聊消息的写入与查询。
- 负责消息的持久化（至少文本/基础消息），为历史消息、审计和统计提供数据。
- 在写入成功后，向 RabbitMQ 发布消息事件（如 `message.created`）。

**核心接口（示意）**

- `rpc PostMessage(PostMessageRequest) returns (PostMessageResponse)`
- `rpc GetHistory(GetHistoryRequest) returns (GetHistoryResponse)`
- `rpc GetLastMessages(GetLastMessagesRequest) returns (GetLastMessagesResponse)` // 列表页摘要

消息结构建议包含：

- `serverMsgId`, `clientMsgId`
- `conversationId`
- `fromUserId`, `toUserId`（单聊）或 `toConversationId`（群聊）
- `body`（类型 + 内容，如 text/image/custom）
- `serverTime` 等。

**数据存储**

- PostgreSQL：
  - `messages` 表：按 `conversation_id + server_time` 排序。
  - 按需增加分表/索引。
- RabbitMQ：
  - 写库成功后发布 `message.created` 事件，供 Delivery/Notification/Analytics 等服务消费。

---

### 5. ConversationService

**职责**

- 管理会话/群组/频道的元数据和成员关系。
- 提供用户的会话列表、成员列表等能力。

**核心接口（示意）**

- `rpc CreateConversation(CreateConversationRequest) returns (CreateConversationResponse)`
- `rpc AddMember(AddMemberRequest) returns (AddMemberResponse)`
- `rpc RemoveMember(RemoveMemberRequest) returns (RemoveMemberResponse)`
- `rpc ListUserConversations(ListUserConversationsRequest) returns (ListUserConversationsResponse)`
- `rpc GetConversation(GetConversationRequest) returns (GetConversationResponse)`
- `rpc ListMembers(ListMembersRequest) returns (ListMembersResponse)`

**数据存储**

- PostgreSQL：
  - `conversations`：会话/群组基础信息。
  - `conversation_members`：成员关系，多对多。
- Redis：
  - 可缓存用户常用的会话列表、群成员列表等热点数据。

---

### 6. AdminAPIService

**职责**

- 运维/运营管理接口，仅对内部或管理后台开放。
- 查询用户状态、强制踢下线、封禁用户/会话、审计消息等。

**协议形式**

- 使用 go-zero `rest` 提供 HTTP API，前端管理后台可直接使用。

**典型接口（示意）**

- `GET /admin/users/{id}`
- `GET /admin/users/{id}/presence`
- `POST /admin/users/{id}/kick`
- `POST /admin/users/{id}/ban`
- `GET /admin/conversations/{id}`
- `GET /admin/conversations/{id}/messages`

---

### 7. 后续可扩展服务

- **NotificationService**
  - 消费 RabbitMQ 消息事件，实现移动推送、邮件、短信等。
- **AnalyticsService**
  - 消费消息与行为事件，做统计与报表。

