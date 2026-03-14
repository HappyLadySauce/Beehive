-- messages: 点对点/群聊消息
CREATE TABLE IF NOT EXISTS messages (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    server_msg_id    TEXT NOT NULL UNIQUE,
    client_msg_id    TEXT NOT NULL DEFAULT '',
    conversation_id  UUID NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
    from_user_id     UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    to_user_id       UUID REFERENCES users (id) ON DELETE SET NULL,  -- 单聊时使用，群聊可为 NULL
    body_type        TEXT NOT NULL DEFAULT 'text',  -- text | image | system | ...
    body_text        TEXT NOT NULL DEFAULT '',
    server_time      BIGINT NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_messages_conversation_server_time ON messages (conversation_id, server_time DESC);
CREATE INDEX IF NOT EXISTS idx_messages_server_msg_id ON messages (server_msg_id);
