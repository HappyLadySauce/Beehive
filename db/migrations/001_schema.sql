-- Beehive 当前版本建表脚本（用户 ID 10 位、会话 ID varchar(20)、联系人表）
-- 依赖顺序：users → user_profiles, user_roles → conversations → conversation_members, messages, conversation_read, contacts

-- users: 账号与认证，id 为 10 位数字字符串
CREATE TABLE IF NOT EXISTS users (
    id            VARCHAR(10) PRIMARY KEY,
    username      TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    status        TEXT NOT NULL DEFAULT 'normal',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);

-- user_profiles: 用户展示信息
CREATE TABLE IF NOT EXISTS user_profiles (
    user_id    VARCHAR(10) PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    nickname   TEXT NOT NULL DEFAULT '',
    avatar_url TEXT NOT NULL DEFAULT '',
    bio        TEXT NOT NULL DEFAULT '',
    status     TEXT NOT NULL DEFAULT 'normal',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- RBAC: 角色与权限
CREATE TABLE IF NOT EXISTS roles (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS permissions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code        TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id       UUID NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions (id) ON DELETE CASCADE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id)
);
CREATE TABLE IF NOT EXISTS user_roles (
    user_id    VARCHAR(10) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role_id    UUID NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles (user_id);

-- conversations: 单聊 id 为 UUID 字符串，群聊 id 为 11 位数字字符串
CREATE TABLE IF NOT EXISTS conversations (
    id            VARCHAR(20) PRIMARY KEY,
    type          TEXT NOT NULL DEFAULT 'single',
    name          TEXT NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_active_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_conversations_type ON conversations (type);
CREATE INDEX IF NOT EXISTS idx_conversations_last_active_at ON conversations (last_active_at DESC);

-- conversation_members
CREATE TABLE IF NOT EXISTS conversation_members (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id VARCHAR(20) NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
    user_id         VARCHAR(10) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role            TEXT NOT NULL DEFAULT 'member',
    status          TEXT NOT NULL DEFAULT 'active',
    joined_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (conversation_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_conversation_members_conversation_id ON conversation_members (conversation_id);
CREATE INDEX IF NOT EXISTS idx_conversation_members_user_id ON conversation_members (user_id);

-- messages
CREATE TABLE IF NOT EXISTS messages (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    server_msg_id    TEXT NOT NULL UNIQUE,
    client_msg_id    TEXT NOT NULL DEFAULT '',
    conversation_id  VARCHAR(20) NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
    from_user_id     VARCHAR(10) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    to_user_id       VARCHAR(10) REFERENCES users (id) ON DELETE SET NULL,
    body_type        TEXT NOT NULL DEFAULT 'text',
    body_text        TEXT NOT NULL DEFAULT '',
    server_time      BIGINT NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_messages_conversation_server_time ON messages (conversation_id, server_time DESC);
CREATE INDEX IF NOT EXISTS idx_messages_server_msg_id ON messages (server_msg_id);

-- conversation_read: Message 服务 GORM 可能自动建表，此处显式创建以保持一致
CREATE TABLE IF NOT EXISTS conversation_read (
    user_id               VARCHAR(10) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    conversation_id       VARCHAR(20) NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
    last_read_server_time BIGINT NOT NULL DEFAULT 0,
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, conversation_id)
);

-- contacts
CREATE TABLE IF NOT EXISTS contacts (
    owner_id        VARCHAR(10) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    contact_user_id VARCHAR(10) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status          TEXT NOT NULL DEFAULT 'accepted',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (owner_id, contact_user_id),
    CONSTRAINT chk_contact_not_self CHECK (owner_id != contact_user_id)
);
CREATE INDEX IF NOT EXISTS idx_contacts_owner_id ON contacts (owner_id);
