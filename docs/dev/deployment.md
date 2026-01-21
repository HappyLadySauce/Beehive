# Beehive IM 部署文档

## 一、部署概述

### 1.1 部署环境

- **开发环境**: Docker Compose（单机部署）
- **生产环境**: Kubernetes（集群部署）

### 1.2 系统要求

**硬件要求**:
- CPU: 4核以上
- 内存: 16GB 以上
- 磁盘: 100GB 以上

**软件要求**:
- Docker 20.10+
- Docker Compose 2.0+
- Kubernetes 1.24+ (生产环境)
- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- RabbitMQ 3.12+
- Elasticsearch 8.11+
- etcd 3.5+

---

## 二、开发环境部署（Docker Compose）

### 2.1 启动基础设施

**文件**: `/opt/Beehive/docker/docker-compose.yml`

```bash
# 进入 docker 目录
cd /opt/Beehive/docker

# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f [service_name]
```

**服务列表**:
- postgres (5432)
- redis (6379)
- rabbitmq (5672, 15672)
- elasticsearch (9200)
- kibana (5601)
- etcd (2379)

### 2.2 初始化数据库

```bash
# 创建数据库
docker exec -it beehive-postgres psql -U postgres -c "CREATE DATABASE beehive;"

# 执行初始化脚本
docker exec -i beehive-postgres psql -U postgres -d beehive < scripts/init_db.sql
```

### 2.3 初始化 Elasticsearch

```bash
# 安装 IK 分词器
docker exec -it beehive-elasticsearch elasticsearch-plugin install https://github.com/medcl/elasticsearch-analysis-ik/releases/download/v8.11.0/elasticsearch-analysis-ik-8.11.0.zip

# 重启 Elasticsearch
docker restart beehive-elasticsearch

# 创建索引
./scripts/init_es.sh
```

### 2.4 初始化 RabbitMQ

```bash
./scripts/init_rabbitmq.sh
```

### 2.5 启动 RPC 服务

```bash
# 生成代码（首次）
./scripts/gen_code.sh

# 启动 User RPC
cd rpc/user
go run user.go -f etc/user.yaml

# 启动 Friend RPC
cd rpc/friend
go run friend.go -f etc/friend.yaml

# 启动 Chat RPC
cd rpc/chat
go run chat.go -f etc/chat.yaml

# 启动 Message RPC
cd rpc/message
go run message.go -f etc/message.yaml

# 启动 File RPC
cd rpc/file
go run file.go -f etc/file.yaml

# 启动 Search RPC
cd rpc/search
go run search.go -f etc/search.yaml
```

### 2.6 启动 API Gateway

```bash
cd api/gateway
go run gateway.go -f etc/gateway.yaml
```

### 2.7 测试

```bash
# 健康检查
curl http://localhost:8888/ping

# 发送验证码
curl -X POST http://localhost:8888/api/v1/auth/send-code \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","purpose":"register"}'

# 注册
curl -X POST http://localhost:8888/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123","code":"123456"}'
```

---

## 三、生产环境部署（Kubernetes）

### 3.1 准备工作

**前提条件**:
- 已有 Kubernetes 集群
- 已安装 kubectl
- 已配置 kubeconfig

**创建命名空间**:

```bash
kubectl create namespace beehive
```

### 3.2 配置 ConfigMap

**文件**: `k8s/configmap.yaml`

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: beehive-config
  namespace: beehive
data:
  # PostgreSQL
  POSTGRES_HOST: beehive-postgres
  POSTGRES_PORT: "5432"
  POSTGRES_DB: beehive
  POSTGRES_USER: postgres
  
  # Redis
  REDIS_HOST: beehive-redis
  REDIS_PORT: "6379"
  
  # RabbitMQ
  RABBITMQ_HOST: beehive-rabbitmq
  RABBITMQ_PORT: "5672"
  
  # Elasticsearch
  ELASTICSEARCH_HOST: beehive-elasticsearch
  ELASTICSEARCH_PORT: "9200"
  
  # etcd
  ETCD_ENDPOINTS: beehive-etcd-0:2379,beehive-etcd-1:2379,beehive-etcd-2:2379
