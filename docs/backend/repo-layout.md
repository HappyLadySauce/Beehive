## 仓库目录结构与公共组件规划

本节规划 Beehive 新版本的代码目录结构，便于后续使用 go-zero/goctl 生成各服务骨架代码并填充实现。

---

### 1. 顶层目录结构（建议）

```text
Beehive/
  docs/                      # 架构与设计文档
  proto/                     # 所有 gRPC 协议定义
    auth.proto
    user.proto
    presence.proto
    message.proto
    conversation.proto
    delivery.proto           # 如需单独的消息投递服务

  services/
    gateway/                 # WebSocket 网关（仅 WS/HTTP）
      cmd/gateway/           # main.go / 启动入口
      internal/              # go-zero 风格内部代码
        config/
        handler/             # WS 消息 handler，按 type 分发
        router/
        ws/                  # 连接/Session 抽象
        client/              # 调用下游服务的 RPC 客户端封装

    auth/
      cmd/auth/              # AuthService main.go
      internal/
        config/
        logic/
        svc/
        model/               # 用户/认证相关 PostgreSQL 表模型

    user/
      cmd/user/
      internal/
        config/
        logic/
        svc/
        model/               # 用户 profile 模型

    presence/
      cmd/presence/
      internal/
        config/
        logic/
        svc/
        store/               # Redis 读写封装

    message/
      cmd/message/
      internal/
        config/
        logic/
        svc/
        model/               # 消息&会话 PostgreSQL 模型
        mq/                  # RabbitMQ 生产端封装

    conversation/
      cmd/conversation/
      internal/
        config/
        logic/
        svc/
        model/               # 会话/成员关系模型

    adminapi/
      cmd/adminapi/
      internal/
        config/
        handler/             # go-zero REST handler
        logic/
        svc/

    delivery/                # 可选：专门消费 MQ 做消息投递
      cmd/delivery/
      internal/
        config/
        logic/
        svc/
        mq/                  # RabbitMQ 消费端封装

  pkg/                       # 公共库
    auth/                    # JWT 封装、当前用户信息解析等
    logger/                  # 日志封装（基于 logx）
    tracing/                 # Tracing 统一封装
    config/                  # 公共配置结构与加载辅助
    errors/                  # 统一错误码与错误结构
```

---

### 2. Go module 策略

- 继续使用单一 Go module：`github.com/HappyLadySauce/Beehive`。
- 各服务通过内部包路径（如 `services/auth/internal/...`）相互引用公共定义。
- 所有 `.proto` 文件统一放在 `proto/` 目录；
  - 使用 `goctl rpc protoc` 为每个服务在 `services/<service>/` 下生成 zrpc 服务端骨架代码（`beehive.<service>.go` + `internal/*`）；
  - 使用 `protoc` 按 `go_package = "services/<service>/pb;xxxpb"` 生成该服务对应的 gRPC client / server stub，物理路径位于 `services/<service>/pb/proto/*.pb.go`，供 Gateway 和其他服务引用。

---

### 3. 公共组件原则

- 尽量让 `pkg/` 只包含真正通用、与业务弱相关的工具（例如 JWT、日志、Tracing 封装）。
- 具体业务逻辑（如消息模型、会话规则）应放在各自服务的 `internal/` 目录中，避免过度抽象导致耦合。
- 关于错误码、统一响应结构：
  - 在 `pkg/errors` 中定义通用错误码与错误结构。
  - Gateway 和 AdminAPI 可以使用统一的错误编码向外暴露。

