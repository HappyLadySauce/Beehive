## 当前里程碑

- **Milestone**：M1 - 网关服务最小可用（Gateway 先于后端 RPC 落地）
- **状态**：进行中
- **时间范围**：2026-03 ～ 2026-04

## 一、整体进展概览

- **网关（优先）**：
  - Gateway：`api/gateway.api` 已定义，`services/gateway` 骨架已由 goctl 生成（`/ws`、`/healthz`）；WebSocket 升级、Auth/Presence 集成待实现（见 `TODO.yml` 中 `gateway-*` 任务）。
- **后端核心 RPC**：
  - AuthService：zrpc 骨架已生成，登录/Token 校验逻辑未实现（见 `TODO.yml` 中任务 `auth-service`）。
  - PresenceService：zrpc 骨架已生成，在线 Session 与 Redis 接入待实现（见 `TODO.yml` 中任务 `presence-service`）。
  - MessageService：zrpc 骨架已生成，消息持久化与 MQ 发布待实现（见 `TODO.yml` 中任务 `message-service`）。
  - ConversationService：zrpc 骨架已生成，会话/成员管理待实现（见 `TODO.yml` 中任务 `conversation-service`）。
  - UserService：zrpc 骨架已生成，用户 Profile 读写待实现（见 `TODO.yml` 中任务 `user-service`）。

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

## 三、正在进行 / 阻塞项

- **Gateway（当前重点，见 `TODO.yml`）**：
  - `gateway-auth-integration`：对接 AuthService，处理 `auth.login` / `auth.tokenLogin`。
  - `gateway-presence-integration`：对接 PresenceService，注册/刷新/注销 Session。
- **后端核心 RPC（在 Gateway 打通后再推进）**：
  - `auth-service`、`presence-service`、`message-service`、`conversation-service`、`user-service` 的业务逻辑实现。

## 四、下一步计划

- **优先顺序**：
  1. ~~实现 `gateway-ws-upgrade`~~（已完成：升级、Hub、读循环、按 type 分发、presence.ping 响应）。
  2. 实现 `gateway-auth-integration` 与 `gateway-presence-integration`（依赖 AuthService、PresenceService 至少具备最小可用实现）。
  3. 随后按需实现 `presence-service`、`auth-service`，再推进 `message-service`、`conversation-service`、`user-service`。

- 实现前请对照：
  - Gateway：`docs/backend/gateway-design.md`、`docs/API/websocket-client-api.md`。
  - 各 RPC 服务：`docs/backend/services-design.md`、`docs/backend/repo-layout.md`、`docs/API/rpc-auth-presence-message-conversation.md`。