```

### 3.3 配置 Secret

**文件**: `k8s/secret.yaml`

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: beehive-secret
  namespace: beehive
type: Opaque
data:
  # Base64 编码
  postgres-password: cG9zdGdyZXM=  # postgres
  redis-password: ""
  rabbitmq-password: Z3Vlc3Q=  # guest
  jwt-secret: eW91ci1zZWNyZXQta2V5
```

### 3.4 部署 PostgreSQL

**文件**: `k8s/postgres.yaml`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: beehive-postgres
  namespace: beehive
spec:
  type: ClusterIP
  ports:
    - port: 5432
      targetPort: 5432
  selector:
    app: postgres
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: beehive-postgres
  namespace: beehive
spec:
  serviceName: beehive-postgres
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
        - name: postgres
          image: postgres:15-alpine
          ports:
            - containerPort: 5432
          env:
            - name: POSTGRES_DB
              valueFrom:
                configMapKeyRef:
                  name: beehive-config
                  key: POSTGRES_DB
            - name: POSTGRES_USER
              valueFrom:
                configMapKeyRef:
                  name: beehive-config
                  key: POSTGRES_USER
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: beehive-secret
                  key: postgres-password
          volumeMounts:
            - name: postgres-data
              mountPath: /var/lib/postgresql/data
  volumeClaimTemplates:
    - metadata:
        name: postgres-data
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 50Gi
```

### 3.5 部署 Redis

**文件**: `k8s/redis.yaml`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: beehive-redis
  namespace: beehive
spec:
  type: ClusterIP
  ports:
    - port: 6379
      targetPort: 6379
  selector:
    app: redis
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: beehive-redis
  namespace: beehive
spec:
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
        - name: redis
          image: redis:7-alpine
          ports:
            - containerPort: 6379
          volumeMounts:
            - name: redis-data
              mountPath: /data
      volumes:
        - name: redis-data
          persistentVolumeClaim:
            claimName: redis-pvc
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: redis-pvc
  namespace: beehive
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```

### 3.6 部署 RabbitMQ

**文件**: `k8s/rabbitmq.yaml`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: beehive-rabbitmq
  namespace: beehive
spec:
  type: ClusterIP
  ports:
    - port: 5672
      targetPort: 5672
      name: amqp
    - port: 15672
      targetPort: 15672
      name: management
  selector:
    app: rabbitmq
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: beehive-rabbitmq
  namespace: beehive
