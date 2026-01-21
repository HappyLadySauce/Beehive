# Beehive IM 消息队列设计文档

## 一、消息队列概述

### 1.1 技术选型

**消息队列**: RabbitMQ 3.x

**选型原因**:
- 成熟稳定的开源消息队列
- 支持多种消息模式（点对点、发布订阅、路由、主题）
- 高可用集群支持
- 消息持久化
- 完善的管理界面
- 支持消息确认机制
- Go 客户端库完善 (github.com/streadway/amqp)

### 1.2 使用场景

1. **消息异步处理**: 消息发送后异步持久化、推送、索引
2. **服务解耦**: Message RPC、Gateway、Search RPC 通过 MQ 解耦
3. **削峰填谷**: 高峰期消息积压在队列中，平滑处理
4. **可靠传输**: 消息确认机制保证消息不丢失

---

## 二、Exchange 和 Queue 设计

### 2.1 Exchange 配置

**Exchange 名称**: `beehive.message.exchange`

**Exchange 类型**: `topic`

**特性**:
- **Durable**: true（持久化，服务器重启后仍存在）
- **AutoDelete**: false（没有队列绑定时不自动删除）
- **Internal**: false（可以被生产者直接发布消息）

**Routing Key 规则**:
- `message.persist`: 消息持久化
- `message.push`: 消息推送
- `message.index`: 消息索引

### 2.2 Queue 配置

#### 2.2.1 持久化队列 (message.persist)

**队列名**: `beehive.message.persist`

**配置**:
```json
{
    "durable": true,
    "autoDelete": false,
    "exclusive": false,
    "arguments": {
        "x-message-ttl": 86400000,  // 消息TTL: 24小时（毫秒）
        "x-max-length": 1000000,    // 最大消息数: 100万
        "x-dead-letter-exchange": "beehive.dlx.exchange",
        "x-dead-letter-routing-key": "message.persist.dlx"
    }
}
```

**绑定**:
- Exchange: `beehive.message.exchange`
- Routing Key: `message.persist`

**消费者**: Message RPC Service（实际上消息已在发布前持久化，此队列可选）

**用途**: 消息持久化到数据库（备用）

---

#### 2.2.2 推送队列 (message.push)

**队列名**: `beehive.message.push`

**配置**:
```json
{
    "durable": true,
    "autoDelete": false,
    "exclusive": false,
    "arguments": {
        "x-message-ttl": 300000,    // 消息TTL: 5分钟（毫秒）
        "x-max-length": 100000,     // 最大消息数: 10万
        "x-dead-letter-exchange": "beehive.dlx.exchange",
        "x-dead-letter-routing-key": "message.push.dlx"
    }
}
```

**绑定**:
- Exchange: `beehive.message.exchange`
- Routing Key: `message.push`

**消费者**: Gateway Service (所有实例)

**用途**: 实时推送消息给在线用户（通过 WebSocket）

---

#### 2.2.3 索引队列 (message.index)

**队列名**: `beehive.message.index`

**配置**:
```json
{
    "durable": true,
    "autoDelete": false,
    "exclusive": false,
    "arguments": {
        "x-message-ttl": 3600000,   // 消息TTL: 1小时（毫秒）
        "x-max-length": 100000,     // 最大消息数: 10万
        "x-dead-letter-exchange": "beehive.dlx.exchange",
        "x-dead-letter-routing-key": "message.index.dlx"
    }
}
```

**绑定**:
- Exchange: `beehive.message.exchange`
- Routing Key: `message.index`

**消费者**: Search RPC Service

**用途**: 将消息索引到 Elasticsearch

---

### 2.3 死信队列 (DLX)

**Exchange 名称**: `beehive.dlx.exchange`

**Exchange 类型**: `direct`

**死信队列**:
- `beehive.message.persist.dlq` (Routing Key: `message.persist.dlx`)
- `beehive.message.push.dlq` (Routing Key: `message.push.dlx`)
- `beehive.message.index.dlq` (Routing Key: `message.index.dlx`)

**用途**: 处理失败的消息（消费失败、TTL 过期、队列满）

