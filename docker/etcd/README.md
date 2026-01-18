# etcd 集群开发环境

这是一个用于本地开发的 etcd 3 节点集群 Docker Compose 配置。

## 快速开始

### 启动集群

```bash
cd docker/etcd
docker-compose up -d
```

### 停止集群

```bash
docker-compose down
```

### 停止并删除数据

```bash
docker-compose down -v
```

## 配置说明

### 端口映射

- **etcd-0**: 
  - Client: `localhost:2379`
  - Peer: `localhost:2380`
- **etcd-1**: 
  - Client: `localhost:23790`
  - Peer: `localhost:23800`
- **etcd-2**: 
  - Client: `localhost:23791`
  - Peer: `localhost:23801`

### 认证信息

默认认证信息（可在 `.env` 文件中修改）：
- **用户名**: `Beehive`
- **密码**: `Beehive`

### 环境变量

创建 `.env` 文件（可参考 `etcd.env`）来自定义配置：

```bash
cp etcd.env .env
# 编辑 .env 文件修改认证信息
```

## 连接方式

### 使用 etcdctl 连接

连接到主节点（etcd-0）：

```bash
docker exec -it etcd-0 etcdctl --endpoints=http://localhost:2379 \
  --user=Beehive:Beehive \
  endpoint health
```

连接到集群（所有节点）：

```bash
docker exec -it etcd-0 etcdctl --endpoints=http://etcd-0:2379,http://etcd-1:2379,http://etcd-2:2379 \
  --user=Beehive:Beehive \
  endpoint health
```

### 从主机连接

```bash
# 需要先安装 etcdctl
etcdctl --endpoints=http://localhost:2379,http://localhost:23790,http://localhost:23791 \
  --user=Beehive:Beehive \
  endpoint health
```

### 在应用中使用

在应用配置中使用以下 endpoints：

```yaml
etcd:
  endpoints: ["localhost:2379", "localhost:23790", "localhost:23791"]
  username: "Beehive"
  password: "Beehive"
  prefix: "/beehive/services"
```

## 健康检查

### 检查集群状态

```bash
docker-compose ps
```

### 查看日志

```bash
# 查看所有节点日志
docker-compose logs -f

# 查看特定节点日志
docker-compose logs -f etcd-0
```

### 检查集群成员

```bash
docker exec -it etcd-0 etcdctl --endpoints=http://localhost:2379 \
  --user=Beehive:Beehive \
  member list
```

## 数据持久化

数据存储在 Docker volumes 中：
- `etcd-0-data`
- `etcd-1-data`
- `etcd-2-data`

查看 volumes：

```bash
docker volume ls | grep etcd
```

删除数据（谨慎操作）：

```bash
docker-compose down -v
```

## 常见操作

### 查看集群信息

```bash
docker exec -it etcd-0 etcdctl --endpoints=http://localhost:2379 \
  --user=Beehive:Beehive \
  endpoint status
```

### 设置键值

```bash
docker exec -it etcd-0 etcdctl --endpoints=http://localhost:2379 \
  --user=Beehive:Beehive \
  put /test/key "test value"
```

### 获取键值

```bash
docker exec -it etcd-0 etcdctl --endpoints=http://localhost:2379 \
  --user=Beehive:Beehive \
  get /test/key
```

### 列出所有键

```bash
docker exec -it etcd-0 etcdctl --endpoints=http://localhost:2379 \
  --user=Beehive:Beehive \
  get --prefix /
```

## 故障排查

### 节点无法启动

1. 检查端口是否被占用：
   ```bash
   netstat -tuln | grep -E '2379|2380|23790|23800|23791|23801'
   ```

2. 查看节点日志：
   ```bash
   docker-compose logs etcd-0
   ```

### 认证失败

如果认证失败，可以重新初始化：

```bash
# 停止集群
docker-compose down -v

# 重新启动
docker-compose up -d

# 等待初始化完成（查看 etcd-init 容器日志）
docker-compose logs etcd-init
```

### 集群状态异常

检查集群健康状态：

```bash
docker exec -it etcd-0 etcdctl --endpoints=http://etcd-0:2379,http://etcd-1:2379,http://etcd-2:2379 \
  --user=Beehive:Beehive \
  endpoint health
```

## 注意事项

1. **生产环境**: 此配置仅适用于开发环境，生产环境需要更严格的安全配置
2. **数据备份**: 定期备份 volumes 数据
3. **资源限制**: 可以根据需要添加资源限制（CPU、内存）
4. **网络安全**: 生产环境应使用 TLS 加密

## 版本信息

- etcd 版本: v3.5.9
- Docker Compose 版本: 3.8+
