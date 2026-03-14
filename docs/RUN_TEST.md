# User / Auth 开发进度与 Gateway 接入情况（敏捷运行测试）

## 一、开发进度概览

### 1. Auth 服务（`services/auth`）

| 能力 | 状态 | 说明 |
|------|------|------|
| **Proto 与骨架** | ✅ | `proto/auth.proto` 已定义，zrpc 已生成 |
| **Login** | ✅ | 用户名密码校验、bcrypt、双 Token 写入 Redis、返回 userId/accessToken/refreshToken/expiresIn |
| **TokenLogin** | ✅ | 凭 accessToken 校验并续期，返回同 LoginResponse |
| **ValidateToken** | ✅ | 校验 token 有效性，返回 valid + userId |
| **Logout** | ✅ | 按 accessToken 删除 Redis 会话 |
| **GetUserRoles** | ✅ | 查系统级角色名列表（依赖 RBAC 表） |
| **CheckPermission** | ✅ | 按用户 + 权限码校验（RBAC） |
| **AssignRoles** | ✅ | 覆盖式设置用户角色（RBAC） |
| **数据与依赖** | ✅ | PostgreSQL（users + RBAC 表）、Redis（token 会话）、GORM + go-redis |

- **配置**：`etc/beehive.auth.yaml`（ListenOn、Etcd 注册 Key: `beehive.auth.rpc`、PostgresDSN、Redis、Token TTL）。
- **与 Gateway**：Gateway 通过 etcd 发现 Auth 并调用 `AuthSvc.Login` / `TokenLogin` / `Logout`，已接入。

### 2. User 服务（`services/user`）

| 能力 | 状态 | 说明 |
|------|------|------|
| **Proto 与骨架** | ✅ | `proto/user.proto` 已定义，zrpc 已生成 |
| **GetUser** | ✅ | Redis 缓存优先，miss 查 PostgreSQL 并回写 |
| **BatchGetUsers** | ✅ | Redis MGet + DB 批查 + Pipeline 回填，按请求 ID 顺序返回 |
| **UpdateUser** | ✅ | 更新 user_profiles，并刷新 Redis 缓存 |
| **数据与依赖** | ✅ | PostgreSQL（users + user_profiles）、Redis（profile 缓存）、GORM + go-redis |

- **配置**：`etc/beehive.user.yaml`（ListenOn、Etcd 注册 Key: `beehive.user.rpc`、PostgresDSN、Redis、UserProfileTTLSeconds）。
- **与 Gateway**：Gateway 通过 etcd 发现 User 后调用 `GetUser`；客户端发送 `user.me` 可获取当前用户资料（昵称、头像、简介等）。

### 3. Gateway 与 Auth/Presence 接入情况

| 接入项 | 状态 | 说明 |
|--------|------|------|
| **AuthService** | ✅ | `ServiceContext` 注入 `AuthSvc`，配置 `AuthRpcConf`（Etcd 发现） |
| **PresenceService** | ✅ | 注入 `PresenceSvc`，配置 `PresenceRpcConf`（Etcd 发现） |
| **UserService** | ✅ | 注入 `UserSvc`，配置 `UserRpcConf`（Etcd 发现）；支持 `user.me` 获取当前用户资料 |
| **WebSocket 消息** | ✅ | `auth.login` / `auth.tokenLogin` → Auth RPC → 成功后 `RegisterSession` → 返回 `auth.login.ok` / `auth.tokenLogin.ok` |
| **auth.logout** | ✅ | 调 Auth.Logout、Presence.UnregisterSession，返回 `auth.logout.ok` |
| **presence.ping** | ✅ | 调用 Presence.RefreshSession 续期会话后返回 `presence.ping.ok` |
| **user.me** | ✅ | 调用 UserService.GetUser 返回当前用户资料（id、nickname、avatarUrl、bio、status） |
| **ConversationService** | ✅ | 可选；配置 ConversationRpcConf 后支持 `conversation.list` |
| **MessageService** | ✅ | 可选；配置 MessageRpcConf 后支持 `message.send`、`message.history` |
| **conversation.list** | ✅ | 调用 ConversationService.ListUserConversations + MessageService.GetLastMessages 聚合返回 |
| **message.send** | ✅ | 调用 MessageService.PostMessage（from_user_id 为当前连接用户） |
| **message.history** | ✅ | 调用 MessageService.GetHistory 返回历史消息 |

- **配置**：`etc/gateway-api.yaml`（Host、Port、GatewayID、AuthRpcConf、PresenceRpcConf、UserRpcConf、ConversationRpcConf、MessageRpcConf；**通过 etcd 发现** 各 RPC 服务）。
- **文档**：`docs/backend/gateway-design.md`、`docs/API/websocket-client-api.md`。