---

## 三、消息格式

### 3.1 消息结构

所有消息使用 JSON 格式：

```json
{
    "message_id": 12345,
    "conversation_id": 100,
    "sender_id": 1001,
    "receiver_ids": [1002, 1003],
    "content_type": 1,
    "content": "Hello",
    "extra_data": {
        "width": 1920,
        "height": 1080
    },
    "created_at": 1705838400
}
```

**字段说明**:
- `message_id`: 消息 ID（数据库主键）
- `conversation_id`: 会话 ID
- `sender_id`: 发送者 ID
- `receiver_ids`: 接收者 ID 列表（群聊有多个）
- `content_type`: 消息类型（1:文本 2:图片 3:语音 4:文件）
- `content`: 消息内容或文件 URL
- `extra_data`: 扩展数据（可选）
- `created_at`: 发送时间（Unix 时间戳）

### 3.2 消息属性

发布消息时设置的 AMQP 属性：

```go
amqp.Publishing{
    ContentType:  "application/json",
    DeliveryMode: amqp.Persistent,  // 持久化
    Priority:     0,
    Timestamp:    time.Now(),
    MessageId:    messageId,
}
```

---

## 四、生产者实现

### 4.1 初始化连接

```go
package mq

import (
    "github.com/streadway/amqp"
    "log"
)

type RabbitMQ struct {
    conn    *amqp.Connection
    channel *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }
    
    channel, err := conn.Channel()
    if err != nil {
        return nil, err
    }
    
    // 声明 Exchange
    err = channel.ExchangeDeclare(
        "beehive.message.exchange", // name
        "topic",                     // type
        true,                        // durable
        false,                       // auto-deleted
        false,                       // internal
        false,                       // no-wait
        nil,                         // arguments
    )
    if err != nil {
        return nil, err
    }
    
    return &RabbitMQ{
        conn:    conn,
        channel: channel,
    }, nil
}

func (r *RabbitMQ) Close() {
    r.channel.Close()
    r.conn.Close()
}
```

### 4.2 发布消息

```go
func (r *RabbitMQ) PublishMessage(msg *Message) error {
    body, err := json.Marshal(msg)
    if err != nil {
        return err
    }
    
    // 发布到3个队列
    routingKeys := []string{
        "message.persist",
        "message.push",
        "message.index",
    }
    
    for _, key := range routingKeys {
        err = r.channel.Publish(
            "beehive.message.exchange", // exchange
            key,                         // routing key
            false,                       // mandatory
            false,                       // immediate
            amqp.Publishing{
                ContentType:  "application/json",
                DeliveryMode: amqp.Persistent,
                Body:         body,
                MessageId:    fmt.Sprintf("%d", msg.MessageId),
                Timestamp:    time.Unix(msg.CreatedAt, 0),
            },
        )
        if err != nil {
            log.Printf("Failed to publish to %s: %v", key, err)
            // 可以考虑重试
        }
    }
    
    return nil
}
```

### 4.3 在 Message RPC 中使用

```go
// internal/svc/servicecontext.go
type ServiceContext struct {
    Config       config.Config
    MessageModel model.MessageModel
    RabbitMQ     *mq.RabbitMQ
}

func NewServiceContext(c config.Config) *ServiceContext {
    rabbitMQ, err := mq.NewRabbitMQ(c.RabbitMQ.Url)
    if err != nil {
        panic(err)
    }
    
    return &ServiceContext{
        Config:       c,
        MessageModel: model.NewMessageModel(sqlx.NewMysql(c.DataSource), c.Cache),
        RabbitMQ:     rabbitMQ,
    }
}

// internal/logic/sendmessagelogic.go
func (l *SendMessageLogic) SendMessage(req *message.SendMessageRequest) (*message.SendMessageResponse, error) {
    // 1. 保存消息到数据库
    result, err := l.svcCtx.MessageModel.Insert(l.ctx, &model.Message{
        ConversationId: req.ConversationId,
        SenderId:       req.SenderId,
        ContentType:    req.ContentType,
        Content:        req.Content,
    })
    if err != nil {
        return nil, err
    }
    
    messageId, _ := result.LastInsertId()
    
    // 2. 发布到 RabbitMQ
    err = l.svcCtx.RabbitMQ.PublishMessage(&mq.Message{
        MessageId:      messageId,
        ConversationId: req.ConversationId,
        SenderId:       req.SenderId,
        ContentType:    req.ContentType,
        Content:        req.Content,
        ExtraData:      req.ExtraData,
        CreatedAt:      time.Now().Unix(),
    })
    if err != nil {
        logx.Errorf("Failed to publish message to MQ: %v", err)
        // 不影响主流程，消息已保存到数据库
    }
    
    return &message.SendMessageResponse{
        MessageId: messageId,
        CreatedAt: time.Now().Unix(),
    }, nil
}
```

