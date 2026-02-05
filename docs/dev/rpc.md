# Beehive IM RPC 服务设计文档

## 一、RPC 服务概述

Beehive IM 系统采用 gRPC 作为微服务间通信协议，使用 Protocol Buffers 定义服务接口。

### 1.1 技术选型

- **RPC 框架**: gRPC
- **接口定义**: Protocol Buffers (proto3)
- **服务发现**: etcd
- **负载均衡**: go-zero 内置负载均衡器
- **代码生成**: goctl

### 1.2 服务列表

| 服务名 | 端口 | 职责 |
|--------|------|------|
| User RPC | 8001 | 用户管理、认证、在线状态 |
| Friend RPC | 8002 | 好友关系、好友申请 |
| Chat RPC | 8004 | 会话管理、会话成员 |
| Message RPC | 8003 | 消息发送、历史消息 |
| File RPC | 8005 | 文件上传、下载、去重 |
| Search RPC | 8006 | 消息全文检索 |

### 1.3 服务注册

所有 RPC 服务启动时注册到 etcd：

```
Key: beehive.rpc.user
Value: 127.0.0.1:8001

Key: beehive.rpc.friend
Value: 127.0.0.1:8002
...
```

## 二、User RPC Service

### 2.1 Proto 定义

**文件路径**: `/opt/Beehive/api/proto/user/v1/user.proto`

```protobuf
syntax = "proto3";

package user;
option go_package = "./user";

// 用户服务
service UserService {
    // 用户注册
    rpc Register(RegisterRequest) returns (RegisterResponse);
    
    // 用户登录
    rpc Login(LoginRequest) returns (LoginResponse);
    
    // 获取用户信息
    rpc GetUserInfo(GetUserInfoRequest) returns (UserInfoResponse);
    
    // 批量获取用户信息
    rpc GetUsersByIds(GetUsersByIdsRequest) returns (UsersResponse);
    
    // 更新用户信息
    rpc UpdateUserInfo(UpdateUserInfoRequest) returns (CommonResponse);
    
    // 发送邮箱验证码
    rpc SendVerificationCode(SendCodeRequest) returns (CommonResponse);
    
    // 验证邮箱验证码
    rpc VerifyCode(VerifyCodeRequest) returns (CommonResponse);
    
    // 更新在线状态
    rpc UpdateOnlineStatus(UpdateOnlineStatusRequest) returns (CommonResponse);
    
    // 校验密码
    rpc VerifyPassword(VerifyPasswordRequest) returns (VerifyPasswordResponse);
}

// 注册请求
message RegisterRequest {
    string username = 1;
    string email = 2;
    string password = 3;
    string code = 4;  // 邮箱验证码
}

// 注册响应
message RegisterResponse {
    int64 user_id = 1;
    string token = 2;
}

// 登录请求
message LoginRequest {
    string account = 1;   // 用户名或邮箱
    string password = 2;
}

// 登录响应
message LoginResponse {
    int64 user_id = 1;
    string username = 2;
    string nickname = 3;
    string avatar = 4;
    string token = 5;
}

// 获取用户信息请求
message GetUserInfoRequest {
    int64 user_id = 1;
}

// 用户信息响应
message UserInfoResponse {
    int64 user_id = 1;
    string username = 2;
    string email = 3;
    string nickname = 4;
    string avatar = 5;
    int32 status = 6;
    bool online = 7;
    int64 last_login_at = 8;
    int64 created_at = 9;
}

// 批量获取用户信息请求
message GetUsersByIdsRequest {
    repeated int64 user_ids = 1;
}

// 批量用户信息响应
message UsersResponse {
    repeated UserInfoResponse users = 1;
}

// 更新用户信息请求
message UpdateUserInfoRequest {
    int64 user_id = 1;
    string nickname = 2;
    string avatar_url = 3;
}

// 发送验证码请求
message SendCodeRequest {
    string email = 1;
    string purpose = 2;  // register, reset_password, bind_email
}

// 验证验证码请求
message VerifyCodeRequest {
    string email = 1;
    string code = 2;
    string purpose = 3;
}

// 更新在线状态请求
message UpdateOnlineStatusRequest {
    int64 user_id = 1;
    bool online = 2;
}

// 校验密码请求
message VerifyPasswordRequest {
    int64 user_id = 1;
    string password = 2;
}

// 校验密码响应
message VerifyPasswordResponse {
    bool valid = 1;
}

// 通用响应
message CommonResponse {
    bool success = 1;
    string message = 2;
}
```

