## 当前里程碑

- **Milestone**：M1 - 打通后端核心 RPC 服务（Auth/Presence/Message/Conversation/User）
- **状态**：进行中
- **时间范围**：2026-03 ～ 2026-04

## 一、整体进展概览

- **服务实现进度**：
  - AuthService：zrpc 骨架已生成，登录/Token 校验逻辑未实现（见 `TODO.yml` 中任务 `auth-service`）。
  - PresenceService：zrpc 骨架已生成，在线 Session 模型与 Redis 接入待实现（见 `TODO.yml` 中任务 `presence-service`）。
  - MessageService：zrpc 骨架已生成，消息持久化与 MQ 发布逻辑待实现（见 `TODO.yml` 中任务 `message-service`）。
  - ConversationService：zrpc 骨架已生成，会话模型/列表与成员管理逻辑待实现（见 `TODO.yml` 中任务 `conversation-service`）。
  - UserService：zrpc 骨架已生成，用户 Profile 读写与数据库映射待实现（见 `TODO.yml` 中任务 `user-service`）。

## 二、已完成工作

- **Proto 与 zrpc 骨架**
  - `proto/auth.proto`, `proto/presence.proto`, `proto/message.proto`, `proto/conversation.proto`, `proto/user.proto` 已与 `docs/API/rpc-auth-presence-message-conversation.md` 中的设计对齐。
  - 使用 goctl 为以上服务生成了 zrpc 代码骨架，入口与目录如下（按服务划分）：
    - AuthService：`services/auth`（入口 `beehive.auth.go`）
    - PresenceService：`services/presence`（入口 `beehive.presence.go`）
    - MessageService：`services/message`（入口 `beehive.message.go`）
    - ConversationService：`services/conversation`（入口 `beehive.conversation.go`）
    - UserService：`services/user`（入口 `beehive.user.go`）

## 三、正在进行 / 阻塞项

- M1 期内关键任务（详情见 `TODO.yml`）：
  - `auth-service`：需要设计并实现基于文档的登录、登出、Token 校验逻辑，结合后续用户存储方案。
  - `presence-service`：需要完成 Redis 连接配置与在线 Session 存储读写，实现 Register/Unregister/Refresh 逻辑。
  - `message-service`：需要完成消息表设计、写入流程以及 `message.created` 事件发布到 RabbitMQ。
  - `conversation-service`：需要完成会话与成员关系表设计，打通创建会话、成员增删与用户会话列表。
  - `user-service`：需要完成用户基本资料表设计与 Get/BatchGet/Update 的最小实现。

## 四、下一步计划

- 优先顺序建议：
  1. 优先实现 `presence-service`，保证 Gateway 有可靠的在线状态查询能力。
  2. 随后实现 `auth-service`，打通登录与 Token 校验链路，为 WebSocket 与 Admin 提供鉴权基础。
  3. 在此基础上推进 `message-service` 与 `conversation-service`，打通最小聊天闭环。
  4. 最后完善 `user-service` 的 Profile 能力，为前端展示提供数据源。

- 每个服务在实现前，应先对照：
  - `docs/backend/services-design.md`：确认服务职责与边界。
  - `docs/backend/repo-layout.md`：确认目录结构与公共组件使用方式。
  - `docs/API/rpc-auth-presence-message-conversation.md`：确认 RPC 接口与字段语义。

