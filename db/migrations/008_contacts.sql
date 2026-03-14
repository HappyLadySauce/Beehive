-- 联系人表：支持添加/列表/移除联系人（owner 与 contact_user 均为 10 位用户 ID）
-- 执行顺序：需在 006 之后执行（users.id 已为 VARCHAR(10)）

CREATE TABLE IF NOT EXISTS contacts (
    owner_id       VARCHAR(10) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    contact_user_id VARCHAR(10) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status         TEXT NOT NULL DEFAULT 'accepted',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (owner_id, contact_user_id),
    CONSTRAINT chk_contact_not_self CHECK (owner_id != contact_user_id)
);

CREATE INDEX IF NOT EXISTS idx_contacts_owner_id ON contacts (owner_id);
