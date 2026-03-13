## 内部 gRPC 接口设计（Auth / Presence / Message / Conversation）

本文件根据 `docs/backend/services-design.md` 描述的服务职责，对核心内部服务的 gRPC 接口进行设计，后续可直接映射为 `proto/*.proto`，使用 goctl 生成 zrpc 服务端骨架，并使用 protoc 生成各服务目录下的 RPC client 代码。

> 说明：以下接口以「伪 proto 风格」描述，实际 `.proto` 文件可在实现时按需要拆分/调整字段类型。

---

### 1. AuthService

**职责参考**：`docs/backend/services-design.md` 中的 AuthService。

#### 1.1 服务定义示意

```proto
service AuthService {
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc TokenLogin(TokenLoginRequest) returns (LoginResponse);
  rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
  rpc Logout(LogoutRequest) returns (LogoutResponse);
  // RBAC：查询用户系统级角色
  rpc GetUserRoles(GetUserRolesRequest) returns (GetUserRolesResponse);
  // RBAC：检查用户是否具备某个权限
  rpc CheckPermission(CheckPermissionRequest) returns (CheckPermissionResponse);
  // RBAC（可选，对内/管理用途）：为用户设置角色
  rpc AssignRoles(AssignRolesRequest) returns (AssignRolesResponse);
}
```

#### 1.2 关键消息体（示意）

```proto
message LoginRequest {
  string username = 1;
  string password = 2;
  string device_id = 3;
}

message LoginResponse {
  string user_id = 1;
  string access_token = 2;
  string refresh_token = 3;
  int64  expires_in = 4;
}

message TokenLoginRequest {
  string access_token = 1;
  string device_id = 2;
}

message ValidateTokenRequest {
  string access_token = 1;
}

message ValidateTokenResponse {
  bool   valid = 1;
  string user_id = 2;
}

// RBAC：用户角色与权限

// 约定若干系统级角色：
// - user：普通用户，默认角色
// - admin：管理员，可访问 Admin 后台的大部分只读接口
// - super_admin：超级管理员，可进行封禁、配置变更等高危操作
//
// 权限使用字符串编码，按照「域.资源.动作」风格命名，例如：
// - admin.user.read
// - admin.user.ban
// - admin.config.write

message GetUserRolesRequest {
  string user_id = 1;
}

message GetUserRolesResponse {
  repeated string roles = 1;
}

message CheckPermissionRequest {
  string user_id = 1;
  string permission = 2;
}

message CheckPermissionResponse {
  bool allowed = 1;
}

message AssignRolesRequest {
  string user_id = 1;
  // 完整覆盖式赋值：实现时建议先清空再写入
  repeated string roles = 2;
}

message AssignRolesResponse {}

> 集成建议：
>
> - Gateway 在处理需要管理权限的 WebSocket 操作前，可调用 `CheckPermission` 校验当前用户是否具备如 `admin.user.ban` 等权限；
> - AdminHTTP 中间件在路由前统一调用 `ValidateToken` + `CheckPermission`，对不同路由绑定不同的权限编码。
```

---

### 2. PresenceService

**职责参考**：`docs/backend/services-design.md` 中的 PresenceService。

#### 2.1 服务定义示意

```proto
service PresenceService {
  rpc RegisterSession(RegisterSessionRequest) returns (RegisterSessionResponse);
  rpc UnregisterSession(UnregisterSessionRequest) returns (UnregisterSessionResponse);
  rpc RefreshSession(RefreshSessionRequest) returns (RefreshSessionResponse);
  rpc GetOnlineSessions(GetOnlineSessionsRequest) returns (GetOnlineSessionsResponse);
  rpc GetUserPresence(GetUserPresenceRequest) returns (GetUserPresenceResponse);
}
```

#### 2.2 关键消息体（示意）

```proto
message RegisterSessionRequest {
  string user_id = 1;
  string gateway_id = 2;
  string conn_id = 3;
  string device_id = 4;
  string device_type = 5;
  string ip = 6;
}

message RegisterSessionResponse {}

message UnregisterSessionRequest {
  string user_id = 1;
  string conn_id = 2;
}

message UnregisterSessionResponse {}

message RefreshSessionRequest {
  string user_id = 1;
  string conn_id = 2;
}

message RefreshSessionResponse {}

message GetOnlineSessionsRequest {
  string user_id = 1;
}

message SessionInfo {
  string gateway_id = 1;
  string conn_id = 2;
  string device_id = 3;
  string device_type = 4;
  int64  last_ping_at = 5;
}

message GetOnlineSessionsResponse {
  repeated SessionInfo sessions = 1;
}

message GetUserPresenceRequest {
  string user_id = 1;
}

message GetUserPresenceResponse {
  bool online = 1;
  repeated SessionInfo sessions = 2;
}
```

