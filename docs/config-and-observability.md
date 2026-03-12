## 配置管理与可观测性设计

本节定义各服务的配置结构，以及日志、Tracing 和指标监控的统一方案。

---

### 1. 配置管理

**目标**

- 每个服务有独立配置文件，但结构风格保持一致。
- 支持不同环境（dev/test/prod）的配置覆盖。

**通用配置字段（示例）**

所有服务的配置结构都包含以下通用部分：

```yaml
Name: beehive-gateway
Mode: dev               # dev/test/prod

Log:
  ServiceName: beehive-gateway
  Mode: console         # console/file
  Encoding: json        # json/plain
  Level: info

Tracing:
  Endpoint: ""          # OTLP / Jaeger 等

Metrics:
  Prometheus:
    Enabled: true
    ListenOn: ":9100"
Etcd:
  Endpoints:
    - "http://etcd-1:2379"
    - "http://etcd-2:2379"
  DialTimeout: 3s
  Prefix: "/beehive/config"   # 各服务在该前缀下读取自己的动态配置
```

**各服务特有字段示例**

- GatewayService

```yaml
Gateway:
  ListenOn: ":8080"

  # 限流参数
  RateLimit:
    LoginPerMin: 30
    MsgPerSecPerUser: 20

  # 下游 RPC 依赖
  AuthRpc:
    Target: "auth.rpc:8081"
  PresenceRpc:
    Target: "presence.rpc:8082"
  MessageRpc:
    Target: "message.rpc:8083"
  ConversationRpc:
    Target: "conversation.rpc:8084"
```

- AuthService / MessageService 等还需要：
  - PostgreSQL 连接信息（DSN/连接池参数）。
  - Redis 连接信息。
  - RabbitMQ 连接信息（写 MQ 的服务）。
  - Etcd 连接信息（用于服务注册发现与动态配置）。

---

### 2. 日志（Logging）

- 使用 go-zero 的 `logx` 作为统一日志框架。
- 要求：
  - 每条日志都包含：`service`、`env`、`traceId`、`spanId`、`userId`（如有）等基础字段。
  - 在 Gateway 中，针对每个 WebSocket 消息记录：
    - `type`、`tid`、`userId`、`latency`、`status` 等。
- 日志分级与输出：
  - dev 环境默认输出到控制台，plain/text 方便调试。
  - prod 环境使用 json 格式，输出到文件或集中日志系统（如 ELK）。

---

### 3. Tracing（分布式链路追踪）

- 使用 OpenTelemetry 作为统一标准，通过 go-zero 提供的集成或自定义中间件注入 trace。
- 要求：
  - 从 Gateway 接收到一条 WebSocket 消息开始创建 trace/span。
  - 调用下游 gRPC 服务时，将 trace context 透传过去。
  - 下游服务在处理请求时继续该 trace，使登录、发消息等完整链路可观测。
- 后端可选：
  - Jaeger / Tempo / OTLP + Grafana。

---

### 4. Metrics（指标监控）

- 采用 Prometheus 作为主要指标采集系统，所有服务暴露 `/metrics` 接口。
- 关键指标示例：
  - Gateway：
    - 当前连接数、登录成功/失败次数。
    - 消息收发 QPS、丢弃/错误计数。
    - 限流触发次数。
  - AuthService：
    - 登录请求 QPS、失败率。
  - MessageService：
    - 消息写入 QPS、失败率。
    - 平均写库延迟、MQ 发布延迟。
  - PresenceService：
    - 在线用户数量、注册/注销频率。
  - RabbitMQ：
    - 队列堆积长度、消费延迟等（通过 MQ 自身监控获取）。

---

### 5. 健康检查与熔断

- 每个服务提供健康检查接口：
  - HTTP：如 `GET /healthz`。
  - gRPC：可使用标准的 health check 协议。
- 依赖检查：
  - 对 PostgreSQL/Redis/RabbitMQ 的连接可在健康检查中进行简要探测，报告为「degraded」或「unhealthy」。
- 在 Gateway 调用下游 RPC 时：
  - 配置合理的超时与重试策略。
  - 在异常情况下向客户端返回合适的错误码和提示。

