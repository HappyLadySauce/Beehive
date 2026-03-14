## 当前里程碑

- **Milestone**：M1 - 网关服务最小可用（Gateway 先于后端 RPC 落地）
- **状态**：进行中
- **时间范围**：2026-03 ～ 2026-04

## 一、整体进展概览

- **网关**：
  - Gateway：`api/gateway.api` 已定义，`services/gateway` 由 goctl 生成；WebSocket 升级、Hub、读循环与 `auth.login`/`auth.tokenLogin`/`presence.ping` 已实现；Presence 集成（RegisterSession/RefreshSession/UnregisterSession）进行中（见 `TODO.yml` 中 `gateway-presence-integration`）。
- **后端核心 RPC**：
  - AuthService：登录/登出/Token 校验/CheckPermission/GetUserRoles 已实现（见 `auth-service`）。
  - PresenceService：zrpc 骨架已生成，在线 Session 与 Redis 接入进行中（见 `presence-service`）。
  - MessageService：zrpc 骨架已生成，消息持久化与 MQ 发布待实现（见 `message-service`）。
  - ConversationService：zrpc 骨架已生成，会话/成员管理待实现（见 `conversation-service`）。
  - UserService：用户 Profile 读写已实现（GetUser/BatchGetUsers/UpdateUser，见 `user-service`）。
- **管理后台 API**：
  - AdminAPIService：独立 HTTP 服务 `services/adminapi`，已按 `docs/backend/development-workflow.md` 先定义 `api/admin.api` 再 goctl 生成并完善；认证中间件（Bearer Token + AuthService.ValidateToken）、用户/会话/消息/配置/运维接口骨架已实现（见 `TODO.yml` 中 `admin-service`）。

## 二、已完成工作

- **Gateway API 与骨架**
  - 在 `api/gateway.api` 中定义了 `/ws`（WebSocket 入口）、`/healthz`（健康检查），与 `docs/backend/gateway-design.md`、`docs/API/websocket-client-api.md` 对齐。
  - 使用 `goctl api go -api api/gateway.api -dir services/gateway -style go_zero` 生成 Gateway 服务骨架，入口 `services/gateway/gateway.go`，handler 位于 `internal/handler/gateway/`（WsEntry、Health）。
- **Gateway WebSocket 升级与连接管理**（`gateway-ws-upgrade`）
  - `internal/ws`：`Envelope`/`ErrBody` 协议结构体，`Connection`（ConnID、UserID、DeviceID、GatewayID）与 `Hub` 连接管理（注册/注销/按 ConnID 查询）。
  - WsEntry Handler：使用 `gorilla/websocket` 完成 HTTP→WebSocket 升级，从 Hub 注册连接并 defer 注销与关闭，将连接交给 Logic 处理。
  - WsEntry Logic：`ServeConn` 读循环解析 JSON Envelope，按 `type` 分发；已实现 `presence.ping` → `presence.ping.ok`，`auth.login`/`auth.tokenLogin`/`auth.logout` 及未知 type 返回错误占位，为后续 auth/presence 集成预留。
- **Proto 与 zrpc 骨架**
  - `proto/auth.proto`, `proto/presence.proto`, `proto/message.proto`, `proto/conversation.proto`, `proto/user.proto` 已与 `docs/API/rpc-auth-presence-message-conversation.md` 对齐。
  - 使用 goctl 为以上服务生成了 zrpc 代码骨架：`services/auth`、`services/presence`、`services/message`、`services/conversation`、`services/user`。
 - **UserService 用户 Profile 读写**（`user-service`）
  - 数据层：在 `db/migrations/001_create_users_and_user_profiles.sql` 中创建 `users` / `user_profiles` 表；在 `internal/model` 中通过 GORM 定义 `User` / `UserProfile` 与 `UserProfileModel`（FindByID/FindByIDs/UpdateProfile）。
  - 配置与依赖：`etc/user.yaml` / `etc/beehive.user.yaml` 配置 `PostgresDSN`、`RedisAddr`、`UserProfileTTLSeconds`；`ServiceContext` 初始化 GORM 与 go-redis 客户端，并注入 `UserProfileMod`。
  - 逻辑层：`GetUser` 先读 Redis 缓存，miss 时查 PostgreSQL 并回写缓存；`BatchGetUsers` 使用 Redis MGet + PostgreSQL 批查 + Pipeline 回填缓存；`UpdateUser` 更新 `user_profiles` 并刷新对应缓存。
- **AdminAPIService 管理后台 HTTP 服务**（`admin-service`）
  - 按 `docs/backend/development-workflow.md`：先定义 `api/admin.api`（与 `docs/API/admin-http-api.md` 对齐），再 `goctl api go` 生成 `services/adminapi` 骨架。
  - 配置与依赖：`etc/admin-api.yaml` 配置 Auth/User/Presence/Message/Conversation 的 RPC（etcd）；`ServiceContext` 注入上述 5 个 RPC 客户端。
  - 认证：`internal/middleware/authmiddleware.go` 从 `Authorization: Bearer <token>` 取 token，调用 AuthService.ValidateToken，将 userId 写入 context；除 `GET /admin/healthz` 外均需认证。
  - 接口实现：GetUser/GetUserSessions/KickUser 调用 User/Presence RPC；GetConversation/ListMembers/ListConversationMessages 调用 Conversation/Message RPC；ListUsers/Ban/Unban/ListConversations/Config/Ops 为占位或空数据，后续可对接 ListUsers RPC、封禁状态、etcd 配置等。

## 三、正在进行 / 阻塞项

- **Gateway**（见 `TODO.yml`）：
  - `gateway-presence-integration`：登录后 RegisterSession、心跳 RefreshSession、断线 UnregisterSession（依赖 PresenceService 最小可用）。
- **后端核心 RPC**：
  - `presence-service`：基于 Redis 的 Session 注册/刷新/注销与 GetOnlineSessions 实现（进行中）。
  - `message-service`、`conversation-service`：待实现。
- **Admin 后续可选**：
  - 按路由注入 CheckPermission 中间件（如 `admin.user.read`、`admin.user.ban`）；ListUsers/封禁/配置等对接具体 RPC 或 etcd。

## 四、下一步计划

- **优先顺序**：
  1. ~~`gateway-ws-upgrade`~~、~~`gateway-auth-integration`~~、~~`auth-service`~~、~~`user-service`~~、~~`admin-service` 骨架与核心接口~~（已完成）。
  2. **当前重点**：完成 `gateway-presence-integration`（依赖 presence-service 最小可用）；并行推进 `presence-service`（Redis Session 管理）。
  3. 随后：`message-service`（持久化 + MQ）、`conversation-service`（会话/成员），再视需求完善 Admin 权限中间件与占位接口。

- 实现前请对照：
  - Gateway：`docs/backend/gateway-design.md`、`docs/API/websocket-client-api.md`。
  - 各 RPC：`docs/backend/services-design.md`、`docs/backend/repo-layout.md`、`docs/API/rpc-auth-presence-message-conversation.md`。
  - 服务开发流程：`docs/backend/development-workflow.md`（先 proto/.api → 再生成 → 再实现）。