spec:
  serviceName: beehive-rabbitmq
  replicas: 3
  selector:
    matchLabels:
      app: rabbitmq
  template:
    metadata:
      labels:
        app: rabbitmq
    spec:
      containers:
        - name: rabbitmq
          image: rabbitmq:3-management-alpine
          ports:
            - containerPort: 5672
            - containerPort: 15672
          env:
            - name: RABBITMQ_DEFAULT_USER
              value: guest
            - name: RABBITMQ_DEFAULT_PASS
              valueFrom:
                secretKeyRef:
                  name: beehive-secret
                  key: rabbitmq-password
          volumeMounts:
            - name: rabbitmq-data
              mountPath: /var/lib/rabbitmq
  volumeClaimTemplates:
    - metadata:
        name: rabbitmq-data
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 20Gi
```

### 3.7 部署 Elasticsearch

**文件**: `k8s/elasticsearch.yaml`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: beehive-elasticsearch
  namespace: beehive
spec:
  type: ClusterIP
  ports:
    - port: 9200
      targetPort: 9200
      name: http
    - port: 9300
      targetPort: 9300
      name: transport
  selector:
    app: elasticsearch
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: beehive-elasticsearch
  namespace: beehive
spec:
  serviceName: beehive-elasticsearch
  replicas: 3
  selector:
    matchLabels:
      app: elasticsearch
  template:
    metadata:
      labels:
        app: elasticsearch
    spec:
      initContainers:
        - name: increase-vm-max-map
          image: busybox
          command: ["sysctl", "-w", "vm.max_map_count=262144"]
          securityContext:
            privileged: true
      containers:
        - name: elasticsearch
          image: docker.elastic.co/elasticsearch/elasticsearch:8.11.0
          ports:
            - containerPort: 9200
            - containerPort: 9300
          env:
            - name: cluster.name
              value: beehive-es
            - name: discovery.seed_hosts
              value: beehive-elasticsearch-0,beehive-elasticsearch-1,beehive-elasticsearch-2
            - name: cluster.initial_master_nodes
              value: beehive-elasticsearch-0,beehive-elasticsearch-1,beehive-elasticsearch-2
            - name: ES_JAVA_OPTS
              value: "-Xms2g -Xmx2g"
            - name: xpack.security.enabled
              value: "false"
          resources:
            requests:
              memory: "4Gi"
              cpu: "1"
            limits:
              memory: "4Gi"
              cpu: "2"
          volumeMounts:
            - name: es-data
              mountPath: /usr/share/elasticsearch/data
  volumeClaimTemplates:
    - metadata:
        name: es-data
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 100Gi
```

### 3.8 部署 etcd

**文件**: `k8s/etcd.yaml`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: beehive-etcd
  namespace: beehive
spec:
  clusterIP: None
  ports:
    - port: 2379
      name: client
    - port: 2380
      name: peer
  selector:
    app: etcd
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: beehive-etcd
  namespace: beehive
spec:
  serviceName: beehive-etcd
  replicas: 3
  selector:
    matchLabels:
      app: etcd
  template:
    metadata:
      labels:
        app: etcd
    spec:
      containers:
        - name: etcd
          image: quay.io/coreos/etcd:v3.5.11
          ports:
            - containerPort: 2379
            - containerPort: 2380
          env:
            - name: ETCD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: ETCD_INITIAL_CLUSTER
              value: beehive-etcd-0=http://beehive-etcd-0.beehive-etcd:2380,beehive-etcd-1=http://beehive-etcd-1.beehive-etcd:2380,beehive-etcd-2=http://beehive-etcd-2.beehive-etcd:2380
            - name: ETCD_INITIAL_CLUSTER_STATE
              value: new
            - name: ETCD_INITIAL_CLUSTER_TOKEN
              value: beehive-etcd-cluster
            - name: ETCD_LISTEN_CLIENT_URLS
              value: http://0.0.0.0:2379
            - name: ETCD_ADVERTISE_CLIENT_URLS
              value: http://$(ETCD_NAME).beehive-etcd:2379
            - name: ETCD_LISTEN_PEER_URLS
              value: http://0.0.0.0:2380
            - name: ETCD_INITIAL_ADVERTISE_PEER_URLS
              value: http://$(ETCD_NAME).beehive-etcd:2380
          volumeMounts:
            - name: etcd-data
              mountPath: /etcd-data
  volumeClaimTemplates:
    - metadata:
        name: etcd-data
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 10Gi
```

### 3.9 部署 RPC 服务

**文件**: `k8s/user-rpc.yaml`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: beehive-user-rpc
  namespace: beehive
spec:
  type: ClusterIP
  ports:
    - port: 8001
      targetPort: 8001
  selector:
    app: user-rpc
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: beehive-user-rpc
  namespace: beehive
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-rpc
  template:
    metadata:
      labels:
        app: user-rpc
    spec:
      containers:
        - name: user-rpc
          image: beehive/user-rpc:latest
          ports:
            - containerPort: 8001
          env:
            - name: DB_HOST
              valueFrom:
                configMapKeyRef:
                  name: beehive-config
                  key: POSTGRES_HOST
            - name: REDIS_HOST
              valueFrom:
                configMapKeyRef:
                  name: beehive-config
                  key: REDIS_HOST
            - name: ETCD_ENDPOINTS
              valueFrom:
                configMapKeyRef:
                  name: beehive-config
                  key: ETCD_ENDPOINTS
          resources:
            requests:
              memory: "512Mi"
              cpu: "500m"
            limits:
              memory: "1Gi"
              cpu: "1"
```

