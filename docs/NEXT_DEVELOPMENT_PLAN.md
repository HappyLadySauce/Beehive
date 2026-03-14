# 下一步开发计划

在 Admin 权限中间件已完成的前提下，建议按以下顺序推进。

---

## 一、优先：测试与文档补齐（推荐先做）

**目标**：任何人按文档即可本地跑全栈并验证 Admin 认证/权限。

### 1.1 在 RUN_TEST.md 中补充 Admin 服务

- **启动顺序**：在「二、本地运行测试」中增加 Admin 服务启动步骤（依赖 Auth/User/Presence/Message/Conversation 已启动，通过 etcd 发现）。
  - 示例：`cd services/adminapi && go run . -f etc/admin-api.yaml`，默认 `0.0.0.0:8888`。
- **获取 Token**：说明先用 WebSocket `auth.login` 或 Auth 的登录接口拿到 `accessToken`，供 Admin 请求头使用。
- **验证方式**：
  - 无 Token 请求 `GET /admin/users` → 期望 `code: 1001`。
  - 有 Token 但用户无 `admin.user.read` 权限 → 期望 `code: 1003`。
  - 有 Token 且具备对应权限 → 进入 handler（当前 ListUsers 仍返回空列表，属预期）。
- **可选**：补充 Admin 健康检查 `GET /admin/healthz` 无需认证，用于负载均衡/就绪探针。

**产出**：更新 [docs/RUN_TEST.md](RUN_TEST.md)，新增「Admin 服务」小节。

### 1.2 部署文档（可选）

- 若需一键本地全量启动：在现有 `docker/docker-compose.yaml` 基础上，增加各 RPC 与 Gateway、Admin 的 service 定义与启动顺序说明，或单独写 `docs/deployment.md` 简述本地 Docker 全栈与端口。
- 本步可后置，与「占位接口对接」并行。

---

## 二、Admin 占位接口对接（按需排期）

以下四项相互独立，可分批实现；每项均依赖对应 RPC 或存储扩展。

| 接口 | 当前状态 | 依赖与方案概要 |
|------|----------|----------------|
| **ListUsers** | 返回空列表 | User 无 ListUsers RPC。方案：在 UserService 新增 `ListUsers(req)`（q/status/online/page/pageSize），查 users/user_profiles 表（可加 status 字段或沿用现有）；Admin ListUsersLogic 调该 RPC 并映射为 Admin API 响应。 |
| **ListConversations** | 返回空列表 | Conversation 仅有 ListUserConversations(user_id)。方案：ConversationService 新增 `ListConversations(type, memberId, page, pageSize)`（按 type/memberId 过滤、分页）；Admin ListConversationsLogic 调该 RPC。 |
| **Ban/Unban** | 仅返回 ok | 需持久化封禁状态并在登录/ValidateToken 时拒绝。方案：Auth 或 User 增加 Ban/Unban 存储与 RPC；ValidateToken 或 Login 时检查封禁状态返回错误；Admin BanUser/UnbanUser 调 RPC。 |
| **ListConfig / PutConfig** | 返回空/ok | 文档约定 etcd。方案：adminapi 内接入 etcd 客户端（或独立 config 服务），按 namespace/key 读写；ListConfigLogic/PutConfigLogic 调 etcd 或该服务。 |

**建议**：先做 **ListUsers** 或 **ListConversations** 其一（与前端用户列表/会话列表页面对齐），再做 Ban/Unban（需先定封禁存储与校验点），Config 可最后或与运维需求一起做。

---

## 三、可选：已读事件发布

- **目标**：Message 服务在执行 MarkRead 时，除写 `conversation_read` 外，可选向 RabbitMQ 发布「已读」事件（如 `message.read`），供其他端（如多端同步、运营统计）消费。
- **范围**：MessageService 内 MarkRead 成功后发布事件；Gateway 或其它消费者按需订阅。协议与 topic 需与现有 `message.created` 风格统一（见 [docs/backend/storage-and-mq.md](backend/storage-and-mq.md)）。
- **优先级**：低于「测试与文档」和「Admin 占位接口」；有明确多端已读同步或统计需求时再做。

---

## 四、建议执行顺序

1. **本周**：完成 **1.1 在 RUN_TEST.md 中补充 Admin 服务**，使 Admin 与权限校验可被任何人复现。
2. **随后**：从 Admin 占位接口中选一项（推荐 ListUsers 或 ListConversations）做 RPC 扩展 + Admin 对接。
3. **再后**：Ban/Unban 与 Config（etcd）按产品/运维需求排期；已读事件发布按需实现。

实现前请继续对照：

- Admin API：[docs/API/admin-http-api.md](API/admin-http-api.md)
- 服务设计：[docs/backend/services-design.md](backend/services-design.md)
- 开发流程：[docs/backend/development-workflow.md](backend/development-workflow.md)
