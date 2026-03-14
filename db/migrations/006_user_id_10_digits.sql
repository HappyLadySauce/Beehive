-- 将 users.id 从 UUID 改为 10 位数字字符串，并更新所有引用表。
-- 执行前需已执行 001~005。若表为空或仅含测试数据，可为每用户分配唯一 10 位 ID。

-- 1. 创建映射表并填充：为每个现有 user 分配唯一 10 位 ID
CREATE TABLE IF NOT EXISTS _migrate_user_id_mapping (
    old_id UUID PRIMARY KEY,
    new_id VARCHAR(10) NOT NULL UNIQUE
);

DO $$
DECLARE
    r RECORD;
    new_id TEXT;
    done BOOLEAN;
BEGIN
    FOR r IN SELECT id FROM users
    LOOP
        done := FALSE;
        WHILE NOT done LOOP
            new_id := (1000000000 + floor(random() * 9000000000)::bigint)::text;
            BEGIN
                INSERT INTO _migrate_user_id_mapping (old_id, new_id) VALUES (r.id, new_id);
                done := TRUE;
            EXCEPTION WHEN unique_violation THEN
                NULL; -- 重试
            END;
        END LOOP;
    END LOOP;
END $$;

-- 2. 删除引用 users(id) 的外键
ALTER TABLE user_profiles DROP CONSTRAINT IF EXISTS user_profiles_user_id_fkey;
ALTER TABLE conversation_members DROP CONSTRAINT IF EXISTS conversation_members_user_id_fkey;
ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_from_user_id_fkey;
ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_to_user_id_fkey;
ALTER TABLE user_roles DROP CONSTRAINT IF EXISTS user_roles_user_id_fkey;
-- conversation_read 若存在（由 Message 服务 GORM 创建）
ALTER TABLE conversation_read DROP CONSTRAINT IF EXISTS conversation_read_user_id_fkey;
ALTER TABLE conversation_read DROP CONSTRAINT IF EXISTS fk_conversation_read_user;

-- 3. 修改 users 表主键
ALTER TABLE users ADD COLUMN IF NOT EXISTS id_new VARCHAR(10);
UPDATE users u SET id_new = m.new_id FROM _migrate_user_id_mapping m WHERE m.old_id = u.id;
ALTER TABLE users DROP CONSTRAINT users_pkey;
ALTER TABLE users DROP COLUMN id;
ALTER TABLE users RENAME COLUMN id_new TO id;
ALTER TABLE users ADD PRIMARY KEY (id);
ALTER TABLE users ALTER COLUMN id SET NOT NULL;

-- 4. 修改 user_profiles
ALTER TABLE user_profiles ADD COLUMN IF NOT EXISTS user_id_new VARCHAR(10);
UPDATE user_profiles up SET user_id_new = m.new_id FROM _migrate_user_id_mapping m WHERE m.old_id = up.user_id;
ALTER TABLE user_profiles DROP CONSTRAINT user_profiles_pkey;
ALTER TABLE user_profiles DROP COLUMN user_id;
ALTER TABLE user_profiles RENAME COLUMN user_id_new TO user_id;
ALTER TABLE user_profiles ADD PRIMARY KEY (user_id);
ALTER TABLE user_profiles ADD CONSTRAINT user_profiles_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- 5. 修改 conversation_members.user_id
ALTER TABLE conversation_members ADD COLUMN IF NOT EXISTS user_id_new VARCHAR(10);
UPDATE conversation_members cm SET user_id_new = m.new_id FROM _migrate_user_id_mapping m WHERE m.old_id = cm.user_id;
ALTER TABLE conversation_members DROP COLUMN user_id;
ALTER TABLE conversation_members RENAME COLUMN user_id_new TO user_id;
ALTER TABLE conversation_members ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE conversation_members ADD CONSTRAINT conversation_members_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- 6. 修改 messages
ALTER TABLE messages ADD COLUMN IF NOT EXISTS from_user_id_new VARCHAR(10);
ALTER TABLE messages ADD COLUMN IF NOT EXISTS to_user_id_new VARCHAR(10);
UPDATE messages msg SET from_user_id_new = m.new_id FROM _migrate_user_id_mapping m WHERE m.old_id = msg.from_user_id;
UPDATE messages msg SET to_user_id_new = m.new_id FROM _migrate_user_id_mapping m WHERE m.old_id = msg.to_user_id;
ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_from_user_id_fkey;
ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_to_user_id_fkey;
ALTER TABLE messages DROP COLUMN from_user_id;
ALTER TABLE messages DROP COLUMN to_user_id;
ALTER TABLE messages RENAME COLUMN from_user_id_new TO from_user_id;
ALTER TABLE messages RENAME COLUMN to_user_id_new TO to_user_id;
ALTER TABLE messages ALTER COLUMN from_user_id SET NOT NULL;
ALTER TABLE messages ADD CONSTRAINT messages_from_user_id_fkey FOREIGN KEY (from_user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE messages ADD CONSTRAINT messages_to_user_id_fkey FOREIGN KEY (to_user_id) REFERENCES users(id) ON DELETE SET NULL;

-- 7. 修改 user_roles（保留复合主键 user_id, role_id）
ALTER TABLE user_roles ADD COLUMN IF NOT EXISTS user_id_new VARCHAR(10);
UPDATE user_roles ur SET user_id_new = m.new_id FROM _migrate_user_id_mapping m WHERE m.old_id = ur.user_id;
ALTER TABLE user_roles DROP CONSTRAINT user_roles_pkey;
ALTER TABLE user_roles DROP COLUMN user_id;
ALTER TABLE user_roles RENAME COLUMN user_id_new TO user_id;
ALTER TABLE user_roles ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE user_roles ADD PRIMARY KEY (user_id, role_id);
ALTER TABLE user_roles ADD CONSTRAINT user_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles (user_id);

-- 8. 若存在 conversation_read，将其 user_id 改为 VARCHAR(10)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'conversation_read') THEN
        ALTER TABLE conversation_read ADD COLUMN IF NOT EXISTS user_id_new VARCHAR(10);
        UPDATE conversation_read cr SET user_id_new = m.new_id FROM _migrate_user_id_mapping m WHERE m.old_id = cr.user_id;
        ALTER TABLE conversation_read DROP CONSTRAINT IF EXISTS conversation_read_pkey;
        ALTER TABLE conversation_read DROP COLUMN user_id;
        ALTER TABLE conversation_read RENAME COLUMN user_id_new TO user_id;
        ALTER TABLE conversation_read ADD PRIMARY KEY (user_id, conversation_id);
        ALTER TABLE conversation_read ADD CONSTRAINT conversation_read_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;

-- 9. 删除映射表
DROP TABLE IF EXISTS _migrate_user_id_mapping;