### 4. Conversation 服务（`services/conversation`）

| 能力 | 状态 | 说明 |
|------|------|------|
| **Proto 与骨架** | ✅ | `proto/conversation.proto` 已定义，zrpc 已生成 |
| **CreateConversation** | ✅ | 创建会话并写入 conversations + conversation_members |
| **AddMember / RemoveMember** | ✅ | 维护会话成员（RemoveMember 置 status=left） |
| **ListUserConversations** | ✅ | 按 user_id 分页返回会话列表（cursor/limit） |
| **GetConversation** | ✅ | 按 id 返回会话信息与成员数 |
| **ListMembers** | ✅ | 按 conversation_id 返回成员列表 |
| **数据与依赖** | ✅ | PostgreSQL（conversations、conversation_members）、GORM |

- **配置**：`etc/beehive.conversation.yaml`（ListenOn、Etcd Key: `beehive.conversation.rpc`、PostgresDSN）。
- **与 Gateway**：Gateway 可选配置 ConversationRpcConf；支持 WebSocket `conversation.list`（聚合 MessageService.GetLastMessages 返回 lastMessage）。

### 5. Message 服务（`services/message`）

| 能力 | 状态 | 说明 |
|------|------|------|
| **Proto 与骨架** | ✅ | `proto/message.proto` 已定义，zrpc 已生成 |
| **PostMessage** | ✅ | 写 messages 表，成功后可选发布 RabbitMQ `message.created` 事件 |
| **GetHistory** | ✅ | 按 conversation_id、before_time、limit 分页拉取历史 |
| **GetLastMessages** | ✅ | 按多 conversation_ids 返回各会话最后一条消息 |
| **数据与依赖** | ✅ | PostgreSQL（messages）、可选 RabbitMQ（im.events） |

- **配置**：`etc/beehive.message.yaml`（ListenOn、Etcd Key: `beehive.message.rpc`、PostgresDSN；可选 RabbitMQURL/Exchange/RouteKey）。
- **与 Gateway**：Gateway 可选配置 MessageRpcConf；支持 `message.send`、`message.history`。

### 6. Presence 服务（`services/presence`）

| 能力 | 状态 | 说明 |
|------|------|------|
| **RegisterSession** | ✅ | 将会话写入 Redis（user 索引 set + session Hash），设置 TTL |
| **UnregisterSession** | ✅ | 从 Redis 删除会话与索引 |
| **RefreshSession** | ✅ | 更新 last_ping_at 并续期 TTL |
| **GetOnlineSessions** | ✅ | 按用户查询在线会话列表，读时清理过期索引 |
| **GetUserPresence** | ✅ | 同上，并返回 online 与 sessions |
| **数据与依赖** | ✅ | Redis（`presence:user:{userId}:conns`、`presence:session:{userId}:{connId}`） |

- **配置**：`etc/beehive.presence.yaml`（ListenOn、Etcd 注册 Key: `beehive.presence.rpc`、Redis、SessionTTLSeconds）。
- **与 Gateway**：登录成功后 Gateway 调 RegisterSession；心跳时调 RefreshSession；登出时调 UnregisterSession。

---

## 二、本地运行测试（敏捷验证）

### 前置条件

1. **PostgreSQL**：已创建库 `beehive`，并执行迁移：
   - `db/migrations/001_create_users_and_user_profiles.sql`
   - `db/migrations/002_create_rbac_tables.sql`
   - `db/migrations/003_seed_test_user.sql`（插入接口测试用用户 testuser / password123）
   - `db/migrations/004_create_conversations_and_members.sql`
   - `db/migrations/005_create_messages.sql`
   - 使用 Docker 时：`docker compose -f docker/docker-compose.yaml up -d postgres`，然后依次执行上述 SQL 文件（如 `docker exec -i beehive-postgres psql -U beehive -d beehive < db/migrations/001_create_users_and_user_profiles.sql` 等）。
   - 或本机 `psql -U beehive -d beehive -f db/migrations/001_...sql` 等。
2. **Redis**：本地 6379 可用（无密码）。Docker：`docker compose -f docker/docker-compose.yaml up -d redis`。
3. **etcd**：用于服务发现。Docker：`docker compose -f docker/docker-compose.yaml up -d etcd`。各服务配置中已使用 `127.0.0.1:2379` 与 Key（如 `beehive.auth.rpc`、`beehive.presence.rpc`）。
4. **可选**：至少插入一名测试用户与角色（见下方「测试用户与角色」）。

### 启动顺序（每个终端在对应目录下执行，或使用 `-f` 指定配置绝对路径）

**先确保 etcd 已启动**（RPC 服务会向 etcd 注册，Gateway 从 etcd 发现）：

