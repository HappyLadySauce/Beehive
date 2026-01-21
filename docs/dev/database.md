# Beehive IM 数据库设计文档

## 一、数据库选型

**数据库**: PostgreSQL 15

**选型原因**:
- 成熟稳定的开源关系型数据库
- 支持 JSONB 类型，灵活存储扩展字段
- 强大的全文检索能力（虽然我们用 Elasticsearch）
- 支持分区表，适合消息表按月分区
- 优秀的并发性能和事务支持
- 丰富的索引类型（B-tree、Hash、GiST、GIN）

## 二、数据库设计原则

1. **范式化设计**: 遵循第三范式，减少数据冗余
2. **合理索引**: 为高频查询字段添加索引，避免过度索引
3. **字段约束**: 使用 NOT NULL、UNIQUE、CHECK 等约束保证数据完整性
4. **外键约束**: 适当使用外键约束，但不过度依赖（考虑性能）
5. **时间戳**: 所有表都有 `created_at` 和 `updated_at` (如需)
6. **软删除**: 重要数据使用软删除（`deleted_at`）
7. **分区策略**: 大表（如消息表）按时间分区

## 三、表结构设计

### 3.1 用户表 (users)

**表说明**: 存储用户基本信息

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    nickname VARCHAR(50),
    avatar_url VARCHAR(255),
    status SMALLINT DEFAULT 1,  -- 1:正常 2:禁用
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP,
    CONSTRAINT chk_username_length CHECK (char_length(username) >= 3),
    CONSTRAINT chk_email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}$')
);

-- 索引
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_status ON users(status);

