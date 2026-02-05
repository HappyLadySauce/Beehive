# Beehive IM API 接口文档

## 一、接口概述

### 1.1 基本信息

- **Base URL**: `http://localhost:8888`
- **API 版本**: v1
- **认证方式**: JWT Bearer Token
- **请求格式**: `application/json`
- **响应格式**: `application/json`
- **字符编码**: UTF-8

### 1.2 通用响应格式

**成功响应**:

```json
{
    "data": { ... },  // 响应数据
    "timestamp": 1705838400
}
```

**错误响应**:

```json
{
    "code": 40001,
    "message": "用户名或密码错误",
    "timestamp": 1705838400
}
```

### 1.3 错误码

| 错误码 | 说明 |
|--------|------|
| 40001 | 参数错误 |
| 40002 | 用户名已存在 |
| 40003 | 邮箱已存在 |
| 40004 | 验证码错误或已过期 |
| 40005 | 用户名或密码错误 |
| 40006 | Token 无效或已过期 |
| 40007 | 权限不足 |
| 40008 | 好友已存在 |
| 40009 | 好友申请不存在 |
| 40010 | 会话不存在 |
| 40011 | 文件上传失败 |
| 50001 | 服务器内部错误 |
| 50002 | 数据库错误 |
| 50003 | 邮件发送失败 |

### 1.4 认证说明

除了注册和登录接口，所有接口都需要在 Header 中携带 JWT Token：

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## 二、用户相关接口

### 2.1 发送邮箱验证码

**接口**: `POST /api/v1/auth/send-code`

**描述**: 发送邮箱验证码，用于注册或重置密码

**请求参数**:

```json
{
    "email": "user@example.com",
    "purpose": "register"  // register | reset_password
}
```

**响应**:

```json
{
    "success": true
}
```

**示例**:

```bash
curl -X POST http://localhost:8888/api/v1/auth/send-code \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "purpose": "register"
  }'
```

---

### 2.2 用户注册

**接口**: `POST /api/v1/auth/register`

**描述**: 用户注册，需要先发送验证码

**请求参数**:

```json
{
    "username": "testuser",
    "email": "user@example.com",
    "password": "password123",
    "code": "123456"
}
```

**参数说明**:
- `username`: 用户名，3-50字符，字母数字下划线
- `email`: 邮箱地址
- `password`: 密码，最少6位
- `code`: 邮箱验证码

**响应**:

```json
{
    "user_id": 1001,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**示例**:

```bash
curl -X POST http://localhost:8888/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "code": "123456"
  }'