```bash
# 0. 基础设施（若用 Docker，在项目根目录）
docker compose -f docker/docker-compose.yaml up -d postgres redis etcd
```

```bash
# 1. Auth（依赖 PostgreSQL + Redis，启动后注册到 etcd Key: beehive.auth.rpc）
cd services/auth && go run . -f etc/beehive.auth.yaml

# 2. Presence（依赖 Redis 存储会话，注册到 etcd Key: beehive.presence.rpc）
cd services/presence && go run . -f etc/beehive.presence.yaml

# 3. Gateway（通过 etcd 发现 Auth / Presence / User）
cd services/gateway && go run . -f etc/gateway-api.yaml
```

若需使用 **user.me**（获取当前用户资料），须先启动 User 服务：

```bash
# 4. User
cd services/user && go run . -f etc/beehive.user.yaml
```

若需使用 **conversation.list**、**message.send**、**message.history**，须先启动 Conversation 与 Message 服务（并已执行 004、005 迁移）：

```bash
# 5. Conversation（依赖 PostgreSQL，注册到 etcd Key: beehive.conversation.rpc）
cd services/conversation && go run . -f etc/beehive.conversation.yaml

# 6. Message（依赖 PostgreSQL，可选 RabbitMQ；注册到 etcd Key: beehive.message.rpc）
cd services/message && go run . -f etc/beehive.message.yaml
```

### 测试用户与角色（可选）

Auth 的 Login 依赖 `users` 表中有用户，且 `GetUserRoles` 会查 RBAC 表（无角色时返回空列表，不影响登录）。若要完整验证 RBAC，可执行：

```sql
-- 插入默认角色与测试用户（密码需为 bcrypt 哈希，此处仅为示例）
-- 密码 'password123' 的 bcrypt 哈希示例（请自行生成替换）：
INSERT INTO roles (id, name, description) VALUES
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'user', 'normal user');
INSERT INTO user_roles (user_id, role_id) VALUES
  ('<你的用户UUID>', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11');
```

生成 bcrypt 哈希示例（Go）：`bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)`。

### 快速 WebSocket 登录验证

1. 用任意 WebSocket 客户端连接 `ws://127.0.0.1:8080/ws`。
2. 发送 JSON：
   ```json
   { "type": "auth.login", "tid": "1", "payload": { "username": "testuser", "password": "password123", "deviceId": "dev-1" } }
   ```
3. 若用户存在且密码正确，应收到 `type: "auth.login.ok"`，payload 含 `userId`、`accessToken`、`refreshToken`、`expiresIn`。
4. 再发 `auth.tokenLogin` 或 `presence.ping` 可验证后续流程。

### 会话与消息 WebSocket 验证（需已启动 Conversation、Message 并已登录）

1. **拉取会话列表** `conversation.list`：
   ```json
   { "type": "conversation.list", "tid": "conv-1", "payload": { "cursor": "", "limit": 50 } }
   ```
   成功响应：`type: "conversation.list.ok"`，payload 含 `items`（id、name、type、lastMessage、lastActiveAt 等）、`nextCursor`。

2. **发送消息** `message.send`（需先通过 Conversation 服务或直接插入获得有效 conversationId）：
   ```json
   { "type": "message.send", "tid": "msg-1", "payload": { "clientMsgId": "local-1", "conversationId": "<会话UUID>", "toUserId": "", "body": { "type": "text", "text": "hello" } } }
   ```
   成功响应：`type: "message.send.ok"`，payload 含 `serverMsgId`、`serverTime`、`conversationId`。

3. **拉取历史消息** `message.history`：
   ```json
   { "type": "message.history", "tid": "hist-1", "payload": { "conversationId": "<会话UUID>", "before": 0, "limit": 50 } }
   ```
   成功响应：`type: "message.history.ok"`，payload 含 `items`、`hasMore`。

---

## 三、小结

- **Auth**：登录/登出/Token 校验/RBAC 已实现，与 Gateway 通过 zrpc 对接完成。
- **User**：GetUser/BatchGetUsers/UpdateUser 已实现，Gateway 支持 `user.me`。
- **Conversation**：CreateConversation、AddMember、RemoveMember、ListUserConversations、GetConversation、ListMembers 已实现；Gateway 支持 `conversation.list`（聚合最后一条消息）。
- **Message**：PostMessage（写库 + 可选 RabbitMQ 事件）、GetHistory、GetLastMessages 已实现；Gateway 支持 `message.send`、`message.history`。
- **Gateway**：已接入 Auth、Presence、User、Conversation、Message；按配置通过 etcd 发现各 RPC；未配置的 Conversation/Message 时对应 WebSocket 类型返回 unavailable。

按上述顺序启动各服务，并执行全部 DB 迁移与测试用户，即可进行敏捷运行测试（含登录、会话列表、发消息、历史消息）。