类似地部署其他 RPC 服务（Friend、Chat、Message、File、Search）。

### 3.10 部署 API Gateway

**文件**: `k8s/gateway.yaml`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: beehive-gateway
  namespace: beehive
spec:
  type: LoadBalancer
  ports:
    - port: 80
      targetPort: 8888
  selector:
    app: gateway
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: beehive-gateway
  namespace: beehive
spec:
  replicas: 3
  selector:
    matchLabels:
      app: gateway
  template:
    metadata:
      labels:
        app: gateway
    spec:
      containers:
        - name: gateway
          image: beehive/gateway:latest
          ports:
            - containerPort: 8888
          env:
            - name: ETCD_ENDPOINTS
              valueFrom:
                configMapKeyRef:
                  name: beehive-config
                  key: ETCD_ENDPOINTS
            - name: JWT_SECRET
              valueFrom:
                secretKeyRef:
                  name: beehive-secret
                  key: jwt-secret
          resources:
            requests:
              memory: "512Mi"
              cpu: "500m"
            limits:
              memory: "1Gi"
              cpu: "1"
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: beehive-gateway-hpa
  namespace: beehive
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: beehive-gateway
  minReplicas: 3
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
```

### 3.11 应用部署

```bash
# 应用所有配置
kubectl apply -f k8s/

# 查看 Pod 状态
kubectl get pods -n beehive

# 查看 Service
kubectl get svc -n beehive

# 查看日志
kubectl logs -f deployment/beehive-gateway -n beehive
```

---

## 四、Docker 镜像构建

### 4.1 User RPC Dockerfile

**文件**: `rpc/user/Dockerfile`

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -o user-rpc rpc/user/user.go

# 运行镜像
FROM alpine:latest

WORKDIR /app

# 安装 ca-certificates
RUN apk --no-cache add ca-certificates

# 从 builder 复制二进制文件
COPY --from=builder /app/user-rpc .
COPY --from=builder /app/rpc/user/etc etc/

EXPOSE 8001

CMD ["./user-rpc", "-f", "etc/user.yaml"]
```

### 4.2 构建镜像

```bash
# User RPC
docker build -t beehive/user-rpc:latest -f rpc/user/Dockerfile .

# Friend RPC
docker build -t beehive/friend-rpc:latest -f rpc/friend/Dockerfile .

# Chat RPC
docker build -t beehive/chat-rpc:latest -f rpc/chat/Dockerfile .

# Message RPC
docker build -t beehive/message-rpc:latest -f rpc/message/Dockerfile .

# File RPC
docker build -t beehive/file-rpc:latest -f rpc/file/Dockerfile .

# Search RPC
docker build -t beehive/search-rpc:latest -f rpc/search/Dockerfile .

# Gateway
docker build -t beehive/gateway:latest -f api/gateway/Dockerfile .
```

### 4.3 推送镜像

```bash
# 登录 Docker Hub
docker login

# 推送镜像
docker push beehive/user-rpc:latest
docker push beehive/friend-rpc:latest
docker push beehive/chat-rpc:latest
docker push beehive/message-rpc:latest
docker push beehive/file-rpc:latest
docker push beehive/search-rpc:latest
docker push beehive/gateway:latest
```

---

## 五、监控和日志

### 5.1 Prometheus + Grafana

**安装 Prometheus**:

```bash
kubectl apply -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/main/bundle.yaml
```

**配置 ServiceMonitor**:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: beehive-monitor
  namespace: beehive
spec:
  selector:
    matchLabels:
      app: gateway
  endpoints:
    - port: metrics
      path: /metrics
