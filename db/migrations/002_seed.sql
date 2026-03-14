-- 测试用户与角色（接口测试等用）
-- testuser 1000000001/ password123；test1/test2/test3 账号密码同用户名，bcrypt 哈希由 gen_bcrypt.go 生成

-- 角色 user
INSERT INTO roles (id, name, description)
VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'user', 'normal user')
ON CONFLICT (id) DO NOTHING;

-- testuser（密码 password123）
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

-- test1 / test1
INSERT INTO users (id, username, password_hash, status)
VALUES (
    '1000000002',
    'test1',
    '$2a$10$n.vxOt3RHiLJx03lVnKHwOpx9IkNb1ESLNFFcw6ZjuvBr86vN.9Hu',
    'normal'
)
ON CONFLICT (id) DO UPDATE SET password_hash = EXCLUDED.password_hash;
INSERT INTO user_roles (user_id, role_id)
VALUES ('1000000002', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11')
ON CONFLICT (user_id, role_id) DO NOTHING;
INSERT INTO user_profiles (user_id, nickname, avatar_url, bio, status)
VALUES ('1000000002', 'Test1', '', '', 'normal')
ON CONFLICT (user_id) DO NOTHING;

-- test2 / test2
INSERT INTO users (id, username, password_hash, status)
VALUES (
    '1000000003',
    'test2',
    '$2a$10$gCwCkHVeYYiaviDzk8LoVONPDUxvlzyLuLhUxowDUgoBc1RSJb8VC',
    'normal'
)
ON CONFLICT (id) DO UPDATE SET password_hash = EXCLUDED.password_hash;
INSERT INTO user_roles (user_id, role_id)
VALUES ('1000000003', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11')
ON CONFLICT (user_id, role_id) DO NOTHING;
INSERT INTO user_profiles (user_id, nickname, avatar_url, bio, status)
VALUES ('1000000003', 'Test2', '', '', 'normal')
ON CONFLICT (user_id) DO NOTHING;

-- test3 / test3
INSERT INTO users (id, username, password_hash, status)
VALUES (
    '1000000004',
    'test3',
    '$2a$10$xfpMPuUNqD3J/bWzIFqQVewTYKoWb8OfjtnV6b0J5QT.7YljWIFFm',
    'normal'
)
ON CONFLICT (id) DO UPDATE SET password_hash = EXCLUDED.password_hash;
INSERT INTO user_roles (user_id, role_id)
VALUES ('1000000004', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11')
ON CONFLICT (user_id, role_id) DO NOTHING;
INSERT INTO user_profiles (user_id, nickname, avatar_url, bio, status)
VALUES ('1000000004', 'Test3', '', '', 'normal')
ON CONFLICT (user_id) DO NOTHING;

-- 群聊：test1 与 test2 共同所在群，群号 11 位 10000000001
INSERT INTO conversations (id, type, name, created_at, last_active_at)
VALUES ('10000000001', 'group', '测试群', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO conversation_members (id, conversation_id, user_id, role, status, joined_at)
VALUES
    (gen_random_uuid(), '10000000001', '1000000002', 'owner', 'active', NOW()),
    (gen_random_uuid(), '10000000001', '1000000003', 'member', 'active', NOW())
ON CONFLICT (conversation_id, user_id) DO NOTHING;

-- 序列推进到 seed 之后，避免新注册与内置用户冲突
SELECT setval('user_id_seq', 1000000005);
