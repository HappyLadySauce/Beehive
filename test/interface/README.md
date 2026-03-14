# Beehive WebSocket 接口测试与模块文档

使用 Node.js（pnpm + ws + Vitest）模拟 WebSocket 客户端，对 Gateway `/ws` 进行自动化接口测试；并通过 AsyncAPI 规范按模块生成接口文档。

## 前置条件

- Node.js 18+
- pnpm（已安装则跳过：`npm install -g pnpm`）
- Gateway 及依赖服务已启动（见 [docs/RUN_TEST.md](../docs/RUN_TEST.md)），默认 `ws://127.0.0.1:8080/ws`
- **测试用户**：需执行数据库迁移以插入接口测试用用户。按 [docs/RUN_TEST.md](../docs/RUN_TEST.md) 执行 `001`、`002`、`003` 后，将存在用户 `testuser` / 密码 `password123`（见 [db/migrations/003_seed_test_user.sql](../db/migrations/003_seed_test_user.sql)）。未执行 003 时，依赖登录的用例会失败。

## 安装

```bash
pnpm install
```

## 运行测试

```bash
# 监听模式（开发时）
pnpm test

# 单次运行（CI）
pnpm test:run
```

可通过环境变量覆盖 WebSocket 地址与测试账号：

- `WS_URL`：默认 `ws://127.0.0.1:8080/ws`
- `TEST_USER` / `TEST_PASSWORD`：可选；默认使用种子用户 `testuser` / `password123`。若使用其他账号，需在库中已存在且密码正确。

## 生成按模块接口文档

从 [asyncapi.yaml](./asyncapi.yaml) 生成 `docs/API/generated/` 下按模块的 Markdown 文档：

```bash
pnpm run docs:generate
```

生成后可在 [docs/API/generated/README.md](../docs/API/generated/README.md) 查看索引。

## 目录说明

| 路径 | 说明 |
|------|------|
| `src/client.ts` | WebSocket 客户端封装（connect、sendAndWait、按 tid 配对） |
| `src/types.ts` | Envelope、ErrBody 及各 type 的 payload 类型 |
| `src/config.ts` | WS_URL 等配置 |
| `tests/*.test.ts` | 按域组织的 Vitest 用例（auth、presence、user） |
| `asyncapi.yaml` | WebSocket 协议 AsyncAPI 描述，按 tag 分模块 |
| `scripts/generate-docs.mjs` | 从 AsyncAPI 生成各模块文档的脚本 |

完整协议说明见 [docs/API/websocket-client-api.md](../docs/API/websocket-client-api.md)。
