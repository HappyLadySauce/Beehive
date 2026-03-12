## Beehive 架构与文档总览

本文档目录用于描述全新的 Beehive IM 系统架构设计，代码可以后续通过 go-zero 代码生成工具拉起骨架再实现。

- `architecture-overview.md`：整体架构、服务划分与数据流。
- `gateway-design.md`：Gateway（WebSocket 网关）连接模型、消息协议与路由。
- `services-design.md`：Auth / User / Presence / Message / Conversation / Admin 等服务职责与接口边界。
- `storage-and-mq.md`：PostgreSQL / Redis / RabbitMQ 的职责划分与使用方式。
- `repo-layout.md`：仓库目录结构与公共组件规划。
- `config-and-observability.md`：配置管理、日志、Tracing、Metrics 设计。
- `migration-strategy.md`：从旧版 Beehive 到新版架构的迁移思路（如有需要）。

