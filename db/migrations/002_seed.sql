-- 测试用户与角色（接口测试等用）
-- 测试用户 id 固定为 10 位：1000000001；密码 password123 的 bcrypt 哈希

-- 角色 user
INSERT INTO roles (id, name, description)
VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'user', 'normal user')
ON CONFLICT (id) DO NOTHING;

-- 测试用户（密码 password123）
INSERT INTO users (id, username, password_hash, status)
VALUES (
    '1000000001',
    'testuser',
    '$2a$10$Btjw1bDKXWA0jG7UUiw9OO51wzHFNBw.IvB.OoZl.35xYwlmBNYYq',
    'normal'
)
ON CONFLICT (id) DO UPDATE SET password_hash = EXCLUDED.password_hash;

INSERT INTO user_roles (user_id, role_id)
VALUES ('1000000001', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11')
ON CONFLICT (user_id, role_id) DO NOTHING;

INSERT INTO user_profiles (user_id, nickname, avatar_url, bio, status)
VALUES (
    '1000000001',
    'Test User',
    '',
    'Interface test user',
    'normal'
)
ON CONFLICT (user_id) DO NOTHING;

-- 序列推进到 seed 之后，避免新注册与 testuser(1000000001) 冲突
SELECT setval('user_id_seq', 1000000002);
