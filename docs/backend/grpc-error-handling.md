# gRPC 服务错误返回规范

本规范约定 Beehive 各 RPC 服务（如 Auth、User、Presence、Message、Conversation 等）在 logic 层向调用方返回错误时的方式，便于客户端根据 gRPC 状态码区分错误类型并做统一处理。

---

## 1. 使用 status 包返回错误

- **不要**在对外 RPC 接口的 logic 中直接返回 `errors.New` 或 `fmt.Errorf`。
- **应当**使用 `google.golang.org/grpc/status` 与 `google.golang.org/grpc/codes` 构造并返回错误，使 gRPC 能正确传递状态码与消息。

示例：

```go
import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 参数错误
return nil, status.Error(codes.InvalidArgument, "user_id is required")

// 资源未找到
return nil, status.Error(codes.NotFound, "user not found")

// 未认证（如 token 无效、用户名密码错误）
return nil, status.Error(codes.Unauthenticated, "invalid username or password")

// 前置条件不满足（如账号被禁用）
return nil, status.Error(codes.FailedPrecondition, "account is not in normal status")

// 服务内部错误（数据库、Redis、下游 RPC 等）
return nil, status.Errorf(codes.Internal, "get user failed: %v", err)
```

---

## 2. 状态码使用约定

| 场景 | 推荐状态码 | 说明 |
|------|------------|------|
| 必填参数缺失或格式明显错误 | `InvalidArgument` | 客户端可提示用户修正输入 |
| 请求的资源不存在 | `NotFound` | 如用户 ID、会话 ID 查不到 |
| 未登录或凭证无效/过期 | `Unauthenticated` | 登录失败、token 无效等 |
| 已认证但状态不允许操作 | `FailedPrecondition` | 如账号被禁用、会话已关闭 |
| 权限不足 | `PermissionDenied` | 无某操作权限 |
| 数据库、Redis、下游 RPC 等异常 | `Internal` | 可带 `%v` 将底层 err 写入消息，便于排查；敏感信息勿写入 |

业务逻辑层内部（如 model、工具函数）可继续使用标准库 `error`；在 logic 层将 error 转为上述 status 再返回。

---

## 3. 安全与文案

- **防枚举**：登录、找回密码等接口中，「用户不存在」与「密码错误」建议使用同一文案（如「用户名或密码错误」）并配合 `Unauthenticated`，避免通过错误信息枚举有效用户名。
- **敏感信息**：`Internal` 错误的 message 可能被日志或监控采集，避免包含密码、token、未脱敏 PII。

---

## 4. 参考实现

- **Auth 服务**：`services/auth/internal/logic` 下各 logic 已统一使用 `status.Error` / `status.Errorf`，可按其方式延用。
- **User 服务**：`services/user/internal/logic` 中 GetUser、UpdateUser、BatchGetUsers 等已使用 `codes.InvalidArgument`、`codes.NotFound`、`codes.Internal`。

后续新增或改造 RPC 服务时，logic 层对外返回错误请遵循本文档。
