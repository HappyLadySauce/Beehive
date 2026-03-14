# 服务开发流程（proto/api → 生成代码 → 完善实现）

本节约定：**凡涉及新增或修改后端服务、接口时，必须严格按「先定义接口、再生成代码、最后完善实现」的顺序进行**，避免手写接口与生成代码脱节。

---

## 1. 原则：先定义、再生成、后实现

- **第一步：在 `proto/` 或 `.api` 中定义接口**  
  先修订 `proto/*.proto`（gRPC/zrpc）或 `*.api`（REST），再动业务代码。接口形态以 proto/api 为唯一来源。
- **第二步：用 goctl 等工具生成代码**  
  根据 proto/api 生成服务骨架、handler、request/response 类型，不手写与接口定义重复的结构。
- **第三步：在生成出的骨架中完善业务逻辑**  
  在 logic、model、中间件等位置补全实现，保持与 proto/api 定义一致。

这样可保证：文档 ↔ proto/api ↔ 生成代码 ↔ 实现 一致，后续改接口时只需改 proto/api 并重新生成，再调整实现即可。

---

## 2. 适用场景

以下情况均应按本流程执行：

- **新增或修改 RPC 服务**（Auth / User / Presence / Message / Conversation 等）  
  → 先改 `proto/*.proto`，再 `goctl rpc protoc` 生成，再完善 logic/model。
- **新增或修改 Admin HTTP 服务/路由**  
  → 先改 `docs/API/admin-http-api.md` 与对应 `.api` 文件，再 `goctl api go` 生成，再完善 handler/logic。
- **为已有服务新增/修改 RPC 方法或 HTTP 路由**  
  → 先更新 proto 或 .api，再重新生成相关代码，最后在生成出的 logic/handler 中完善实现。

仅做**纯业务逻辑修补**（不增删接口、不改请求/响应结构）时，可直接改 logic，无需重新走「定义 → 生成」步骤。

---

## 3. 具体步骤

### 3.1 gRPC / zrpc 服务（proto）

1. **修订 `docs/API/rpc-auth-presence-message-conversation.md`**（可选但推荐）  
   先在设计文档中写清方法签名、请求/响应字段，便于与 proto 对齐。
2. **修订 `proto/<service>.proto`**  
   在对应服务的 `.proto` 中新增/修改 `service` 与 `message`，保证方法名、字段与文档一致。
3. **生成代码**  
   使用 goctl 等工具生成 zrpc 服务端与 pb 客户端，例如：
   - `goctl rpc protoc proto/xxx.proto --go_out=... --go-grpc_out=... --zrpc_out=services/xxx`
4. **完善实现**  
   在 `services/<service>/internal/logic/` 下实现各 RPC 的 business logic，不擅自改生成出的 request/response 结构。

### 3.2 Admin HTTP API（.api + rest）

1. **修订 `docs/API/admin-http-api.md`**  
   先在设计文档中写清路由、Method、Query/Body、响应结构及错误码。
2. **修订或新增 `.api` 文件**  
   在 `api/` 或约定目录下编写/更新 Admin 的 `.api` 定义（路由、请求/响应类型），与 admin-http-api.md 一致。
3. **生成代码**  
   使用 goctl 生成 REST 骨架，例如：
   - `goctl api go -api admin.api -dir services/adminapi`
4. **完善实现**  
   在生成的 handler/logic 中注入 RPC 客户端、实现认证/鉴权中间件、补全业务逻辑；不手写与 .api 重复的路由或类型定义。

### 3.3 Gateway（WebSocket）

- WebSocket 消息类型与 Envelope 以 **`docs/API/websocket-client-api.md`** 为准。
- 若新增消息类型或字段，先更新该文档，再在 Gateway 的 handler/router 中实现解析与转发，必要时再调整下游 proto。

---

## 4. 记录与检查

- 在 Code Review 或任务说明中，若涉及接口变更，应标明：  
  「已按 development-workflow：先更新 proto/api → 已生成代码 → 在 xxx 中完成实现」。
- 若发现**未先更新 proto/api 就手写接口或 handler**，应回退为：先补充/修订 proto 或 .api，重新生成，再在生成代码上完善。

---

## 5. 相关文档

- 接口设计总览与步骤：`docs/API/README.md`、`.cursor/skills/beehive-api-design/SKILL.md`
- 服务职责与目录结构：`docs/backend/services-design.md`、`docs/backend/repo-layout.md`