### 2.2 服务实现要点

#### 2.2.1 用户注册

1. 验证邮箱验证码（5分钟有效期）
2. 检查用户名和邮箱是否已存在
3. 密码使用 bcrypt 加密（cost=12）
4. 生成 JWT Token（7天有效期）
5. 缓存用户信息到 Redis

#### 2.2.2 用户登录

1. 支持用户名或邮箱登录
2. 验证密码（bcrypt.CompareHashAndPassword）
3. 生成 JWT Token
4. 更新最后登录时间
5. 缓存用户信息和在线状态到 Redis

#### 2.2.3 邮箱验证码

1. 生成 6 位随机数字验证码
2. 保存到数据库（有效期 5 分钟）
3. 通过网易邮箱 SMTP 发送
4. 限制：同一邮箱 1 分钟内只能发送一次

#### 2.2.4 在线状态管理

- Redis Key: `user:online:{user_id}`
- Value: `{"online": true, "last_active": 1705838400}`
- TTL: 5 分钟（心跳续期）

---

## 三、Friend RPC Service

### 3.1 Proto 定义

**文件路径**: `/opt/Beehive/api/proto/friend/v1/friend.proto`

```protobuf
syntax = "proto3";

package friend;
option go_package = "./friend";

// 好友服务
service FriendService {
    // 发送好友申请
    rpc SendFriendRequest(SendFriendRequestRequest) returns (SendFriendRequestResponse);
    
    // 处理好友申请
    rpc HandleFriendRequest(HandleFriendRequestRequest) returns (CommonResponse);
    
    // 获取好友申请列表
    rpc GetFriendRequests(GetFriendRequestsRequest) returns (FriendRequestsResponse);
    
    // 获取好友列表
    rpc GetFriends(GetFriendsRequest) returns (FriendsResponse);
    
    // 删除好友
    rpc DeleteFriend(DeleteFriendRequest) returns (CommonResponse);
    
    // 检查是否为好友
    rpc IsFriend(IsFriendRequest) returns (IsFriendResponse);
    
    // 更新好友备注
    rpc UpdateFriendRemark(UpdateFriendRemarkRequest) returns (CommonResponse);
}

// 发送好友申请请求
message SendFriendRequestRequest {
    int64 from_user_id = 1;
    int64 to_user_id = 2;
    string message = 3;
}

// 发送好友申请响应
message SendFriendRequestResponse {
    int64 request_id = 1;
}

// 处理好友申请请求
message HandleFriendRequestRequest {
    int64 request_id = 1;
    int64 user_id = 2;  // 当前操作用户ID
    bool accept = 3;    // true: 同意, false: 拒绝
}

// 获取好友申请列表请求
message GetFriendRequestsRequest {
    int64 user_id = 1;
}

// 好友申请列表响应
message FriendRequestsResponse {
    repeated FriendRequestInfo requests = 1;
}

// 好友申请信息
message FriendRequestInfo {
    int64 request_id = 1;
    int64 from_user_id = 2;
    int64 to_user_id = 3;
    string message = 4;
    int32 status = 5;  // 0: 待处理, 1: 已同意, 2: 已拒绝, 3: 已过期
    int64 created_at = 6;
}

// 获取好友列表请求
message GetFriendsRequest {
    int64 user_id = 1;
}

// 好友列表响应
message FriendsResponse {
    repeated FriendInfo friends = 1;
}

// 好友信息
message FriendInfo {
    int64 user_id = 1;
    int64 friend_id = 2;
    string remark = 3;
    int64 created_at = 4;
}

// 删除好友请求
message DeleteFriendRequest {
    int64 user_id = 1;
    int64 friend_id = 2;
}

// 检查是否为好友请求
message IsFriendRequest {
    int64 user_id = 1;
    int64 friend_id = 2;
}

// 检查是否为好友响应
message IsFriendResponse {
    bool is_friend = 1;
}

// 更新好友备注请求
message UpdateFriendRemarkRequest {
    int64 user_id = 1;
    int64 friend_id = 2;
    string remark = 3;
}

// 通用响应
message CommonResponse {
    bool success = 1;
    string message = 2;
}
```