---

## 五、消费者实现

### 5.1 消费推送队列 (Gateway)

```go
// internal/consumer/messageConsumer.go
package consumer

import (
    "context"
    "encoding/json"
    "github.com/streadway/amqp"
    "github.com/zeromicro/go-zero/core/logx"
)

type MessageConsumer struct {
    channel *amqp.Channel
    wsManager *ws.Manager
}

func NewMessageConsumer(channel *amqp.Channel, wsManager *ws.Manager) *MessageConsumer {
    return &MessageConsumer{
        channel:   channel,
        wsManager: wsManager,
    }
}

func (c *MessageConsumer) Start(ctx context.Context) error {
    // 声明队列
    queue, err := c.channel.QueueDeclare(
        "beehive.message.push", // name
        true,                   // durable
        false,                  // delete when unused
        false,                  // exclusive
        false,                  // no-wait
        map[string]interface{}{
            "x-message-ttl":           300000,
            "x-max-length":            100000,
            "x-dead-letter-exchange":  "beehive.dlx.exchange",
            "x-dead-letter-routing-key": "message.push.dlx",
        },
    )
    if err != nil {
        return err
    }
    
    // 绑定队列到 Exchange
    err = c.channel.QueueBind(
        queue.Name,
        "message.push",
        "beehive.message.exchange",
        false,
        nil,
    )
    if err != nil {
        return err
    }
    
    // 设置 QoS（每次获取1条消息）
    err = c.channel.Qos(1, 0, false)
    if err != nil {
        return err
    }
    
    // 开始消费
    msgs, err := c.channel.Consume(
        queue.Name, // queue
        "",         // consumer
        false,      // auto-ack（手动确认）
        false,      // exclusive
        false,      // no-local
        false,      // no-wait
        nil,        // args
    )
    if err != nil {
        return err
    }
    
    logx.Info("MessageConsumer started")
    
    // 消费循环
    go func() {
        for {
            select {
            case <-ctx.Done():
                logx.Info("MessageConsumer stopped")
                return
            case msg, ok := <-msgs:
                if !ok {
                    return
                }
                c.handleMessage(msg)
            }
        }
    }()
    
    return nil
}

func (c *MessageConsumer) handleMessage(msg amqp.Delivery) {
    var message Message
    err := json.Unmarshal(msg.Body, &message)
    if err != nil {
        logx.Errorf("Failed to unmarshal message: %v", err)
        msg.Nack(false, false) // 拒绝并丢弃（或进入死信队列）
        return
    }
    
    // 推送给在线用户
    for _, receiverId := range message.ReceiverIds {
        err = c.wsManager.SendToUser(receiverId, &ws.Message{
            Type: "new_message",
            Data: message,
        })
        if err != nil {
            logx.Errorf("Failed to push message to user %d: %v", receiverId, err)
        }
    }
    
    // 确认消息
    msg.Ack(false)
}
```

### 5.2 消费索引队列 (Search RPC)