---

### 3. MessageService

**职责参考**：负责消息写入/历史查询，并在写入成功后发布 `message.created` 事件到 RabbitMQ。

#### 3.1 服务定义示意

```proto
service MessageService {
  rpc PostMessage(PostMessageRequest) returns (PostMessageResponse);
  rpc GetHistory(GetHistoryRequest) returns (GetHistoryResponse);
  rpc GetLastMessages(GetLastMessagesRequest) returns (GetLastMessagesResponse);
}
```

#### 3.2 关键消息体（示意）

```proto
message MessageBody {
  string type = 1;   // text / image / system / ...
  string text = 2;   // 当 type = text 时使用
  // 后续可扩展 image_url 等字段
}

message PostMessageRequest {
  string client_msg_id = 1;
  string conversation_id = 2;
  string from_user_id = 3;
  string to_user_id = 4;      // 单聊可用，群聊可为空
  MessageBody body = 5;
}

message PostMessageResponse {
  string server_msg_id = 1;
  string conversation_id = 2;
  int64  server_time = 3;
}

message GetHistoryRequest {
  string conversation_id = 1;
  int64  before_time = 2;     // 可选：某个时间之前
  int32  limit = 3;
}

message MessageRecord {
  string server_msg_id = 1;
  string conversation_id = 2;
  string from_user_id = 3;
  string to_user_id = 4;
  MessageBody body = 5;
  int64  server_time = 6;
}

message GetHistoryResponse {
  repeated MessageRecord items = 1;
  bool has_more = 2;
}

message GetLastMessagesRequest {
  repeated string conversation_ids = 1;
}

message GetLastMessagesResponse {
  map<string, MessageRecord> last_messages = 1;
}
```

---

### 4. ConversationService

**职责参考**：负责会话/群组元数据与成员关系，提供会话列表等能力。

#### 4.1 服务定义示意

```proto
service ConversationService {
  rpc CreateConversation(CreateConversationRequest) returns (CreateConversationResponse);
  rpc AddMember(AddMemberRequest) returns (AddMemberResponse);
  rpc RemoveMember(RemoveMemberRequest) returns (RemoveMemberResponse);
  rpc ListUserConversations(ListUserConversationsRequest) returns (ListUserConversationsResponse);
  rpc GetConversation(GetConversationRequest) returns (GetConversationResponse);
  rpc ListMembers(ListMembersRequest) returns (ListMembersResponse);
}
```

#### 4.2 关键消息体（示意）

```proto
message CreateConversationRequest {
  string type = 1;         // single / group / channel
  string name = 2;
  repeated string member_ids = 3;
}

message CreateConversationResponse {
  string conversation_id = 1;
}

message AddMemberRequest {
  string conversation_id = 1;
  string user_id = 2;
  string role = 3;         // owner / admin / member
}

message AddMemberResponse {}

message RemoveMemberRequest {
  string conversation_id = 1;
  string user_id = 2;
}

message RemoveMemberResponse {}

message ConversationInfo {
  string id = 1;
  string type = 2;
  string name = 3;
  int32  member_count = 4;
  int64  created_at = 5;
  int64  last_active_at = 6;
}

message ListUserConversationsRequest {
  string user_id = 1;
  string cursor = 2;
  int32  limit = 3;
}

message ListUserConversationsResponse {
  repeated ConversationInfo items = 1;
  string next_cursor = 2;
}

message GetConversationRequest {
  string id = 1;
}

message GetConversationResponse {
  ConversationInfo conversation = 1;
}

message ListMembersRequest {
  string conversation_id = 1;
}

message MemberInfo {
  string user_id = 1;
  string role = 2;
  int64  joined_at = 3;
  string status = 4;   // active / left / banned
}

message ListMembersResponse {
  repeated MemberInfo items = 1;
}
```

---

### 5. 设计落地建议

- 每个服务单独维护一个 `.proto` 文件（如 `auth.proto`, `presence.proto`, `message.proto`, `conversation.proto`），放在仓库 `proto/` 目录。
- 使用 `goctl rpc protoc` 生成对应 zrpc 服务代码，放到 `services/<service>/` 下：
  - 例如：
    - `goctl rpc protoc proto/auth.proto --go_out=. --go-grpc_out=. --zrpc_out=services/auth`
- 上述消息体仅为初始设计，可在实现阶段根据具体字段需求微调，但应保持与 `docs/API/websocket-client-api.md` 中的字段语义一致。