### 3.2 服务实现要点

#### 3.2.1 发送好友申请

1. 验证目标用户是否存在（调用 User RPC）
2. 检查是否已经是好友
3. 检查是否已有待处理的申请
4. 限制：24小时内不能重复发送申请
5. 插入好友申请记录

#### 3.2.2 处理好友申请

1. 验证申请是否存在且待处理
2. 验证当前用户是否为被申请人
3. 更新申请状态
4. 如果同意：
   - 插入双向好友关系（两条记录）
   - 删除好友列表缓存
   - 推送通知给申请人

#### 3.2.3 删除好友

1. 删除双向好友关系
2. 删除好友列表缓存
3. 不删除历史会话和消息

---

## 四、Chat RPC Service

### 4.1 Proto 定义

**文件路径**: `/opt/Beehive/api/proto/chat/v1/chat.proto`

```protobuf
syntax = "proto3";

package chat;
option go_package = "./chat";

// 会话服务
service ChatService {
    // 创建会话
    rpc CreateConversation(CreateConversationRequest) returns (CreateConversationResponse);
    
    // 获取会话列表
    rpc GetConversations(GetConversationsRequest) returns (ConversationsResponse);
    
    // 获取会话详情
    rpc GetConversationDetail(GetConversationDetailRequest) returns (ConversationDetailResponse);
    
    // 添加会话成员
    rpc AddMembers(AddMembersRequest) returns (CommonResponse);
    
    // 移除会话成员
    rpc RemoveMembers(RemoveMembersRequest) returns (CommonResponse);
    
    // 更新未读数
    rpc UpdateUnreadCount(UpdateUnreadCountRequest) returns (CommonResponse);
    
    // 标记已读
    rpc MarkRead(MarkReadRequest) returns (CommonResponse);
    
    // 更新会话信息
    rpc UpdateConversationInfo(UpdateConversationInfoRequest) returns (CommonResponse);
    
    // 检查用户是否在会话中
    rpc CheckMember(CheckMemberRequest) returns (CheckMemberResponse);
}

// 创建会话请求
message CreateConversationRequest {
    int32 type = 1;  // 1: 单聊, 2: 群聊
    string name = 2;
    int64 owner_id = 3;
    repeated int64 member_ids = 4;
}

// 创建会话响应
message CreateConversationResponse {
    int64 conversation_id = 1;
}

// 获取会话列表请求
message GetConversationsRequest {
    int64 user_id = 1;
}

// 会话列表响应
message ConversationsResponse {
    repeated ConversationInfo conversations = 1;
}

// 会话信息
message ConversationInfo {
    int64 conversation_id = 1;
    int32 type = 2;
    string name = 3;
    string avatar = 4;
    int32 unread_count = 5;
    int64 last_message_at = 6;
}

// 获取会话详情请求
message GetConversationDetailRequest {
    int64 conversation_id = 1;
    int64 user_id = 2;
}

// 会话详情响应
message ConversationDetailResponse {
    int64 conversation_id = 1;
    int32 type = 2;
    string name = 3;
    string avatar = 4;
    int64 owner_id = 5;
    repeated MemberInfo members = 6;
}

// 会话成员信息
message MemberInfo {
    int64 user_id = 1;
    int32 role = 2;  // 1: 普通成员, 2: 管理员, 3: 群主
    int64 joined_at = 3;
}

// 添加成员请求
message AddMembersRequest {
    int64 conversation_id = 1;
    repeated int64 user_ids = 2;
}

// 移除成员请求
message RemoveMembersRequest {
    int64 conversation_id = 1;
    repeated int64 user_ids = 2;
}

// 更新未读数请求
message UpdateUnreadCountRequest {
    int64 conversation_id = 1;
    int64 user_id = 2;
    int32 increment = 3;  // 增加的数量（可以是负数）
}

// 标记已读请求
message MarkReadRequest {
    int64 conversation_id = 1;
    int64 user_id = 2;
}

// 更新会话信息请求
message UpdateConversationInfoRequest {
    int64 conversation_id = 1;
    string name = 2;
    string avatar = 3;
}

// 检查成员请求
message CheckMemberRequest {
    int64 conversation_id = 1;
    int64 user_id = 2;
}

// 检查成员响应
message CheckMemberResponse {
    bool is_member = 1;
    int32 role = 2;
}

// 通用响应
message CommonResponse {
    bool success = 1;
    string message = 2;
}
```