```go
// internal/consumer/indexConsumer.go
package consumer

import (
    "context"
    "encoding/json"
    "github.com/streadway/amqp"
    "github.com/zeromicro/go-zero/core/logx"
)

type IndexConsumer struct {
    channel *amqp.Channel
    esClient *elasticsearch.Client
}

func NewIndexConsumer(channel *amqp.Channel, esClient *elasticsearch.Client) *IndexConsumer {
    return &IndexConsumer{
        channel:  channel,
        esClient: esClient,
    }
}

func (c *IndexConsumer) Start(ctx context.Context) error {
    // 声明队列
    queue, err := c.channel.QueueDeclare(
        "beehive.message.index",
        true,
        false,
        false,
        false,
        map[string]interface{}{
            "x-message-ttl":           3600000,
            "x-max-length":            100000,
            "x-dead-letter-exchange":  "beehive.dlx.exchange",
            "x-dead-letter-routing-key": "message.index.dlx",
        },
    )
    if err != nil {
        return err
    }
    
    // 绑定队列
    err = c.channel.QueueBind(
        queue.Name,
        "message.index",
        "beehive.message.exchange",
        false,
        nil,
    )
    if err != nil {
        return err
    }
    
    // 设置 QoS（批量消费，每次10条）
    err = c.channel.Qos(10, 0, false)
    if err != nil {
        return err
    }
    
    // 开始消费
    msgs, err := c.channel.Consume(
        queue.Name,
        "",
        false,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return err
    }
    
    logx.Info("IndexConsumer started")
    
    // 消费循环（批量处理）
    go func() {
        buffer := make([]amqp.Delivery, 0, 10)
        ticker := time.NewTicker(time.Second * 5)
        defer ticker.Stop()
        
        for {
            select {
            case <-ctx.Done():
                c.flushBuffer(buffer)
                logx.Info("IndexConsumer stopped")
                return
            case msg, ok := <-msgs:
                if !ok {
                    return
                }
                buffer = append(buffer, msg)
                if len(buffer) >= 10 {
                    c.flushBuffer(buffer)
                    buffer = buffer[:0]
                }
            case <-ticker.C:
                if len(buffer) > 0 {
                    c.flushBuffer(buffer)
                    buffer = buffer[:0]
                }
            }
        }
    }()
    
    return nil
}

func (c *IndexConsumer) flushBuffer(buffer []amqp.Delivery) {
    if len(buffer) == 0 {
        return
    }
    
    // 批量索引到 Elasticsearch
    bulk := c.esClient.Bulk()
    
    for _, msg := range buffer {
        var message Message
        err := json.Unmarshal(msg.Body, &message)
        if err != nil {
            logx.Errorf("Failed to unmarshal message: %v", err)
            msg.Nack(false, false)
            continue
        }
        
        // 只索引文本消息
        if message.ContentType == 1 {
            doc := map[string]interface{}{
                "message_id":      message.MessageId,
                "conversation_id": message.ConversationId,
                "sender_id":       message.SenderId,
                "content":         message.Content,
                "created_at":      message.CreatedAt,
            }
            
            bulk.Index("messages", doc)
        }
        
        msg.Ack(false)
    }
    
    // 执行批量索引
    _, err := bulk.Do(context.Background())
    if err != nil {
        logx.Errorf("Failed to bulk index: %v", err)
    }
}
```

---

## 六、消息可靠性保证

### 6.1 生产者确认

```go
// 开启 Confirm 模式
err = channel.Confirm(false)
if err != nil {
    return err
}

// 发布消息
err = channel.Publish(...)

// 等待确认
confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))
confirmed := <-confirms
if !confirmed.Ack {
    logx.Error("Message not confirmed")
}
```

### 6.2 消费者确认

- 手动 ACK（`autoAck=false`）
- 处理成功后调用 `msg.Ack(false)`
- 处理失败调用 `msg.Nack(false, true)` 重新入队
- 多次失败进入死信队列

### 6.3 消息持久化

- Exchange: `durable=true`
- Queue: `durable=true`
- Message: `deliveryMode=amqp.Persistent`

### 6.4 死信队列

- 消费失败多次的消息进入死信队列
- 定期人工处理死信队列
- 记录失败日志

---

## 七、监控和告警

### 7.1 RabbitMQ 管理界面

访问：http://localhost:15672

用户名/密码：guest/guest

**监控指标**:
- 队列长度
- 消费速率
- 发布速率
- 未确认消息数

