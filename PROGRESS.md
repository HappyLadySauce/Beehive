## 当前里程碑

- **Milestone**：M2 - IM 客户端界面重设计与群管/好友流程
- **状态**：进行中
- **时间范围**：2026-03 ～ 2026-04
- **参考**：M1 已完成；M2 见 `docs/frontend/im-client-redesign.md`，执行计划见 TODO.yml 中 `frontend-client` 与相关 backend/gateway 项。

## 一、整体进展概览

- **网关**：
  - Gateway：`api/gateway.api` 已定义，`services/gateway` 由 goctl 生成；WebSocket 升级、Hub、读循环已实现；已集成 Auth、Presence、User、Conversation、Message（`auth.login`/`auth.tokenLogin`/`auth.logout`、`presence.ping`、`user.me`、`conversation.list`、`message.send`、`message.history`）；可选配置 RabbitMQ 消费后，向本实例在线连接推送 `message.push`。
- **后端核心 RPC**：
  - AuthService：登录/登出/Token 校验/CheckPermission/GetUserRoles 已实现（见 `auth-service`）。
  - PresenceService：Redis Session 注册/刷新/注销、GetOnlineSessions/GetUserPresence 已实现（见 `presence-service`）。
  - MessageService：PostMessage（写库 + 可选 RabbitMQ 事件）、GetHistory、GetLastMessages 已实现（见 `message-service`）。
  - ConversationService：CreateConversation、AddMember、RemoveMember、ListUserConversations、GetConversation、ListMembers 已实现（见 `conversation-service`）。
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
  - WsEntry Logic：`ServeConn` 读循环解析 JSON Envelope，按 `type` 分发；已实现 `presence.ping`、`auth.login`/`auth.tokenLogin`/`auth.logout`、`user.me`、`conversation.list`、`message.send`、`message.history` 及错误占位。
- **Gateway 与 Conversation/Message 集成**
  - 配置：`ConversationRpcConf`、`MessageRpcConf`（Etcd 发现），未配置时对应 WebSocket 类型返回 unavailable。
  - `conversation.list`：调用 ConversationService.ListUserConversations + MessageService.GetLastMessages 聚合成会话列表（含 lastMessage）。
  - `message.send`：调用 MessageService.PostMessage（from_user_id 为当前连接用户），返回 `message.send.ok`。
  - `message.history`：调用 MessageService.GetHistory，返回 `message.history.ok`。
- **消息实时投递 `message.push`**（`message-delivery-push`）
  - Gateway 可选配置 RabbitMQ 消费（RabbitMQURL、Exchange、Queue、RouteKey）；`internal/push/consumer.go` 消费 `message.created`，按 ConversationService.ListMembers 取会话成员，按 PresenceService.GetOnlineSessions 取在线会话，仅向本实例（session.GatewayId == 本机）连接通过 Hub 推送 `message.push`；多实例时每实例独立队列绑定同一 exchange，各推各连接。
  - 配置示例见 `etc/gateway-api.yaml` 注释；RUN_TEST.md 已补充启用说明。
- **数据库迁移 004/005**
  - `004_create_conversations_and_members.sql`：`conversations`、`conversation_members` 表及索引。
  - `005_create_messages.sql`：`messages` 表及 `(conversation_id, server_time)` 索引。
- **ConversationService 完整实现**（`conversation-service`）
  - 配置：PostgresDSN；ServiceContext 注入 GORM + ConversationModel。
  - 逻辑：CreateConversation（事务写会话与成员）、AddMember/RemoveMember、ListUserConversations（cursor/limit）、GetConversation、ListMembers；错误使用 gRPC status 返回。
  - 配置示例：`etc/beehive.conversation.yaml`。
- **MessageService 完整实现**（`message-service`）
  - 配置：PostgresDSN、可选 RabbitMQ（URL/Exchange/RouteKey）；ServiceContext 注入 GORM + MessageModel + 可选 MQ Publisher。
  - 逻辑：PostMessage（写库 + 可选发布 `message.created`）、GetHistory（before_time/limit）、GetLastMessages（多会话最后一条）；错误使用 gRPC status 返回。
  - 配置示例：`etc/beehive.message.yaml`。
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

## 三、待完善内容（Admin 放最后）