### 4.2 服务实现要点

#### 4.2.1 创建会话

**单聊**:
1. 检查两个用户是否为好友
2. 检查是否已存在单聊会话（防止重复创建）
3. 创建会话记录
4. 添加两个成员到 conversation_members

**群聊**:
1. 验证所有成员是否存在
2. 创建会话记录（群主为创建者）
3. 添加所有成员到 conversation_members
4. 群主角色为 3，其他成员为 1

#### 4.2.2 更新未读数

- 每次收到新消息时，会话所有成员（除发送者）未读数 +1
- 用户打开会话查看消息后，调用 MarkRead 清零

---

## 五、Message RPC Service

### 5.1 Proto 定义

**文件路径**: `/opt/Beehive/api/proto/message/v1/message.proto`

```protobuf
syntax = "proto3";

package message;
option go_package = "./message";

// 消息服务
service MessageService {
    // 发送消息
    rpc SendMessage(SendMessageRequest) returns (SendMessageResponse);
    
    // 获取历史消息
    rpc GetMessages(GetMessagesRequest) returns (MessagesResponse);
    
    // 更新消息状态
    rpc UpdateMessageStatus(UpdateMessageStatusRequest) returns (CommonResponse);
    
    // 获取消息详情
    rpc GetMessageDetail(GetMessageDetailRequest) returns (MessageInfo);
}

// 发送消息请求
message SendMessageRequest {
    int64 conversation_id = 1;
    int64 sender_id = 2;
    int32 content_type = 3;  // 1: 文本, 2: 图片, 3: 语音, 4: 文件
    string content = 4;
    string extra_data = 5;  // JSON 格式扩展数据
}

// 发送消息响应
message SendMessageResponse {
    int64 message_id = 1;
    int64 created_at = 2;
}

// 获取历史消息请求
message GetMessagesRequest {
    int64 conversation_id = 1;
    int64 before_id = 2;  // 消息ID，获取该ID之前的消息（游标分页）
    int32 limit = 3;      // 每页数量，默认50
}

// 历史消息响应
message MessagesResponse {
    repeated MessageInfo messages = 1;
    bool has_more = 2;
}

// 消息信息
message MessageInfo {
    int64 message_id = 1;
    int64 conversation_id = 2;
    int64 sender_id = 3;
    int32 content_type = 4;
    string content = 5;
    string extra_data = 6;
    int32 status = 7;  // 0: 发送中, 1: 已送达, 2: 已读
    int64 created_at = 8;
}

// 更新消息状态请求
message UpdateMessageStatusRequest {
    int64 message_id = 1;
    int32 status = 2;
}

// 获取消息详情请求
message GetMessageDetailRequest {
    int64 message_id = 1;
}

// 通用响应
message CommonResponse {
    bool success = 1;
    string message = 2;
}
```

### 5.2 服务实现要点

#### 5.2.1 发送消息

1. 验证会话存在且用户是会话成员
2. 保存消息到数据库
3. 发布消息到 RabbitMQ（3个队列）：
   - `message.persist`: 持久化队列（已在步骤2完成）
   - `message.push`: 推送队列（Gateway 消费）
   - `message.index`: 索引队列（Search RPC 消费）
4. 更新会话未读数（调用 Chat RPC）
5. 返回消息 ID 和时间戳

#### 5.2.2 获取历史消息

- 使用游标分页（based on message_id）
- 按时间倒序返回
- 先查 Redis 缓存（热点消息）
- 缓存未命中则查数据库

---

## 六、File RPC Service

### 6.1 Proto 定义

**文件路径**: `/opt/Beehive/api/proto/file/v1/file.proto`

