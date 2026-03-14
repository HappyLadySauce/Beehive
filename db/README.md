# 数据库迁移

## 当前脚本

- `migrations/001_schema.sql`：当前版本建表（用户 10 位 ID、会话 varchar(20)、联系人等）
- `migrations/002_seed.sql`：测试用户与角色（testuser / password123，id=1000000001）

## 重置并导入

在**仓库根目录**执行：

```bash
go run ./db/cmd/migrate
```

或指定 DSN：

```bash
go run ./db/cmd/migrate -dsn "postgres://user:pass@host:5432/dbname?sslmode=disable"
```

也可设置环境变量 `BEEHIVE_POSTGRES_DSN`，不传 `-dsn` 时优先使用该值；未设置时默认使用 `postgres://beehive:beehive@127.0.0.1:5432/beehive?sslmode=disable`。

程序会先执行 `DROP SCHEMA public CASCADE; CREATE SCHEMA public;` 清空库，再按文件名顺序执行 `migrations/*.sql`。