```

---

### 2.3 用户登录

**接口**: `POST /api/v1/auth/login`

**描述**: 用户登录，支持用户名或邮箱登录

**请求参数**:

```json
{
    "account": "testuser",  // 用户名或邮箱
    "password": "password123"
}
```

**响应**:

```json
{
    "user_id": 1001,
    "username": "testuser",
    "nickname": "测试用户",
    "avatar": "https://example.com/avatar.jpg",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**示例**:

```bash
curl -X POST http://localhost:8888/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "account": "testuser",
    "password": "password123"
  }'
```

---

### 2.4 获取用户信息

**接口**: `GET /api/v1/users/:user_id`

**描述**: 获取指定用户信息

**请求 Header**:

```
Authorization: Bearer {token}
```

**响应**:

```json
{
    "user_id": 1001,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "测试用户",
    "avatar": "https://example.com/avatar.jpg",
    "status": 1,
    "online": true,
    "last_login_at": 1705838400,
    "created_at": 1705752000
}
```

**示例**:

```bash
curl -X GET http://localhost:8888/api/v1/users/1001 \
  -H "Authorization: Bearer {token}"
```

---

### 2.5 更新用户信息

**接口**: `PUT /api/v1/users/me`

**描述**: 更新当前登录用户信息

**请求 Header**:

```
Authorization: Bearer {token}
```

**请求参数**:

```json
{
    "nickname": "新昵称",
    "avatar_url": "https://example.com/new-avatar.jpg"
}
```

**响应**:

```json
{
    "success": true
}
```

**示例**:

```bash
curl -X PUT http://localhost:8888/api/v1/users/me \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "nickname": "新昵称"
  }'
```

---

## 三、好友相关接口

### 3.1 发送好友申请

**接口**: `POST /api/v1/friends/request`

**描述**: 向指定用户发送好友申请

**请求 Header**:

```
Authorization: Bearer {token}
```

**请求参数**:

```json
{
    "to_user_id": 1002,
    "message": "你好，我想加你为好友"
}
```

**响应**:

```json
{
    "request_id": 5001
}
```

**示例**:

```bash
curl -X POST http://localhost:8888/api/v1/friends/request \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "to_user_id": 1002,
    "message": "你好，我想加你为好友"
  }'
```

---

### 3.2 处理好友申请

**接口**: `POST /api/v1/friends/request/handle`

**描述**: 同意或拒绝好友申请

**请求 Header**:

```
Authorization: Bearer {token}
```

**请求参数**:

```json
{
    "request_id": 5001,
    "accept": true  // true: 同意, false: 拒绝
}
```

**响应**:

```json
{
    "success": true
}
```

**示例**:

```bash
curl -X POST http://localhost:8888/api/v1/friends/request/handle \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "request_id": 5001,
    "accept": true
  }'
```

---

### 3.3 获取好友申请列表

**接口**: `GET /api/v1/friends/requests`

**描述**: 获取当前用户收到的好友申请列表

**请求 Header**:

```
Authorization: Bearer {token}
```

**响应**:

```json
{
    "requests": [
        {
            "request_id": 5001,
            "from_user_id": 1003,
            "from_user": {
                "user_id": 1003,
                "username": "friend1",
                "nickname": "朋友1",
                "avatar": "https://example.com/avatar1.jpg",
                "online": true
            },
            "message": "你好，我想加你为好友",
            "status": 0,  // 0: 待处理, 1: 已同意, 2: 已拒绝
            "created_at": 1705838400
        }
    ]
}
```

**示例**:

```bash
curl -X GET http://localhost:8888/api/v1/friends/requests \
  -H "Authorization: Bearer {token}"
```

---

### 3.4 获取好友列表

**接口**: `GET /api/v1/friends`

**描述**: 获取当前用户的好友列表

**请求 Header**:

```
Authorization: Bearer {token}
```

**响应**:

```json
{
    "friends": [
        {
            "user_id": 1002,
            "username": "friend1",
            "nickname": "朋友1",
            "avatar": "https://example.com/avatar1.jpg",
            "remark": "同事",
            "online": true,
            "created_at": 1705752000
        },
        {
            "user_id": 1003,
            "username": "friend2",
            "nickname": "朋友2",
            "avatar": "https://example.com/avatar2.jpg",
            "remark": "",
            "online": false,
            "created_at": 1705838400
        }
    ]
}
```

**示例**:

```bash
curl -X GET http://localhost:8888/api/v1/friends \
  -H "Authorization: Bearer {token}"
```

---

### 3.5 删除好友

**接口**: `DELETE /api/v1/friends/:friend_id`

**描述**: 删除指定好友

**请求 Header**:

```
Authorization: Bearer {token}
```

**响应**:

```json
{
    "success": true
}
```

**示例**:

```bash
curl -X DELETE http://localhost:8888/api/v1/friends/1002 \
  -H "Authorization: Bearer {token}"
```

---

## 四、会话相关接口

### 4.1 创建会话

**接口**: `POST /api/v1/conversations`

**描述**: 创建单聊或群聊会话

**请求 Header**:

```
Authorization: Bearer {token}
```

**请求参数**:

```json
{
    "type": 1,  // 1: 单聊, 2: 群聊
    "name": "技术交流群",  // 群聊名称（单聊不需要）
    "member_ids": [1002, 1003, 1004]  // 成员ID列表
}
```

**响应**:

```json
{
    "conversation_id": 2001
}
```

**示例**:

```bash
# 创建单聊
curl -X POST http://localhost:8888/api/v1/conversations \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "type": 1,
    "member_ids": [1002]
  }'

# 创建群聊
curl -X POST http://localhost:8888/api/v1/conversations \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "type": 2,
    "name": "技术交流群",
    "member_ids": [1002, 1003, 1004]
  }'
```

---

### 4.2 获取会话列表

**接口**: `GET /api/v1/conversations`

**描述**: 获取当前用户的所有会话

**请求 Header**:

```
Authorization: Bearer {token}
```

**响应**:

```json
{
    "conversations": [
        {
            "conversation_id": 2001,
            "type": 1,
            "name": "friend1",  // 单聊显示对方昵称，群聊显示群名
            "avatar": "https://example.com/avatar1.jpg",
            "last_message": "你好",
            "last_message_at": 1705838400,
            "unread_count": 3,
            "members": [
                {
                    "user_id": 1002,
                    "username": "friend1",
                    "nickname": "朋友1",
                    "avatar": "https://example.com/avatar1.jpg",
                    "role": 1
                }
            ]
        },
        {
            "conversation_id": 2002,
            "type": 2,
            "name": "技术交流群",
            "avatar": "https://example.com/group-avatar.jpg",
            "last_message": "大家好",
            "last_message_at": 1705838300,
            "unread_count": 0,
            "members": [
                {
                    "user_id": 1002,
                    "username": "friend1",
                    "nickname": "朋友1",
                    "avatar": "https://example.com/avatar1.jpg",
                    "role": 1
                },
                {
                    "user_id": 1003,
                    "username": "friend2",
                    "nickname": "朋友2",
                    "avatar": "https://example.com/avatar2.jpg",
                    "role": 1
                }
            ]
        }
    ]
}
```

**示例**:

```bash
curl -X GET http://localhost:8888/api/v1/conversations \
  -H "Authorization: Bearer {token}"
```

---

### 4.3 获取会话详情

**接口**: `GET /api/v1/conversations/:conversation_id`

**描述**: 获取指定会话的详细信息

**请求 Header**:

```
Authorization: Bearer {token}
```

**响应**:

```json
{
    "conversation_id": 2002,
    "type": 2,
    "name": "技术交流群",
    "avatar": "https://example.com/group-avatar.jpg",
    "last_message": "大家好",
    "last_message_at": 1705838300,
    "unread_count": 0,
    "members": [
        {
            "user_id": 1001,
            "username": "testuser",
            "nickname": "测试用户",
            "avatar": "https://example.com/avatar.jpg",
            "role": 3  // 1: 普通成员, 2: 管理员, 3: 群主
        },
        {
            "user_id": 1002,
            "username": "friend1",
            "nickname": "朋友1",
            "avatar": "https://example.com/avatar1.jpg",
            "role": 1
        }
    ]
}
```

**示例**:

```bash
curl -X GET http://localhost:8888/api/v1/conversations/2002 \
  -H "Authorization: Bearer {token}"
```

---

### 4.4 标记消息已读

**接口**: `POST /api/v1/conversations/mark-read`

**描述**: 标记指定会话的消息为已读（清空未读数）

**请求 Header**:

```
Authorization: Bearer {token}
```

**请求参数**:

```json
{
    "conversation_id": 2001
}
```

**响应**:

```json
{
    "success": true
}
```

**示例**:

```bash
curl -X POST http://localhost:8888/api/v1/conversations/mark-read \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "conversation_id": 2001
  }'
```

---

## 五、消息相关接口

### 5.1 获取历史消息

**接口**: `GET /api/v1/conversations/:conversation_id/messages`

**描述**: 获取指定会话的历史消息（分页）

**请求 Header**:

```
Authorization: Bearer {token}
```

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| before_id | int64 | 否 | 消息ID，获取该ID之前的消息（分页游标） |
| limit | int32 | 否 | 每页数量，默认50 |

**响应**:

```json
{
    "messages": [
        {
            "message_id": 3001,
            "conversation_id": 2001,
            "sender_id": 1002,
            "sender_name": "朋友1",
            "sender_avatar": "https://example.com/avatar1.jpg",
            "content_type": 1,  // 1: 文本, 2: 图片, 3: 语音, 4: 文件
            "content": "你好",
            "extra_data": null,
            "status": 2,  // 0: 发送中, 1: 已送达, 2: 已读
            "created_at": 1705838400
        },
        {
            "message_id": 3002,
            "conversation_id": 2001,
            "sender_id": 1001,
            "sender_name": "测试用户",
            "sender_avatar": "https://example.com/avatar.jpg",
            "content_type": 2,
            "content": "https://example.com/images/photo.jpg",
            "extra_data": "{\"width\":1920,\"height\":1080,\"size\":204800}",
            "status": 1,
            "created_at": 1705838500
        }
    ],
    "has_more": true
}
```

**示例**:

```bash
# 获取最新50条消息
curl -X GET "http://localhost:8888/api/v1/conversations/2001/messages?limit=50" \
  -H "Authorization: Bearer {token}"

# 获取消息ID 3000 之前的50条消息
curl -X GET "http://localhost:8888/api/v1/conversations/2001/messages?before_id=3000&limit=50" \
  -H "Authorization: Bearer {token}"
```

---

### 5.2 搜索消息

**接口**: `GET /api/v1/messages/search`

**描述**: 全文搜索历史消息

**请求 Header**:

```
Authorization: Bearer {token}
```

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 是 | 搜索关键词 |
| conversation_id | int64 | 否 | 会话ID（不填则搜索所有会话） |
| page | int32 | 否 | 页码，默认1 |
| page_size | int32 | 否 | 每页数量，默认20 |

**响应**:

```json
{
    "messages": [
        {
            "message_id": 3001,
            "conversation_id": 2001,
            "sender_id": 1002,
            "sender_name": "朋友1",
            "sender_avatar": "https://example.com/avatar1.jpg",
            "content_type": 1,
            "content": "这是一条包含关键词的消息",
            "extra_data": null,
            "status": 2,
            "created_at": 1705838400
        }
    ],
    "total": 15,
    "page": 1,
    "page_size": 20
}
```

**示例**:

```bash
# 搜索所有会话
curl -X GET "http://localhost:8888/api/v1/messages/search?keyword=关键词" \
  -H "Authorization: Bearer {token}"

# 搜索指定会话
curl -X GET "http://localhost:8888/api/v1/messages/search?keyword=关键词&conversation_id=2001" \
  -H "Authorization: Bearer {token}"
```

---

## 六、文件相关接口

### 6.1 上传文件

**接口**: `POST /api/v1/files/upload`

**描述**: 上传文件（头像、图片、语音、文件）

**请求 Header**:

```
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | 文件 |
| file_type | string | 是 | 文件类型（image/audio/file） |

**响应**:

```json
{
    "file_id": 4001,
    "file_url": "https://example.com/files/abc123.jpg",
    "file_hash": "sha256_hash_value"
}
```

**示例**:

```bash
curl -X POST http://localhost:8888/api/v1/files/upload \
  -H "Authorization: Bearer {token}" \
  -F "file=@/path/to/image.jpg" \
  -F "file_type=image"
```

---

### 6.2 批量上传文件

**接口**: `POST /api/v1/files/batch-upload`

**描述**: 批量上传多个文件

**请求 Header**:

```
Authorization: Bearer {token}
Content-Type: application/json
```

**请求参数**:

```json
{
    "files": [
        {
            "file_name": "image1.jpg",
            "file_data": "base64_encoded_data",
            "file_type": "image"
        },
        {
            "file_name": "image2.jpg",
            "file_data": "base64_encoded_data",
            "file_type": "image"
        }
    ]
}
```

**响应**:

```json
{
    "files": [
        {
            "file_id": 4001,
            "file_url": "https://example.com/files/abc123.jpg",
            "file_hash": "sha256_hash_value1"
        },
        {
            "file_id": 4002,
            "file_url": "https://example.com/files/def456.jpg",
            "file_hash": "sha256_hash_value2"
        }
    ]
}
```

---

### 6.3 下载文件

**接口**: `GET /api/v1/files/:file_id`

**描述**: 下载指定文件

**请求 Header**:

```
Authorization: Bearer {token}
```

**响应**: 文件流

**示例**:

```bash
curl -X GET http://localhost:8888/api/v1/files/4001 \
  -H "Authorization: Bearer {token}" \
  -o downloaded_file.jpg
```

---

## 七、WebSocket 接口

### 7.1 建立连接

**接口**: `ws://localhost:8888/api/v1/ws?token={jwt_token}`

**描述**: 建立 WebSocket 长连接，用于实时消息推送

**连接参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| token | string | 是 | JWT Token |

**示例**:

```javascript
const ws = new WebSocket('ws://localhost:8888/api/v1/ws?token=' + jwt_token);

ws.onopen = function() {
    console.log('WebSocket 连接成功');
};

ws.onmessage = function(event) {
    const message = JSON.parse(event.data);
    console.log('收到消息:', message);
};

ws.onclose = function() {
    console.log('WebSocket 连接关闭');
};
```

---

### 7.2 消息格式

所有 WebSocket 消息统一格式：

```json
{
    "type": "message_type",
    "data": { ... },
    "msg_id": "unique_message_id",
    "time": 1705838400
}
```

---

### 7.3 发送消息（客户端 -> 服务器）

**消息类型**: `send_message`

**数据格式**:

```json
{
    "type": "send_message",
    "data": {
        "conversation_id": 2001,
        "content_type": 1,  // 1: 文本, 2: 图片, 3: 语音, 4: 文件
        "content": "你好",
        "extra_data": null
    },
    "msg_id": "client_generated_id",
    "time": 1705838400
}
```

**响应**:

```json
{
    "type": "send_message_ack",
    "data": {
        "message_id": 3001,
        "created_at": 1705838400
    },
    "msg_id": "client_generated_id",
    "time": 1705838400
}
```

---

### 7.4 接收新消息（服务器 -> 客户端）

**消息类型**: `new_message`

**数据格式**:

```json
{
    "type": "new_message",
    "data": {
        "message": {
            "message_id": 3001,
            "conversation_id": 2001,
            "sender_id": 1002,
            "sender_name": "朋友1",
            "sender_avatar": "https://example.com/avatar1.jpg",
            "content_type": 1,
            "content": "你好",
            "extra_data": null,
            "status": 1,
            "created_at": 1705838400
        }
    },
    "msg_id": "server_generated_id",
    "time": 1705838400
}
```

---

### 7.5 好友申请通知（服务器 -> 客户端）

**消息类型**: `friend_request`

**数据格式**:

```json
{
    "type": "friend_request",
    "data": {
        "request": {
            "request_id": 5001,
            "from_user_id": 1003,
            "from_user": {
                "user_id": 1003,
                "username": "friend1",
                "nickname": "朋友1",
                "avatar": "https://example.com/avatar1.jpg"
            },
            "message": "你好，我想加你为好友",
            "status": 0,
            "created_at": 1705838400
        }
    },
    "msg_id": "server_generated_id",
    "time": 1705838400
}
```

---

### 7.6 心跳检测（双向）

**客户端 -> 服务器**:

```json
{
    "type": "heartbeat",
    "data": {
        "ping": "ping"
    },
    "msg_id": "client_generated_id",
    "time": 1705838400
}
```

**服务器 -> 客户端**:

```json
{
    "type": "heartbeat",
    "data": {
        "pong": "pong"
    },
    "msg_id": "server_generated_id",
    "time": 1705838400
}
```

**心跳间隔**: 30秒

---

## 八、接口测试

### 8.1 Postman Collection

完整的 Postman Collection 位于：`/opt/Beehive/docs/api/Beehive.postman_collection.json`

### 8.2 测试流程

1. **注册用户**
   ```
   POST /api/v1/auth/send-code (发送验证码)
   POST /api/v1/auth/register (注册)
   ```

2. **登录获取 Token**
   ```
   POST /api/v1/auth/login
   ```

3. **添加好友**
   ```
   POST /api/v1/friends/request (发送申请)
   POST /api/v1/friends/request/handle (处理申请)
   GET /api/v1/friends (查看好友列表)
   ```

4. **创建会话**
   ```
   POST /api/v1/conversations (创建单聊/群聊)
   GET /api/v1/conversations (查看会话列表)
   ```

5. **发送消息**
   ```
   建立 WebSocket 连接
   发送 send_message 消息
   ```

6. **搜索消息**
   ```
   GET /api/v1/messages/search?keyword=关键词
   ```

---

## 九、常见问题

### Q1: Token 过期怎么办？

A: Token 默认 7 天过期，需要重新登录获取新 Token。后续会增加 Refresh Token 机制。

### Q2: WebSocket 断线重连怎么处理？

A: 客户端检测到断线后，使用指数退避策略重连，重连成功后同步离线消息。

### Q3: 如何实现消息撤回？

A: 后续版本将支持消息撤回功能（2分钟内）。

### Q4: 文件上传大小限制？

A: 默认限制 10MB，可在配置文件中修改。

### Q5: 如何实现群聊 @ 某人？

A: 后续版本将支持 @ 功能，客户端在消息内容中添加 @user_id 标记。

---

## 十、版本历史

| 版本 | 日期 | 说明 |
|------|------|------|
| v1.0.0 | 2026-01-21 | 初始版本，支持基本的 IM 功能 |

---

## 十一、联系方式

- **项目地址**: https://github.com/HappyLadySauce/Beehive
- **文档地址**: /opt/Beehive/docs/
- **作者邮箱**: 13452552349@163.com
