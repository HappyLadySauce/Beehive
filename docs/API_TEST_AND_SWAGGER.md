# 接口测试与 API 文档（go-zero）

本文说明 Beehive 的接口测试方式，以及如何用 go-zero 的 goctl 生成并查看 API 文档（Swagger）。

---

## 一、接口测试方式

### 1. Admin 管理后台 HTTP API

- **服务**：`services/adminapi`，默认端口 `8888`，基础路径 `/admin`。
- **认证**：除 `GET /admin/healthz` 外，请求头需带 `Authorization: Bearer <accessToken>`（token 由 Auth 登录或 Gateway WebSocket `auth.login` 获取）。

**示例（curl）**：

```bash
# 健康检查（无需认证）
curl -s http://127.0.0.1:8888/admin/healthz

# 需要认证的接口（先取 token，再请求）
curl -s -H "Authorization: Bearer <your-access-token>" http://127.0.0.1:8888/admin/users/某个用户ID
```

也可使用 **Postman / Insomnia**：新建请求，URL 填 `http://127.0.0.1:8888/admin/...`，在 Headers 中加 `Authorization: Bearer <token>`。

### 2. Gateway WebSocket（/ws）

- **入口**：`ws://127.0.0.1:8080/ws`
- **协议**：统一 JSON Envelope，见 [docs/API/websocket-client-api.md](API/websocket-client-api.md)。

**测试方式**：

- **wscat**（命令行）：`npm i -g wscat` 后执行 `wscat -c ws://127.0.0.1:8080/ws`，再在终端里手输 JSON 消息（如 `{"type":"auth.login","tid":"1","payload":{"username":"alice","password":"xxx","deviceId":"dev1"}}`）。
- **浏览器**：用前端项目或任意 WebSocket 客户端页面连 `ws://127.0.0.1:8080/ws`，按文档发送 `type`/`tid`/`payload`。
- **Postman**：新建 WebSocket 请求，URL 填 `ws://127.0.0.1:8080/ws`，连接后发送 JSON。

### 3. RPC 服务（Auth / User / Presence 等）

- 内部 gRPC，一般不直接对 HTTP 暴露。可用 **grpcurl** 或 **BloomRPC** 等工具连对应端口做调试；或通过 Gateway / AdminAPI 的 HTTP 或 WebSocket 间接调用。

---

## 二、go-zero 生成 API 文档（Swagger）

go-zero 支持从 **`.api` 文件** 生成 **Swagger 2.0** 文档，便于在线浏览与调试 HTTP 接口。

### 1. 命令（需 goctl >= 1.8.2）

```bash
# 为 Admin API 生成 Swagger JSON（输出到 docs/swagger/，生成文件为 admin.json）
goctl api swagger --api api/admin.api --dir docs/swagger

# 生成 YAML 格式（可选）
goctl api swagger --api api/admin.api --dir docs/swagger --yaml
```

- `--api`：指定 `.api` 定义文件。
- `--dir`：生成文件的输出目录（如 `docs/swagger`）。
- `--filename`：可选，生成的文件名（不含扩展名）。
- `--yaml`：生成 YAML，不写则默认 JSON。

### 2. .api 中的 info 配置

在 `api/admin.api` 的 `info` 块中可配置 Swagger 展示用信息，例如：

- `title`、`desc`/`description`、`version`：标题、描述、版本。
- `host`、`basePath`：如 `127.0.0.1:8888`、`/admin`，方便在 Swagger UI 里直接“Try it out”。
- `schemes`、`consumes`、`produces`：协议与内容类型（见 [go-zero Swagger 参考](https://go-zero.dev/zh-cn/reference/cli-guide/swagger/)）。

当前 `api/admin.api` 已包含 `title`、`desc`、`version`、`host`、`basePath`，可直接用于生成。

### 3. 查看 Swagger 文档

- **Swagger Editor**：打开 [https://editor.swagger.io](https://editor.swagger.io)，将生成的 `swagger.json` 内容粘贴进去即可浏览和调试。
- **本地 Swagger UI**：用 Docker 快速起一个 UI 服务指向生成的 `admin.json`：
  ```bash
  docker run -d -p 8081:8080 -e SWAGGER_JSON=/foo/admin.json -v ${PWD}/docs/swagger:/foo swaggerapi/swagger-ui
  ```
  浏览器访问 `http://127.0.0.1:8081` 即可（PowerShell 下 `${PWD}` 为当前目录；Linux/macOS 可用 `$(pwd)`）。

---

## 三、小结

| 接口类型           | 测试方式                     | 文档来源 |
|--------------------|------------------------------|----------|
| Admin HTTP API     | curl / Postman / Swagger UI  | goctl api swagger 从 `api/admin.api` 生成 |
| Gateway WebSocket  | wscat / 浏览器 / Postman WS | [websocket-client-api.md](API/websocket-client-api.md) |
| 内部 RPC           | grpcurl / BloomRPC 等       | proto 与 [rpc-auth-presence-message-conversation.md](API/rpc-auth-presence-message-conversation.md) |

生成 Swagger 后，可在 Swagger UI 中对 Admin API 做“Try it out”接口测试；Bearer Token 可在 Security 中配置为 `Authorization: Bearer <token>`。