```

### 5.2 ELK 日志收集

**安装 Filebeat**:

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: filebeat
  namespace: beehive
spec:
  selector:
    matchLabels:
      app: filebeat
  template:
    metadata:
      labels:
        app: filebeat
    spec:
      containers:
        - name: filebeat
          image: docker.elastic.co/beats/filebeat:8.11.0
          volumeMounts:
            - name: config
              mountPath: /usr/share/filebeat/filebeat.yml
              subPath: filebeat.yml
            - name: varlibdockercontainers
              mountPath: /var/lib/docker/containers
              readOnly: true
      volumes:
        - name: config
          configMap:
            name: filebeat-config
        - name: varlibdockercontainers
          hostPath:
            path: /var/lib/docker/containers
```

---

## 六、备份和恢复

### 6.1 数据库备份

```bash
# 定时备份
0 2 * * * kubectl exec -it beehive-postgres-0 -n beehive -- pg_dump -U postgres beehive > backup_$(date +\%Y\%m\%d).sql
```

### 6.2 Elasticsearch 备份

```bash
# 创建快照
kubectl exec -it beehive-elasticsearch-0 -n beehive -- \
  curl -X PUT "localhost:9200/_snapshot/beehive_backup/snapshot_$(date +\%Y\%m\%d)"
```

---

## 七、常见运维操作

### 7.1 扩容

```bash
# 扩容 Gateway
kubectl scale deployment beehive-gateway --replicas=5 -n beehive

# 扩容 User RPC
kubectl scale deployment beehive-user-rpc --replicas=5 -n beehive
```

### 7.2 滚动更新

```bash
# 更新镜像
kubectl set image deployment/beehive-gateway gateway=beehive/gateway:v1.1.0 -n beehive

# 查看更新状态
kubectl rollout status deployment/beehive-gateway -n beehive

# 回滚
kubectl rollout undo deployment/beehive-gateway -n beehive
```

### 7.3 故障排查

```bash
# 查看 Pod 日志
kubectl logs -f pod/beehive-gateway-xxx -n beehive

# 进入 Pod
kubectl exec -it pod/beehive-gateway-xxx -n beehive -- /bin/sh

# 查看 Pod 详情
kubectl describe pod beehive-gateway-xxx -n beehive

# 查看事件
kubectl get events -n beehive
```

---

## 八、安全加固

### 8.1 网络策略

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: beehive-network-policy
  namespace: beehive
spec:
  podSelector:
    matchLabels:
      app: user-rpc
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: gateway
      ports:
        - protocol: TCP
          port: 8001
  egress:
    - to:
        - podSelector:
            matchLabels:
              app: postgres
      ports:
        - protocol: TCP
          port: 5432
```

### 8.2 RBAC

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: beehive-sa
  namespace: beehive
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: beehive-role
  namespace: beehive
rules:
  - apiGroups: [""]
    resources: ["pods", "services"]
    verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: beehive-rolebinding
  namespace: beehive
subjects:
  - kind: ServiceAccount
    name: beehive-sa
roleRef:
  kind: Role
  name: beehive-role
  apiGroup: rbac.authorization.k8s.io
```

---

## 九、性能调优

### 9.1 资源限制

- **CPU Requests**: 实际使用量的 70%
- **CPU Limits**: Requests 的 2倍
- **Memory Requests**: 实际使用量的 80%
- **Memory Limits**: Requests 的 1.5倍

### 9.2 连接池配置

- **数据库连接池**: MaxOpenConns=20, MaxIdleConns=10
- **Redis 连接池**: PoolSize=50
- **gRPC 连接池**: PoolSize=5

---

## 十、总结

本文档详细描述了 Beehive IM 系统的部署方案，包括开发环境和生产环境。

**部署建议**:
- 开发环境使用 Docker Compose
- 生产环境使用 Kubernetes
- 做好监控和告警
- 定期备份数据
- 制定灾难恢复计划
