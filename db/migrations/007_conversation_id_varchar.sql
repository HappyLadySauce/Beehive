-- 将会话 ID 从 UUID 改为 VARCHAR(20)，以支持单聊 UUID 与群聊 11 位数字。
-- 现有数据保持原 UUID 字符串形式存入；新群聊由应用层生成 11 位 ID。

-- 1. 删除引用 conversations(id) 的外键
ALTER TABLE conversation_members DROP CONSTRAINT IF EXISTS conversation_members_conversation_id_fkey;
ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_conversation_id_fkey;

-- 2. 修改 conversations 表主键
ALTER TABLE conversations ADD COLUMN IF NOT EXISTS id_new VARCHAR(20);
UPDATE conversations SET id_new = id::text WHERE id_new IS NULL;
ALTER TABLE conversations DROP CONSTRAINT conversations_pkey;
ALTER TABLE conversations DROP COLUMN id;
ALTER TABLE conversations RENAME COLUMN id_new TO id;
ALTER TABLE conversations ADD PRIMARY KEY (id);
ALTER TABLE conversations ALTER COLUMN id SET NOT NULL;

-- 3. 修改 conversation_members.conversation_id
ALTER TABLE conversation_members ADD COLUMN IF NOT EXISTS conversation_id_new VARCHAR(20);
UPDATE conversation_members SET conversation_id_new = conversation_id::text WHERE conversation_id_new IS NULL;
ALTER TABLE conversation_members DROP CONSTRAINT IF EXISTS conversation_members_conversation_id_fkey;
ALTER TABLE conversation_members DROP COLUMN conversation_id;
ALTER TABLE conversation_members RENAME COLUMN conversation_id_new TO conversation_id;
ALTER TABLE conversation_members ALTER COLUMN conversation_id SET NOT NULL;
ALTER TABLE conversation_members ADD CONSTRAINT conversation_members_conversation_id_fkey FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE;

-- 4. 修改 messages.conversation_id
ALTER TABLE messages ADD COLUMN IF NOT EXISTS conversation_id_new VARCHAR(20);
UPDATE messages SET conversation_id_new = conversation_id::text WHERE conversation_id_new IS NULL;
ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_conversation_id_fkey;
ALTER TABLE messages DROP COLUMN conversation_id;
ALTER TABLE messages RENAME COLUMN conversation_id_new TO conversation_id;
ALTER TABLE messages ALTER COLUMN conversation_id SET NOT NULL;
ALTER TABLE messages ADD CONSTRAINT messages_conversation_id_fkey FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE;

-- 5. 若 conversation_read 表已存在（由 Message 服务 GORM 创建），则同步修改类型
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'conversation_read') THEN
        ALTER TABLE conversation_read ADD COLUMN IF NOT EXISTS conversation_id_new VARCHAR(20);
        UPDATE conversation_read SET conversation_id_new = conversation_id::text WHERE conversation_id_new IS NULL;
        ALTER TABLE conversation_read DROP CONSTRAINT IF EXISTS conversation_read_pkey;
        ALTER TABLE conversation_read DROP COLUMN conversation_id;
        ALTER TABLE conversation_read RENAME COLUMN conversation_id_new TO conversation_id;
        ALTER TABLE conversation_read ADD PRIMARY KEY (user_id, conversation_id);
        -- 若存在 FK 则重建（GORM 可能未建 FK）
        IF EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'conversation_read_conversation_id_fkey') THEN
            ALTER TABLE conversation_read DROP CONSTRAINT conversation_read_conversation_id_fkey;
        END IF;
        ALTER TABLE conversation_read ADD CONSTRAINT conversation_read_conversation_id_fkey FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE;
    END IF;
END $$;