```protobuf
syntax = "proto3";

package file;
option go_package = "./file";

// 文件服务
service FileService {
    // 上传文件
    rpc UploadFile(UploadFileRequest) returns (UploadFileResponse);
    
    // 获取文件信息
    rpc GetFileInfo(GetFileInfoRequest) returns (FileInfoResponse);
    
    // 删除文件
    rpc DeleteFile(DeleteFileRequest) returns (CommonResponse);
    
    // 检查文件是否存在（去重）
    rpc CheckFileExists(CheckFileExistsRequest) returns (CheckFileExistsResponse);
    
    // 初始化分片上传
    rpc InitMultipartUpload(InitMultipartUploadRequest) returns (InitMultipartUploadResponse);
    
    // 上传文件分片
    rpc UploadPart(UploadPartRequest) returns (UploadPartResponse);
    
    // 完成分片上传
    rpc CompleteMultipartUpload(CompleteMultipartUploadRequest) returns (UploadFileResponse);
}

// 上传文件请求
message UploadFileRequest {
    int64 user_id = 1;
    string file_name = 2;
    int64 file_size = 3;
    string file_type = 4;  // MIME类型
    string file_hash = 5;  // SHA256哈希
    bytes file_data = 6;
}

// 上传文件响应
message UploadFileResponse {
    int64 file_id = 1;
    string file_url = 2;
    string file_hash = 3;
    bool is_duplicate = 4;  // 是否为去重文件
}

// 获取文件信息请求
message GetFileInfoRequest {
    int64 file_id = 1;
}

// 文件信息响应
message FileInfoResponse {
    int64 file_id = 1;
    string file_name = 2;
    int64 file_size = 3;
    string file_type = 4;
    string file_url = 5;
    string file_hash = 6;
}

// 删除文件请求
message DeleteFileRequest {
    int64 file_id = 1;
}

// 检查文件是否存在请求
message CheckFileExistsRequest {
    string file_hash = 1;
}

// 检查文件是否存在响应
message CheckFileExistsResponse {
    bool exists = 1;
    int64 file_id = 2;
    string file_url = 3;
}

// 初始化分片上传请求
message InitMultipartUploadRequest {
    int64 user_id = 1;
    string file_name = 2;
    int64 file_size = 3;
    string file_type = 4;
    string file_hash = 5;
    int32 total_parts = 6;  // 总分片数
}

// 初始化分片上传响应
message InitMultipartUploadResponse {
    string upload_id = 1;
}

// 上传分片请求
message UploadPartRequest {
    string upload_id = 1;
    int32 part_number = 2;  // 分片编号（从1开始）
    bytes part_data = 3;
}

// 上传分片响应
message UploadPartResponse {
    bool success = 1;
    string etag = 2;  // 分片标识
}

// 完成分片上传请求
message CompleteMultipartUploadRequest {
    string upload_id = 1;
    repeated PartInfo parts = 2;
}

// 分片信息
message PartInfo {
    int32 part_number = 1;
    string etag = 2;
}

// 通用响应
message CommonResponse {
    bool success = 1;
    string message = 2;
}
```

### 6.2 服务实现要点

#### 6.2.1 文件去重

1. 客户端计算文件 SHA256 哈希
2. 调用 `CheckFileExists` 检查哈希是否存在
3. 如果存在：
   - 直接返回已有文件 URL
   - `usage_count` +1
   - 跳过文件上传
4. 如果不存在：
   - 上传文件
   - 保存文件记录

#### 6.2.2 断点续传

1. 客户端调用 `InitMultipartUpload` 初始化
2. 客户端将文件切分为多个分片（每个5MB）
3. 并发上传各个分片 `UploadPart`
4. 所有分片上传完成后，调用 `CompleteMultipartUpload`
5. 服务端合并分片，保存文件

---

## 七、Search RPC Service

### 7.1 Proto 定义

**文件路径**: `/opt/Beehive/api/proto/search/v1/search.proto`

