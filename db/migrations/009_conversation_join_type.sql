-- 群加入方式：approval=需审批加入，direct=直接加入（申请即入群）
ALTER TABLE conversations ADD COLUMN IF NOT EXISTS join_type TEXT NOT NULL DEFAULT 'approval';
