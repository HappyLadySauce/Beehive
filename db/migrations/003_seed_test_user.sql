-- 接口测试用测试用户（test/interface 默认 TEST_USER=testuser, TEST_PASSWORD=password123）
-- 执行顺序：001、002 后再执行本文件。password_hash 由 db/scripts/gen_bcrypt.go 生成。

-- 角色 user（若不存在）
INSERT INTO roles (id, name, description)
VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'user', 'normal user')
ON CONFLICT (id) DO NOTHING;

-- 测试用户（密码 password123 的 bcrypt 哈希）
INSERT INTO users (id, username, password_hash, status)
VALUES (
    'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'testuser',
    '$2a$10$SMzAjyVPlzmRkxBb6NjnH.rTyhgMGFrfkKM5fhshqeQ8F4kUSAgb.',
    'normal'
)
ON CONFLICT (id) DO NOTHING;

-- 绑定角色
INSERT INTO user_roles (user_id, role_id)
VALUES ('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11')
ON CONFLICT (user_id, role_id) DO NOTHING;

-- 用户资料（供 UserService.GetUser / user.me 返回）
INSERT INTO user_profiles (user_id, nickname, avatar_url, bio, status)
VALUES (
    'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'Test User',
    '',
    'Interface test user',
    'normal'
)
ON CONFLICT (user_id) DO NOTHING;