### 7.2 告警规则

- 队列长度 > 10000：消费者处理能力不足
- 未确认消息数 > 1000：消费者卡住
- 死信队列消息数 > 100：消费逻辑有问题

### 7.3 Prometheus 指标

```go
// 自定义 Prometheus 指标
var (
    messagePublished = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "rabbitmq_message_published_total",
            Help: "Total number of messages published",
        },
        []string{"routing_key"},
    )
    
    messageConsumed = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "rabbitmq_message_consumed_total",
            Help: "Total number of messages consumed",
        },
        []string{"queue"},
    )
)
```

---

## 八、性能优化

### 8.1 批量消费

- 设置 QoS `prefetchCount=10`
- 批量处理消息后一次性 ACK

### 8.2 连接池

- 使用连接池复用连接
- 每个连接创建多个 Channel

### 8.3 集群部署

- RabbitMQ 集群（3节点）
- 镜像队列（高可用）

### 8.4 消息压缩

- 大消息使用 gzip 压缩
- 设置 `content-encoding: gzip`

---

## 九、最佳实践

1. **幂等性**: 消费者实现幂等性，防止重复消费
2. **重试机制**: 失败消息重试3次后进入死信队列
3. **消息顺序**: 同一会话的消息使用相同的 Routing Key
4. **队列拆分**: 不同业务使用不同队列
5. **监控告警**: 及时发现队列积压和消费异常

---

## 十、常见问题

### Q1: 消息丢失怎么办？

A: 
1. 开启生产者确认
2. 消息持久化
3. 手动 ACK
4. 死信队列兜底

### Q2: 消息重复消费怎么办？

A: 消费者实现幂等性，使用 `message_id` 去重。

### Q3: 消息积压怎么办？

A:
1. 增加消费者实例
2. 优化消费逻辑
3. 批量消费

### Q4: RabbitMQ 性能瓶颈？

A:
1. 使用集群
2. 消息压缩
3. 减少消息体大小
4. 使用 Kafka（高吞吐场景）

---

## 十一、初始化脚本

```bash
#!/bin/bash
# scripts/init_rabbitmq.sh

RABBITMQ_HOST="localhost"
RABBITMQ_PORT="15672"
RABBITMQ_USER="guest"
RABBITMQ_PASS="guest"

# 声明 Exchange
curl -i -u $RABBITMQ_USER:$RABBITMQ_PASS -H "content-type:application/json" \
  -XPUT http://$RABBITMQ_HOST:$RABBITMQ_PORT/api/exchanges/%2F/beehive.message.exchange \
  -d'{"type":"topic","durable":true}'

# 声明队列
curl -i -u $RABBITMQ_USER:$RABBITMQ_PASS -H "content-type:application/json" \
  -XPUT http://$RABBITMQ_HOST:$RABBITMQ_PORT/api/queues/%2F/beehive.message.push \
  -d'{"durable":true,"arguments":{"x-message-ttl":300000,"x-max-length":100000}}'

curl -i -u $RABBITMQ_USER:$RABBITMQ_PASS -H "content-type:application/json" \
  -XPUT http://$RABBITMQ_HOST:$RABBITMQ_PORT/api/queues/%2F/beehive.message.index \
  -d'{"durable":true,"arguments":{"x-message-ttl":3600000,"x-max-length":100000}}'

# 绑定队列
curl -i -u $RABBITMQ_USER:$RABBITMQ_PASS -H "content-type:application/json" \
  -XPOST http://$RABBITMQ_HOST:$RABBITMQ_PORT/api/bindings/%2F/e/beehive.message.exchange/q/beehive.message.push \
  -d'{"routing_key":"message.push"}'

curl -i -u $RABBITMQ_USER:$RABBITMQ_PASS -H "content-type:application/json" \
  -XPOST http://$RABBITMQ_HOST:$RABBITMQ_PORT/api/bindings/%2F/e/beehive.message.exchange/q/beehive.message.index \
  -d'{"routing_key":"message.index"}'

echo "RabbitMQ initialized successfully!"
```
