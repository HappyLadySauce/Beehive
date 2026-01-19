# Beehive IM Docker 部署文档

本目录包含 Beehive IM 系统的 Docker 容器配置文件。

## 快速启动

### 1. 启动所有基础设施服务

```bash
cd /opt/Beehive/docker
docker-compose up -d
```

这将启动以下服务：
- **PostgreSQL** (5432): 主数据库
- **Redis** (6379): 缓存和会话存储
- **RabbitMQ** (5672, 15672): 消息队列
- **Elasticsearch** (9200, 9300): 消息全文搜索
- **Kibana** (5601): Elasticsearch 数据可视化
- **etcd** (2379, 2380): 服务注册与发现

### 2. 查看服务状态

```bash
docker-compose ps
```

### 3. 查看服务日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f elasticsearch
docker-compose logs -f postgres
```

### 4. 停止所有服务

```bash
docker-compose down
```

### 5. 停止并删除所有数据

```bash
docker-compose down -v
```

## 服务说明

### PostgreSQL

- **端口**: 5432
- **用户名**: postgres
- **密码**: postgres
- **数据库**: beehive
- **数据卷**: postgres_data

**连接字符串**:
```
postgresql://postgres:postgres@localhost:5432/beehive?sslmode=disable
```

### Redis

- **端口**: 6379
- **密码**: 无
- **数据卷**: redis_data

**连接字符串**:
```
redis://localhost:6379/0
```

### RabbitMQ

- **AMQP 端口**: 5672
- **管理界面**: http://localhost:15672
- **用户名**: guest
- **密码**: guest
- **数据卷**: rabbitmq_data

**连接字符串**:
```
amqp://guest:guest@localhost:5672/
```

### Elasticsearch

- **HTTP 端口**: 9200
- **传输端口**: 9300
- **数据卷**: es_data
- **内存限制**: 512MB (可根据需要调整)

**连接地址**:
```
http://localhost:9200
```

**健康检查**:
```bash
curl http://localhost:9200/_cluster/health
```

**安装 IK 中文分词插件**:
```bash
docker exec -it beehive-elasticsearch bash
elasticsearch-plugin install https://github.com/medcl/elasticsearch-analysis-ik/releases/download/v8.11.0/elasticsearch-analysis-ik-8.11.0.zip
exit
docker restart beehive-elasticsearch
```

### Kibana

- **端口**: 5601
- **访问地址**: http://localhost:5601

Kibana 是 Elasticsearch 的可视化界面，可以用来：
- 查看索引数据
- 执行搜索查询
- 分析日志和指标
- 管理 Elasticsearch 集群

### etcd

- **客户端端口**: 2379
- **对等端口**: 2380
- **数据卷**: etcd_data

**连接地址**:
```
http://localhost:2379
```

**使用 etcdctl**:
```bash
# 进入容器
docker exec -it beehive-etcd sh

# 查看所有键
etcdctl get "" --prefix

# 查看服务注册信息
etcdctl get /beehive/services/ --prefix
```

## 初始化 Elasticsearch 索引

启动 Elasticsearch 后，需要创建消息索引：

```bash
curl -X PUT "http://localhost:9200/beehive-messages" -H 'Content-Type: application/json' -d'
{
  "settings": {
    "number_of_shards": 3,
    "number_of_replicas": 1,
    "analysis": {
      "analyzer": {
        "ik_max_word_analyzer": {
          "type": "custom",
          "tokenizer": "ik_max_word"
        },
        "ik_smart_analyzer": {
          "type": "custom",
          "tokenizer": "ik_smart"
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "message_id": { "type": "keyword" },
      "type": { "type": "keyword" },
      "from_id": { "type": "keyword" },
      "to_id": { "type": "keyword" },
      "group_id": { "type": "keyword" },
      "content": {
        "type": "text",
        "analyzer": "ik_max_word_analyzer",
        "search_analyzer": "ik_smart_analyzer",
        "fields": {
          "keyword": {
            "type": "keyword",
            "ignore_above": 256
          }
        }
      },
      "message_type": { "type": "keyword" },
      "status": { "type": "keyword" },
      "created_at": {
        "type": "date",
        "format": "epoch_second"
      },
      "updated_at": {
        "type": "date",
        "format": "epoch_second"
      }
    }
  }
}
'
```

验证索引创建成功：
```bash
curl http://localhost:9200/beehive-messages
```

## 初始化数据库

```bash
# 使用 psql 连接数据库
docker exec -it beehive-postgres psql -U postgres -d beehive

# 或者从宿主机连接
psql -h localhost -U postgres -d beehive
```

## 生产环境配置

### 1. PostgreSQL

```yaml
postgres:
  environment:
    POSTGRES_PASSWORD: <strong-password>
  volumes:
    - /data/postgres:/var/lib/postgresql/data
```

### 2. Redis

```yaml
redis:
  command: redis-server --requirepass <redis-password>
```

### 3. RabbitMQ

```yaml
rabbitmq:
  environment:
    RABBITMQ_DEFAULT_USER: <username>
    RABBITMQ_DEFAULT_PASS: <strong-password>
```

### 4. Elasticsearch

```yaml
elasticsearch:
  environment:
    - xpack.security.enabled=true
    - ELASTIC_PASSWORD=<strong-password>
    - "ES_JAVA_OPTS=-Xms2g -Xmx2g"  # 生产环境增加内存
```

### 5. etcd 集群

生产环境建议使用 etcd 集群（3节点或5节点），参考 `etcd/docker-compose.yml` 中的集群配置。

## 监控和维护

### 查看资源使用

```bash
docker stats
```

### 备份数据

```bash
# 备份 PostgreSQL
docker exec beehive-postgres pg_dump -U postgres beehive > backup.sql

# 备份 Elasticsearch
curl -X PUT "http://localhost:9200/_snapshot/my_backup" -H 'Content-Type: application/json' -d'
{
  "type": "fs",
  "settings": {
    "location": "/usr/share/elasticsearch/data/backup"
  }
}
'
```

### 清理日志

```bash
docker-compose logs --no-log-prefix > logs.txt
docker-compose restart
```

## 故障排查

### Elasticsearch 启动失败

**问题**: vm.max_map_count 太小

**解决方案**:
```bash
sudo sysctl -w vm.max_map_count=262144
# 永久生效
echo "vm.max_map_count=262144" | sudo tee -a /etc/sysctl.conf
```

### PostgreSQL 无法连接

**检查容器状态**:
```bash
docker logs beehive-postgres
```

**检查端口占用**:
```bash
netstat -tunlp | grep 5432
```

### RabbitMQ 管理界面无法访问

**等待服务完全启动** (约30秒):
```bash
docker logs beehive-rabbitmq
```

### Elasticsearch 内存不足

**调整内存限制**:
```yaml
elasticsearch:
  environment:
    - "ES_JAVA_OPTS=-Xms1g -Xmx1g"
```

## 参考资料

- [Docker Compose 文档](https://docs.docker.com/compose/)
- [PostgreSQL Docker 镜像](https://hub.docker.com/_/postgres)
- [Redis Docker 镜像](https://hub.docker.com/_/redis)
- [RabbitMQ Docker 镜像](https://hub.docker.com/_/rabbitmq)
- [Elasticsearch Docker 文档](https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html)
- [etcd Docker 文档](https://etcd.io/docs/latest/op-guide/container/)
