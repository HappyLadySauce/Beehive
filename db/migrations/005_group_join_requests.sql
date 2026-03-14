-- 入群申请表：群聊、申请人、留言、状态、审批人与时间
CREATE TABLE IF NOT EXISTS group_join_requests (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id VARCHAR(36) NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
    user_id         VARCHAR(10) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    message         TEXT NOT NULL DEFAULT '',
    status          TEXT NOT NULL DEFAULT 'pending',
    processed_at    TIMESTAMPTZ,
    processed_by    VARCHAR(10) REFERENCES users (id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (conversation_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_group_join_requests_conversation_status ON group_join_requests (conversation_id, status);
CREATE INDEX IF NOT EXISTS idx_group_join_requests_user_id ON group_join_requests (user_id);