-- 注释
COMMENT ON TABLE users IS '用户表';
COMMENT ON COLUMN users.id IS '用户ID';
COMMENT ON COLUMN users.username IS '用户名，唯一';
COMMENT ON COLUMN users.email IS '邮箱，唯一';
COMMENT ON COLUMN users.password_hash IS '密码哈希（bcrypt）';
COMMENT ON COLUMN users.nickname IS '昵称';
COMMENT ON COLUMN users.avatar_url IS '头像URL';
COMMENT ON COLUMN users.status IS '状态: 1-正常 2-禁用';
COMMENT ON COLUMN users.email_verified IS '邮箱是否已验证';
COMMENT ON COLUMN users.created_at IS '创建时间';
COMMENT ON COLUMN users.updated_at IS '更新时间';
COMMENT ON COLUMN users.last_login_at IS '最后登录时间';
```

**字段说明**:
- `id`: 主键，自增
- `username`: 用户名，3-50字符，唯一
- `email`: 邮箱，唯一，用于登录和找回密码
- `password_hash`: bcrypt 加密后的密码
- `nickname`: 昵称，可为空，默认为用户名
- `avatar_url`: 头像 URL
- `status`: 用户状态，1=正常，2=禁用
- `email_verified`: 邮箱是否已验证
- `created_at`: 注册时间
- `updated_at`: 最后更新时间
- `last_login_at`: 最后登录时间

**索引说明**:
- `idx_users_username`: 用户名查询
- `idx_users_email`: 邮箱查询
- `idx_users_created_at`: 按注册时间排序
- `idx_users_status`: 查询正常/禁用用户

---

### 3.2 邮箱验证码表 (email_verification_codes)

**表说明**: 存储邮箱验证码，用于注册、重置密码等

```sql
CREATE TABLE email_verification_codes (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(100) NOT NULL,
    code VARCHAR(6) NOT NULL,
    purpose VARCHAR(20) NOT NULL,  -- register, reset_password, bind_email
    expired_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX idx_email_code ON email_verification_codes(email, code, purpose);
CREATE INDEX idx_expired_at ON email_verification_codes(expired_at);
CREATE INDEX idx_created_at ON email_verification_codes(created_at);

-- 注释
COMMENT ON TABLE email_verification_codes IS '邮箱验证码表';
COMMENT ON COLUMN email_verification_codes.email IS '邮箱地址';
COMMENT ON COLUMN email_verification_codes.code IS '6位数字验证码';
COMMENT ON COLUMN email_verification_codes.purpose IS '用途: register, reset_password, bind_email';
COMMENT ON COLUMN email_verification_codes.expired_at IS '过期时间';
COMMENT ON COLUMN email_verification_codes.used IS '是否已使用';
COMMENT ON COLUMN email_verification_codes.created_at IS '创建时间';
```

**字段说明**:
- `email`: 邮箱地址
- `code`: 6位数字验证码
- `purpose`: 用途，`register` / `reset_password` / `bind_email`
- `expired_at`: 过期时间，通常为 5 分钟后
- `used`: 是否已使用，验证后标记为 true
- `created_at`: 创建时间

**清理策略**:
- 定时任务清理过期验证码（created_at < now() - interval '1 day'）

---

### 3.3 好友关系表 (friends)

**表说明**: 存储好友关系（双向）

```sql
CREATE TABLE friends (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    friend_id BIGINT NOT NULL,
    remark VARCHAR(50),  -- 备注名
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_friend FOREIGN KEY (friend_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_not_self CHECK (user_id != friend_id),
    UNIQUE(user_id, friend_id)
);

-- 索引
CREATE INDEX idx_friends_user_id ON friends(user_id);
CREATE INDEX idx_friends_friend_id ON friends(friend_id);
CREATE INDEX idx_friends_created_at ON friends(created_at);

-- 注释
COMMENT ON TABLE friends IS '好友关系表';
COMMENT ON COLUMN friends.user_id IS '用户ID';
COMMENT ON COLUMN friends.friend_id IS '好友ID';
COMMENT ON COLUMN friends.remark IS '备注名';
COMMENT ON COLUMN friends.created_at IS '添加时间';
```

**字段说明**:
- `user_id`: 用户 ID
- `friend_id`: 好友 ID
- `remark`: 备注名，可自定义好友显示名称
- `created_at`: 添加好友时间

**约束说明**:
- `UNIQUE(user_id, friend_id)`: 防止重复添加
- `CHECK (user_id != friend_id)`: 不能添加自己为好友
- 外键级联删除：用户删除后，好友关系自动删除

**业务逻辑**:
- 好友关系是双向的，A 添加 B 为好友时，需要插入两条记录：
  - (user_id=A, friend_id=B)
  - (user_id=B, friend_id=A)

---

### 3.4 好友申请表 (friend_requests)

**表说明**: 存储好友申请记录

```sql
CREATE TABLE friend_requests (
    id BIGSERIAL PRIMARY KEY,
    from_user_id BIGINT NOT NULL,
    to_user_id BIGINT NOT NULL,
    message TEXT,
    status SMALLINT DEFAULT 0,  -- 0:待处理 1:已同意 2:已拒绝 3:已过期
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    handled_at TIMESTAMP,
    CONSTRAINT fk_from_user FOREIGN KEY (from_user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_to_user FOREIGN KEY (to_user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_not_self_request CHECK (from_user_id != to_user_id)
);

-- 索引
CREATE INDEX idx_friend_requests_to_user ON friend_requests(to_user_id, status);
CREATE INDEX idx_friend_requests_from_user ON friend_requests(from_user_id);
CREATE INDEX idx_friend_requests_created_at ON friend_requests(created_at);

-- 注释
COMMENT ON TABLE friend_requests IS '好友申请表';
COMMENT ON COLUMN friend_requests.from_user_id IS '申请人ID';
COMMENT ON COLUMN friend_requests.to_user_id IS '被申请人ID';
COMMENT ON COLUMN friend_requests.message IS '申请消息';
COMMENT ON COLUMN friend_requests.status IS '状态: 0-待处理 1-已同意 2-已拒绝 3-已过期';
COMMENT ON COLUMN friend_requests.created_at IS '申请时间';
COMMENT ON COLUMN friend_requests.handled_at IS '处理时间';
```

**字段说明**:
- `from_user_id`: 发起申请的用户 ID
- `to_user_id`: 接收申请的用户 ID
- `message`: 申请附加消息，如"你好，我想加你为好友"
- `status`: 申请状态
  - 0: 待处理
  - 1: 已同意
  - 2: 已拒绝
  - 3: 已过期（7天后自动过期）
- `created_at`: 申请时间
- `handled_at`: 处理时间（同意或拒绝）

**索引说明**:
- `idx_friend_requests_to_user`: 查询某用户收到的待处理申请
- `idx_friend_requests_from_user`: 查询某用户发出的申请
- `idx_friend_requests_created_at`: 按时间排序

---

### 3.5 会话表 (conversations)

**表说明**: 存储会话（单聊/群聊）

```sql
CREATE TABLE conversations (
    id BIGSERIAL PRIMARY KEY,
    type SMALLINT NOT NULL,  -- 1:单聊 2:群聊
    name VARCHAR(100),       -- 群聊名称
    avatar VARCHAR(255),     -- 群头像
    owner_id BIGINT,         -- 群主ID
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_owner FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 索引
CREATE INDEX idx_conversations_type ON conversations(type);
CREATE INDEX idx_conversations_owner ON conversations(owner_id);
CREATE INDEX idx_conversations_created_at ON conversations(created_at);

-- 注释
COMMENT ON TABLE conversations IS '会话表';
COMMENT ON COLUMN conversations.type IS '类型: 1-单聊 2-群聊';
COMMENT ON COLUMN conversations.name IS '群聊名称';
COMMENT ON COLUMN conversations.avatar IS '群头像URL';
COMMENT ON COLUMN conversations.owner_id IS '群主ID';
COMMENT ON COLUMN conversations.created_at IS '创建时间';
COMMENT ON COLUMN conversations.updated_at IS '更新时间';
```

**字段说明**:
- `type`: 会话类型，1=单聊，2=群聊
- `name`: 群聊名称（单聊为空）
- `avatar`: 群聊头像（单聊为空）
- `owner_id`: 群主 ID（仅群聊）
- `created_at`: 创建时间
- `updated_at`: 最后更新时间（群信息修改、最后一条消息时间）

**业务逻辑**:
- 单聊：两个用户之间只能有一个单聊会话
- 群聊：可以创建多个群聊
- 群主退出群聊：群主转让或解散群聊

---

### 3.6 会话成员表 (conversation_members)

**表说明**: 存储会话成员关系

```sql
CREATE TABLE conversation_members (
    id BIGSERIAL PRIMARY KEY,
    conversation_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    role SMALLINT DEFAULT 1,  -- 1:普通成员 2:管理员 3:群主
    unread_count INT DEFAULT 0,
    last_read_at TIMESTAMP,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_conversation FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    CONSTRAINT fk_member_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(conversation_id, user_id)
);

-- 索引
CREATE INDEX idx_conversation_members_conv ON conversation_members(conversation_id);
CREATE INDEX idx_conversation_members_user ON conversation_members(user_id);
CREATE INDEX idx_conversation_members_user_unread ON conversation_members(user_id, unread_count);

-- 注释
COMMENT ON TABLE conversation_members IS '会话成员表';
COMMENT ON COLUMN conversation_members.conversation_id IS '会话ID';
COMMENT ON COLUMN conversation_members.user_id IS '用户ID';
COMMENT ON COLUMN conversation_members.role IS '角色: 1-普通成员 2-管理员 3-群主';
COMMENT ON COLUMN conversation_members.unread_count IS '未读消息数';
COMMENT ON COLUMN conversation_members.last_read_at IS '最后阅读时间';
COMMENT ON COLUMN conversation_members.joined_at IS '加入时间';
```

**字段说明**:
- `conversation_id`: 会话 ID
- `user_id`: 用户 ID
- `role`: 成员角色
  - 1: 普通成员
  - 2: 管理员
  - 3: 群主
- `unread_count`: 未读消息数（每次收到新消息 +1，标记已读后清零）
- `last_read_at`: 最后阅读消息的时间
- `joined_at`: 加入会话的时间

**索引说明**:
- `idx_conversation_members_conv`: 查询会话所有成员
- `idx_conversation_members_user`: 查询用户所有会话
- `idx_conversation_members_user_unread`: 查询用户有未读消息的会话

---

### 3.7 消息表 (messages)

**表说明**: 存储所有消息（建议分区）

```sql
CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    conversation_id BIGINT NOT NULL,
    sender_id BIGINT NOT NULL,
    content_type SMALLINT NOT NULL,  -- 1:文本 2:图片 3:语音 4:文件
    content TEXT NOT NULL,           -- 文本内容或文件URL
    extra_data JSONB,                -- 扩展数据（图片尺寸、语音时长等）
    status SMALLINT DEFAULT 0,       -- 0:发送中 1:已送达 2:已读
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_msg_conversation FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    CONSTRAINT fk_msg_sender FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE
) PARTITION BY RANGE (created_at);

-- 创建分区（按月）
CREATE TABLE messages_2026_01 PARTITION OF messages
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');

CREATE TABLE messages_2026_02 PARTITION OF messages
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');

-- 索引
CREATE INDEX idx_messages_conversation ON messages(conversation_id, created_at DESC);
CREATE INDEX idx_messages_sender ON messages(sender_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);

-- 注释
COMMENT ON TABLE messages IS '消息表（分区表）';
COMMENT ON COLUMN messages.conversation_id IS '会话ID';
COMMENT ON COLUMN messages.sender_id IS '发送者ID';
COMMENT ON COLUMN messages.content_type IS '内容类型: 1-文本 2-图片 3-语音 4-文件';
COMMENT ON COLUMN messages.content IS '消息内容或文件URL';
COMMENT ON COLUMN messages.extra_data IS '扩展数据（JSON格式）';
COMMENT ON COLUMN messages.status IS '状态: 0-发送中 1-已送达 2-已读';
COMMENT ON COLUMN messages.created_at IS '发送时间';
```

**字段说明**:
- `conversation_id`: 所属会话 ID
- `sender_id`: 发送者 ID
- `content_type`: 消息类型
  - 1: 文本消息
  - 2: 图片消息
  - 3: 语音消息
  - 4: 文件消息
- `content`: 
  - 文本消息：消息文本内容
  - 图片/语音/文件：文件 URL
- `extra_data`: 扩展数据（JSONB 格式）
  - 图片：`{"width": 1920, "height": 1080, "size": 204800}`
  - 语音：`{"duration": 10, "size": 102400}`
  - 文件：`{"filename": "doc.pdf", "size": 1048576}`
- `status`: 消息状态
  - 0: 发送中
  - 1: 已送达
  - 2: 已读
- `created_at`: 发送时间

**分区说明**:
- 消息表数据量大，建议按月分区
- 每月自动创建新分区
- 旧分区可以归档到冷存储

**索引说明**:
- `idx_messages_conversation`: 查询会话消息（按时间倒序）
- `idx_messages_sender`: 查询用户发送的消息
- `idx_messages_created_at`: 按时间范围查询

**extra_data 示例**:

```json
// 图片消息
{
    "width": 1920,
    "height": 1080,
    "size": 204800,
    "thumbnail_url": "https://example.com/thumb.jpg"
}

// 语音消息
{
    "duration": 10,
    "size": 102400,
    "format": "mp3"
}

// 文件消息
{
    "filename": "document.pdf",
    "size": 1048576,
    "mime_type": "application/pdf"
}
```

---

### 3.8 文件表 (files)

**表说明**: 存储文件信息（支持去重）

```sql
CREATE TABLE files (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    file_hash VARCHAR(64) UNIQUE NOT NULL,  -- SHA256哈希
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    file_type VARCHAR(50) NOT NULL,         -- MIME类型: image/jpeg, audio/mp3等
    storage_path VARCHAR(512) NOT NULL,     -- 存储路径
    url VARCHAR(512) NOT NULL,              -- 访问URL
    usage_count INT DEFAULT 1,              -- 引用计数
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_file_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 索引
CREATE INDEX idx_files_hash ON files(file_hash);
CREATE INDEX idx_files_user ON files(user_id);
CREATE INDEX idx_files_created_at ON files(created_at);
CREATE INDEX idx_files_usage_count ON files(usage_count);

-- 注释
COMMENT ON TABLE files IS '文件表';
COMMENT ON COLUMN files.user_id IS '上传用户ID';
COMMENT ON COLUMN files.file_hash IS '文件SHA256哈希，用于去重';
COMMENT ON COLUMN files.file_name IS '原始文件名';
COMMENT ON COLUMN files.file_size IS '文件大小（字节）';
COMMENT ON COLUMN files.file_type IS 'MIME类型';
COMMENT ON COLUMN files.storage_path IS '存储路径';
COMMENT ON COLUMN files.url IS '访问URL';
COMMENT ON COLUMN files.usage_count IS '引用计数';
COMMENT ON COLUMN files.created_at IS '上传时间';
```

**字段说明**:
- `user_id`: 上传文件的用户 ID
- `file_hash`: 文件 SHA256 哈希，用于去重
- `file_name`: 原始文件名
- `file_size`: 文件大小（字节）
- `file_type`: MIME 类型，如 `image/jpeg`、`audio/mp3`、`video/mp4`
- `storage_path`: 文件存储路径（服务器本地或 OSS）
- `url`: 文件访问 URL
- `usage_count`: 引用计数，多个用户上传相同文件时 +1
- `created_at`: 首次上传时间

**去重逻辑**:
1. 客户端上传文件前计算 SHA256 哈希
2. 调用接口检查该哈希是否已存在
3. 如果存在，直接返回已有文件 URL，`usage_count` +1
4. 如果不存在，保存文件并插入记录

**清理策略**:
- 当 `usage_count` 为 0 时，可以删除文件

---

## 四、ER 图

```
┌──────────────┐          ┌──────────────────┐
│    users     │◄─────────│email_verification│
│              │          │      _codes      │
└──────┬───────┘          └──────────────────┘
       │
       │ 1:N
       │
       ▼
┌──────────────┐          ┌──────────────────┐
│   friends    │          │  friend_requests │
│              │          │                  │
└──────────────┘          └──────────────────┘
       │                           │
       │                           │
       │ 1:N                       │ 1:N
       │                           │
       ▼                           ▼
┌──────────────┐          ┌──────────────────┐
│conversation_ │◄─────────│  conversations   │
│   members    │  N:1     │                  │
└──────┬───────┘          └─────────┬────────┘
       │                            │
       │ 1:N                        │ 1:N
       │                            │
       └────────────┬───────────────┘
                    │
                    ▼
            ┌──────────────┐
            │   messages   │
            │              │
            └──────────────┘
                    │
                    │ N:1
                    │
                    ▼
            ┌──────────────┐
            │    files     │
            │              │
            └──────────────┘
```

## 五、索引策略

### 5.1 索引原则

1. **高频查询字段**: WHERE、JOIN、ORDER BY 字段
2. **唯一约束**: 添加 UNIQUE 索引
3. **外键字段**: 添加普通索引
4. **组合索引**: 多条件查询使用组合索引
5. **避免过度索引**: 索引会影响写入性能

### 5.2 索引列表

| 表名 | 索引名 | 字段 | 类型 | 说明 |
|------|--------|------|------|------|
| users | PRIMARY | id | B-tree | 主键 |
| users | idx_users_username | username | B-tree | 用户名查询 |
| users | idx_users_email | email | B-tree | 邮箱查询 |
| users | idx_users_created_at | created_at | B-tree | 按时间排序 |
| friends | PRIMARY | id | B-tree | 主键 |
| friends | idx_friends_user_id | user_id | B-tree | 查询用户好友 |
| friends | idx_friends_friend_id | friend_id | B-tree | 反向查询 |
| friend_requests | idx_friend_requests_to_user | to_user_id, status | B-tree | 查询待处理申请 |
| conversations | PRIMARY | id | B-tree | 主键 |
| conversation_members | idx_conversation_members_conv | conversation_id | B-tree | 查询会话成员 |
| conversation_members | idx_conversation_members_user | user_id | B-tree | 查询用户会话 |
| messages | idx_messages_conversation | conversation_id, created_at | B-tree | 查询会话消息 |
| messages | idx_messages_created_at | created_at | B-tree | 按时间查询 |
| files | idx_files_hash | file_hash | B-tree | 文件去重 |

## 六、数据库初始化脚本

完整的初始化脚本位于：`/opt/Beehive/scripts/init_db.sql`

执行方式：

```bash
# 创建数据库
psql -U postgres -c "CREATE DATABASE beehive;"

# 执行初始化脚本
psql -U postgres -d beehive -f scripts/init_db.sql
```

## 七、缓存策略

### 7.1 Redis 缓存设计

| 缓存键 | 数据类型 | TTL | 说明 |
|--------|---------|-----|------|
| `user:info:{user_id}` | Hash | 1小时 | 用户信息 |
| `user:online:{user_id}` | String | 5分钟 | 用户在线状态 |
| `user:session:{user_id}` | Hash | 7天 | 用户会话信息 |
| `conv:list:{user_id}` | List | 30分钟 | 会话列表 |
| `conv:detail:{conv_id}` | Hash | 1小时 | 会话详情 |
| `conv:members:{conv_id}` | Set | 1小时 | 会话成员列表 |
| `friend:list:{user_id}` | Set | 1小时 | 好友列表 |
| `msg:hot:{conv_id}` | List | 10分钟 | 热点消息（最近100条） |
| `file:info:{file_hash}` | Hash | 永久 | 文件信息（去重） |

### 7.2 缓存更新策略

- **Cache Aside**: 先查缓存，缓存未命中再查数据库，写入时先更新数据库再删除缓存
- **Write Through**: 写入时同时更新数据库和缓存
- **延迟双删**: 更新数据库后删除缓存，延迟500ms再次删除

## 八、数据备份策略

### 8.1 全量备份

- **频率**: 每天凌晨 2:00
- **保留时间**: 30 天
- **工具**: `pg_dump`

```bash
pg_dump -U postgres -d beehive -F c -f backup_$(date +%Y%m%d).dump
```

### 8.2 增量备份

- **频率**: 每小时
- **保留时间**: 7 天
- **工具**: PostgreSQL WAL 归档

### 8.3 灾难恢复

- **RTO**: 1小时
- **RPO**: 10分钟
- **主从复制**: 实时同步

## 九、性能优化

### 9.1 分区表

消息表按月分区：

```sql
-- 自动创建未来分区的函数
CREATE OR REPLACE FUNCTION create_messages_partition()
RETURNS void AS $$
DECLARE
    start_date DATE;
    end_date DATE;
    partition_name TEXT;
BEGIN
    start_date := date_trunc('month', CURRENT_DATE);
    end_date := start_date + interval '1 month';
    partition_name := 'messages_' || to_char(start_date, 'YYYY_MM');
    
    EXECUTE format('CREATE TABLE IF NOT EXISTS %I PARTITION OF messages FOR VALUES FROM (%L) TO (%L)',
                   partition_name, start_date, end_date);
END;
$$ LANGUAGE plpgsql;

-- 定时任务每月创建分区
SELECT cron.schedule('create-messages-partition', '0 0 1 * *', 'SELECT create_messages_partition()');
```

### 9.2 读写分离

- **主库**: 处理所有写操作
- **从库**: 处理所有读操作
- **复制延迟**: < 1秒

### 9.3 连接池

- **最大连接数**: 100
- **空闲连接数**: 10
- **连接超时**: 30秒

### 9.4 慢查询监控

```sql
-- 开启慢查询日志
ALTER DATABASE beehive SET log_min_duration_statement = 1000;  -- 1秒

-- 查看慢查询
SELECT * FROM pg_stat_statements ORDER BY total_time DESC LIMIT 10;
```

## 十、数据安全

### 10.1 密码加密

- 使用 bcrypt 算法
- Cost 因子: 12
- 不可逆加密

### 10.2 敏感数据加密

- 邮箱可见，但查询时脱敏
- 手机号中间4位脱敏

### 10.3 数据权限

- 用户只能访问自己的数据
- 好友可以查看对方信息
- 会话成员可以查看会话消息

## 十一、数据库维护

### 11.1 定期清理

- 清理过期验证码（每天）
- 清理已删除的文件（引用计数为0）
- 归档历史消息（6个月前）

### 11.2 VACUUM

```sql
-- 定期 VACUUM
VACUUM ANALYZE messages;
VACUUM ANALYZE conversation_members;
```

### 11.3 索引维护

```sql
-- 重建索引
REINDEX TABLE messages;
```

## 十二、数据库版本管理

使用数据库迁移工具管理表结构变更：

- **工具**: golang-migrate 或 goose
- **迁移文件**: `migrations/000001_init_schema.up.sql`
- **版本控制**: Git

## 十三、FAQ

### Q1: 为什么选择 PostgreSQL 而不是 MySQL？

A: PostgreSQL 支持 JSONB、分区表、更强的并发性能，适合复杂业务场景。

### Q2: 消息表为什么要分区？

A: 消息数据量大，按月分区可以提高查询性能，方便历史数据归档。

### Q3: 好友关系为什么是双向的？

A: 为了查询效率，双向存储可以避免 JOIN 查询。

### Q4: 文件表的引用计数有什么用？

A: 用于去重，多个用户上传相同文件时只保存一份，节省存储空间。

### Q5: 为什么不使用 MongoDB？

A: IM 系统涉及复杂的关系查询（好友、会话、消息），关系型数据库更合适。
