## Beehive API 总览

本目录用于描述 Beehive 新架构下对外/对内的各类接口设计，包括：

- WebSocket 客户端协议（IM Web 客户端通过 Gateway `/ws` 使用）
- Admin 管理后台使用的 HTTP API
- 内部各业务服务（Auth / User / Presence / Message / Conversation 等）的 gRPC / zrpc 接口

### 文档结构

- `websocket-client-api.md`  
  - 定义 WebSocket 的统一消息 Envelope 结构（`type` / `tid` / `payload` / `error`）  
  - 约定认证、心跳、消息发送/推送、会话列表等客户端可见的消息类型与字段
- `admin-http-api.md`  
  - 定义 Admin 管理后台使用的 HTTP API 路由、请求/响应模型与错误码
- `rpc-auth-presence-message-conversation.md`  
  - 归档内部核心服务的 gRPC 接口设计（将来会对应 `proto/*.proto`）

### 设计原则

- **统一 Envelope**：客户端与 Gateway 间所有 WebSocket 消息都使用统一的 JSON Envelope，便于前后端对齐与追踪。
- **服务边界清晰**：业务语义由后端服务负责，Gateway 只做路由和编排；文档中会标出每个接口属于哪个服务。
- **可落地到 .proto / .api**：接口设计应能较为直接地映射到 go-zero 的 `.proto` / `.api` 文件，方便后续使用 goctl 生成代码。

