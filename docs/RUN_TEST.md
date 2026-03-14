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

- **配置**：`etc/gateway-api.yaml`（Host、Port、GatewayID、AuthRpcConf、PresenceRpcConf、UserRpcConf；**通过 etcd 发现** Auth/Presence/User）。
- **文档**：`docs/backend/gateway-design.md`、`docs/API/websocket-client-api.md`。

### 4. Presence 服务（`services/presence`）

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
   - 使用 Docker 时：`docker compose -f docker/docker-compose.yaml up -d postgres`，然后：
     - `docker exec -i beehive-postgres psql -U beehive -d beehive < db/migrations/001_create_users_and_user_profiles.sql`
     - `docker exec -i beehive-postgres psql -U beehive -d beehive < db/migrations/002_create_rbac_tables.sql`
     - `docker exec -i beehive-postgres psql -U beehive -d beehive < db/migrations/003_seed_test_user.sql`
   - 或本机 `psql -U beehive -d beehive -f db/migrations/001_...sql`、`002_...sql`、`003_...sql`。
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

若需使用 **user.me**（获取当前用户资料），须先启动 User 服务；如需同时验证 AdminAPI，可再起：

```bash
# 4. User
cd services/user && go run . -f etc/beehive.user.yaml
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

---

## 三、小结

- **Auth**：登录/登出/Token 校验/RBAC 已实现，与 Gateway 通过 zrpc 对接完成。
- **User**：GetUser/BatchGetUsers/UpdateUser 已实现，Gateway 未接入；AdminAPI 等已用 UserService。
- **Gateway**：已接入 Auth + Presence；登录流程（auth.login / auth.tokenLogin → RegisterSession → auth.xxx.ok）可跑通；Presence 的 Redis Session 实现仍在进行中，当前 RegisterSession/RefreshSession 为占位。

按上述顺序启动 Auth → Presence → Gateway，并准备好 DB/Redis 与测试用户，即可进行敏捷运行测试。