```protobuf
syntax = "proto3";

package search;
option go_package = "./search";

// 搜索服务
service SearchService {
    // 搜索消息
    rpc SearchMessages(SearchMessagesRequest) returns (SearchMessagesResponse);
    
    // 索引消息
    rpc IndexMessage(IndexMessageRequest) returns (CommonResponse);
    
    // 删除消息索引
    rpc DeleteMessageIndex(DeleteMessageIndexRequest) returns (CommonResponse);
}

// 搜索消息请求
message SearchMessagesRequest {
    string keyword = 1;
    int64 user_id = 2;
    int64 conversation_id = 3;  // 可选，按会话过滤
    int32 page = 4;
    int32 page_size = 5;
}

// 搜索消息响应
message SearchMessagesResponse {
    repeated SearchMessageInfo messages = 1;
    int64 total = 2;
}

// 搜索消息信息
message SearchMessageInfo {
    int64 message_id = 1;
    int64 conversation_id = 2;
    int64 sender_id = 3;
    string content = 4;
    string highlight = 5;  // 高亮内容
    int64 created_at = 6;
}

// 索引消息请求
message IndexMessageRequest {
    int64 message_id = 1;
    int64 conversation_id = 2;
    int64 sender_id = 3;
    string content = 4;
    int64 created_at = 5;
}

// 删除消息索引请求
message DeleteMessageIndexRequest {
    int64 message_id = 1;
}

// 通用响应
message CommonResponse {
    bool success = 1;
    string message = 2;
}
```

### 7.2 服务实现要点

#### 7.2.1 消息索引

1. 监听 RabbitMQ `message.index` 队列
2. 接收到消息后调用 `IndexMessage`
3. 索引消息到 Elasticsearch
4. 仅索引文本消息内容

#### 7.2.2 消息搜索

1. 使用 Elasticsearch 全文检索
2. 支持高亮显示
3. 支持按会话过滤
4. 分页返回结果

---

## 八、代码生成

### 8.1 生成 RPC 代码

```bash
# 进入 proto 目录
cd /opt/Beehive/api/proto

# 生成 User RPC
cd user/v1
goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=../../../rpc/user

# 生成 Friend RPC
cd ../friend/v1
goctl rpc protoc friend.proto --go_out=. --go-grpc_out=. --zrpc_out=../../../rpc/friend

# 生成 Chat RPC
cd ../chat/v1
goctl rpc protoc chat.proto --go_out=. --go-grpc_out=. --zrpc_out=../../../rpc/chat

# 生成 Message RPC
cd ../message/v1
goctl rpc protoc message.proto --go_out=. --go-grpc_out=. --zrpc_out=../../../rpc/message

# 生成 File RPC
cd ../file/v1
goctl rpc protoc file.proto --go_out=. --go-grpc_out=. --zrpc_out=../../../rpc/file

# 生成 Search RPC
cd ../search/v1
goctl rpc protoc search.proto --go_out=. --go-grpc_out=. --zrpc_out=../../../rpc/search
```

### 8.2 生成的文件结构

```
rpc/user/
├── etc/
│   └── user.yaml          # 配置文件
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── logic/
│   │   ├── registerlogic.go
│   │   ├── loginlogic.go
│   │   └── ...
│   ├── server/
│   │   └── userserviceserver.go
│   └── svc/
│       └── servicecontext.go
├── model/                 # goctl model 生成
│   ├── usersmodel.go
│   └── emailverificationcodesmodel.go
├── user/                  # pb.go
│   ├── user.pb.go
│   └── user_grpc.pb.go
└── user.go               # main 入口
```

---

## 九、RPC 调用示例

### 9.1 Gateway 调用 User RPC

```go
// 在 Gateway 的 ServiceContext 中注入 User RPC Client
type ServiceContext struct {
    Config    config.Config
    UserRpc   userclient.UserService
    FriendRpc friendclient.FriendService
    // ...
}

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config:    c,
        UserRpc:   userclient.NewUserService(zrpc.MustNewClient(c.UserRpc)),
        FriendRpc: friendclient.NewFriendService(zrpc.MustNewClient(c.FriendRpc)),
    }
}

// 在 Logic 中调用
func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
    resp, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginRequest{
        Account:  req.Account,
        Password: req.Password,
    })
    if err != nil {
        return nil, err
    }
    
    return &types.LoginResp{
        UserId:   resp.UserId,
        Username: resp.Username,
        Nickname: resp.Nickname,
        Avatar:   resp.Avatar,
        Token:    resp.Token,
    }, nil
}
```

