-- 群公告：conversations 表增加 announcement 字段
ALTER TABLE conversations ADD COLUMN IF NOT EXISTS announcement TEXT NOT NULL DEFAULT '';
