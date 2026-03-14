-- 好友申请表：发起方、接收方、状态（pending/accepted/declined）、可选留言
CREATE TABLE IF NOT EXISTS contact_requests (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_user_id VARCHAR(10) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    to_user_id   VARCHAR(10) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status       TEXT NOT NULL DEFAULT 'pending',
    message      TEXT NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_contact_request_not_self CHECK (from_user_id != to_user_id),
    UNIQUE (from_user_id, to_user_id)
);
CREATE INDEX IF NOT EXISTS idx_contact_requests_to_user_status ON contact_requests (to_user_id, status);
CREATE INDEX IF NOT EXISTS idx_contact_requests_from_user_id ON contact_requests (from_user_id);