### 9.2 Message RPC 调用 Chat RPC

```go
// 在 Message RPC 的 ServiceContext 中注入 Chat RPC Client
type ServiceContext struct {
    Config  config.Config
    ChatRpc chatclient.ChatService
}

// 发送消息后更新未读数
func (l *SendMessageLogic) SendMessage(req *message.SendMessageRequest) (*message.SendMessageResponse, error) {
    // 1. 保存消息到数据库
    result, err := l.svcCtx.MessageModel.Insert(l.ctx, &model.Message{
        ConversationId: req.ConversationId,
        SenderId:       req.SenderId,
        ContentType:    req.ContentType,
        Content:        req.Content,
        ExtraData:      sql.NullString{String: req.ExtraData, Valid: true},
    })
    // ...
    
    // 2. 更新会话未读数
    _, err = l.svcCtx.ChatRpc.UpdateUnreadCount(l.ctx, &chat.UpdateUnreadCountRequest{
        ConversationId: req.ConversationId,
        UserId:         req.SenderId,
        Increment:      1,
    })
    
    // 3. 发布到 RabbitMQ
    // ...
    
    return &message.SendMessageResponse{
        MessageId: messageId,
        CreatedAt: time.Now().Unix(),
    }, nil
}
```

---

## 十、性能优化

### 10.1 连接池

go-zero 内置连接池，配置示例：

```yaml
# etc/gateway.yaml
UserRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: beehive.rpc.user
  Timeout: 5000  # 超时时间（毫秒）
```

### 10.2 负载均衡

go-zero 默认使用轮询负载均衡，支持：
- p2c (Power of Two Choices) - 默认
- random
- roundrobin
- consistent hash

### 10.3 熔断器

go-zero 自带自适应熔断器，自动开启，无需配置。

### 10.4 超时控制

所有 RPC 调用建议设置超时：

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
defer cancel()

resp, err := l.svcCtx.UserRpc.GetUserInfo(ctx, &user.GetUserInfoRequest{
    UserId: userId,
})
```

---

## 十一、监控和日志

### 11.1 RPC 日志

go-zero 自动记录 RPC 调用日志：

```
{"@timestamp":"2026-01-21T10:00:00.000+08:00","caller":"internal/server.go:123","content":"rpc call","duration":"10ms","level":"info","method":"/user.UserService/Login"}
```

### 11.2 链路追踪

go-zero 集成 OpenTelemetry，自动生成 TraceID：

```go
// 从 Context 中提取 TraceID
traceId := trace.TraceIDFromContext(ctx)
```

### 11.3 Prometheus 指标

go-zero 自动暴露 Prometheus 指标：

```
# RPC 请求总数
rpc_server_requests_total{method="/user.UserService/Login"} 1000

# RPC 请求延迟
rpc_server_request_duration_ms_bucket{method="/user.UserService/Login",le="100"} 950
```

---

## 十二、常见问题

### Q1: RPC 调用超时怎么办？

A: 
1. 检查服务是否启动
2. 检查 etcd 服务发现是否正常
3. 增加超时时间
4. 检查网络连接

### Q2: 如何调试 RPC 服务？

A:
1. 查看日志文件
2. 使用 grpcurl 测试：`grpcurl -plaintext localhost:8001 list`
3. 使用 go-zero 内置的 debug 端点

### Q3: 如何热更新配置？

A: go-zero 暂不支持热更新，需要重启服务。建议使用配置中心（如 Nacos）。

### Q4: 如何实现 RPC 服务的灰度发布？

A: 使用 Kubernetes + Istio 实现流量管理和灰度发布。

---

## 十三、最佳实践

1. **错误处理**: 使用 go-zero 的 `errorx` 包统一错误码
2. **参数校验**: 在 RPC 层做参数校验
3. **幂等性**: 对于可能重复调用的接口，实现幂等性
4. **版本管理**: Proto 文件使用版本号（v1, v2）
5. **文档维护**: 及时更新 Proto 文件注释