### 核心 IM 能力（优先）

1. **会话创建/成员管理的 WebSocket 暴露**（见 `conversation-ws-create-member`）  
   - **已完成**：Gateway 已增加 `conversation.create`、`conversation.addMember`、`conversation.removeMember` 的 type 分发与 RPC 调用，协议 6.2 已补充请求/响应 JSON。

2. **会话列表未读计数 `unreadCount`**  
   - **已完成**：Message 服务新增 `conversation_read` 表与 MarkRead/GetUnreadCounts RPC；Gateway 在 `conversation.list` 聚合时调用 GetUnreadCounts 填入各 item 的 `unreadCount`。

3. **已读回执 `message.read`**（见 `message-read-receipt`）  
   - **已完成**：Gateway 处理 `message.read` 并调用 MessageService.MarkRead，返回 `message.read.ok`；与未读计数共用 `conversation_read` 存储。

### 体验与稳定性（随后）

4. **单聊会话解析**  
   - **已完成**：ConversationService 新增 FindOrCreateSingleConversation RPC；`message.send` 在仅传 `toUserId` 时先查/建单聊再 PostMessage。

5. **限流**  
   - **已完成**：Gateway 配置 Redis 与 RateLimitMessageSendPerMinute，`handleMessageSend` 入口检查限流，超限返回 `message.send.error`（code: `rate_limited`）。

### Admin 相关（最后）

6. **Admin 按路由权限校验**（见 `TODO.yml` 中 `admin-permission-middleware`）：**已完成**。CheckPermission 中间件已实现，路由按 admin.user.read/ban、admin.conversation.read、admin.message.read、admin.config.read/write、admin.ops.use 分组；ListUsers/封禁/配置等占位接口可后续对接 RPC 或 etcd。

### M2 - IM 客户端重设计与群管/好友（进行中）

**分阶段执行计划**（详见 TODO.yml 执行顺序）：

- **阶段 0**：前端设计文档与执行计划（不写代码）  
  - 完善 `docs/frontend/im-client-redesign.md`：与参考图对应、操作与 WS 依赖表。  
  - 定好 PROGRESS.md / TODO.yml 为唯一执行计划入口。

- **阶段 1 - 后端**（优先）：  
  7. 会话列表 memberCount 透传（Gateway conversation.list.ok）。  
  8. 群公告全链路（DB 已有 003；Conversation 读/填 announcement；Gateway 透传）。  
  9. conversation.get / conversation.listMembers（Gateway WS，群详情与成员含 role）。  
  10. 好友申请与通过（contact_requests + User RPC + Gateway contact.request/requestList/accept/decline）。  
  11. 群申请与审批（可选）：group_join_requests + RPC + Gateway group.apply/joinRequestList/approve/decline。

- **阶段 2 - 前端**：  
  12. 左侧栏 Beehive+头像（收起仅保留头像）。  
  13. 消息/联系人/设置 导航与列表。  
  14. 群聊窗口标题(人数)+右侧群公告与成员。  
  15. 好友/群通知与退出群聊、删除好友 操作入口。

---

## 四、下一步计划

- **优先顺序**：
  1. ~~网关与五大 RPC 基础集成~~、~~消息投递 `message.push`~~（已完成）。
  2. ~~会话 WS 暴露~~、~~未读计数~~、~~已读回执 `message.read`~~、~~单聊会话解析~~、~~限流~~（已完成）。
  3. **M2**：**先**完成后端（memberCount、群公告、conversation.get/listMembers、好友申请、群申请可选），**再**按 im-client-redesign 做前端（左侧栏、消息/联系人/设置、群聊右侧栏、通知与操作）。执行时按 TODO.yml 中「设计文档 → 后端项 → 前端项」顺序勾选。
  4. ~~Admin 权限中间件~~（已完成）；占位接口对接可后续迭代。

- 实现前请对照：
  - Gateway：`docs/backend/gateway-design.md`、`docs/API/websocket-client-api.md`。
  - 各 RPC：`docs/backend/services-design.md`、`docs/backend/repo-layout.md`、`docs/API/rpc-auth-presence-message-conversation.md`。
  - 服务开发流程：`docs/backend/development-workflow.md`（先 proto/.api → 再生成 → 再实现）。
