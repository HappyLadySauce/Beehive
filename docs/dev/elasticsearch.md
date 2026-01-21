# Beehive IM Elasticsearch 设计文档

## 一、Elasticsearch 概述

### 1.1 技术选型

**搜索引擎**: Elasticsearch 8.11.0

**选型原因**:
- 强大的全文检索能力
- 分布式架构，水平扩展
- 近实时搜索
- 支持中文分词（IK 分词器）
- 丰富的查询 DSL
- 高亮显示支持
- RESTful API

### 1.2 使用场景

1. **历史消息全文检索**: 根据关键词搜索历史消息
2. **搜索结果高亮**: 关键词高亮显示
3. **按会话过滤**: 在指定会话中搜索
4. **分页查询**: 支持大数据量分页

---

## 二、索引设计

### 2.1 消息索引 (messages)

**索引名称**: `messages`

**分片配置**:
- **主分片**: 3
- **副本分片**: 1
- **刷新间隔**: 5s（近实时）

**索引设置**:

```json
{
    "settings": {
        "number_of_shards": 3,
        "number_of_replicas": 1,
        "refresh_interval": "5s",
        "max_result_window": 10000,
        "analysis": {
            "analyzer": {
                "ik_smart_analyzer": {
                    "type": "custom",
                    "tokenizer": "ik_smart"
                },
                "ik_max_word_analyzer": {
                    "type": "custom",
                    "tokenizer": "ik_max_word"
                }
            }
        }
    },
    "mappings": {
        "properties": {
            "message_id": {
                "type": "long"
            },
            "conversation_id": {
                "type": "long"
            },
            "sender_id": {
                "type": "long"
            },
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
            "content_type": {
                "type": "integer"
            },
            "created_at": {
                "type": "date",
                "format": "epoch_second"
            }
        }
    }
}
```

**字段说明**:

| 字段 | 类型 | 说明 |
|------|------|------|
| message_id | long | 消息ID（主键） |
| conversation_id | long | 会话ID（用于过滤） |
| sender_id | long | 发送者ID |
| content | text | 消息内容（全文检索字段） |
| content_type | integer | 消息类型（1:文本 2:图片 3:语音） |
| created_at | date | 发送时间（Unix时间戳） |

**索引策略**:
- 仅索引文本消息（content_type = 1）
- 图片、语音消息不索引内容，但保留元数据
- 按月创建索引（messages-2026-01）

---

### 2.2 IK 分词器安装

**IK 分词器**: 中文分词插件

**安装步骤**:

```bash
# 下载 IK 分词器
cd /usr/share/elasticsearch/plugins
wget https://github.com/medcl/elasticsearch-analysis-ik/releases/download/v8.11.0/elasticsearch-analysis-ik-8.11.0.zip

# 解压
unzip elasticsearch-analysis-ik-8.11.0.zip -d ik

# 重启 Elasticsearch
docker restart beehive-elasticsearch
```

**分词器说明**:
- **ik_smart**: 粗粒度分词，适合搜索
- **ik_max_word**: 细粒度分词，适合索引

**测试分词**:

```bash
curl -X POST "localhost:9200/_analyze" -H 'Content-Type: application/json' -d'
{
  "analyzer": "ik_smart",
  "text": "我想和你一起吃饭"
}'

# 结果：["我", "想", "和", "你", "一起", "吃饭"]
```

---

## 三、索引操作

### 3.1 创建索引

```bash
curl -X PUT "localhost:9200/messages" -H 'Content-Type: application/json' -d'
{
    "settings": {
        "number_of_shards": 3,
        "number_of_replicas": 1,
        "refresh_interval": "5s",
        "analysis": {
            "analyzer": {
                "ik_smart_analyzer": {
                    "type": "custom",
                    "tokenizer": "ik_smart"
                },
                "ik_max_word_analyzer": {
                    "type": "custom",
                    "tokenizer": "ik_max_word"
                }
            }
        }
    },
    "mappings": {
        "properties": {
            "message_id": {"type": "long"},
            "conversation_id": {"type": "long"},
            "sender_id": {"type": "long"},
            "content": {
                "type": "text",
                "analyzer": "ik_max_word_analyzer",
                "search_analyzer": "ik_smart_analyzer"
            },
            "content_type": {"type": "integer"},
            "created_at": {"type": "date", "format": "epoch_second"}
        }
    }
}'
```

### 3.2 索引文档（单条）

```bash
curl -X POST "localhost:9200/messages/_doc/3001" -H 'Content-Type: application/json' -d'
{
    "message_id": 3001,
    "conversation_id": 2001,
    "sender_id": 1002,
    "content": "这是一条测试消息",
    "content_type": 1,
    "created_at": 1705838400
}'
```

### 3.3 批量索引（Bulk API）

```bash
curl -X POST "localhost:9200/_bulk" -H 'Content-Type: application/json' -d'
{"index":{"_index":"messages","_id":"3001"}}
{"message_id":3001,"conversation_id":2001,"sender_id":1002,"content":"你好","content_type":1,"created_at":1705838400}
{"index":{"_index":"messages","_id":"3002"}}
{"message_id":3002,"conversation_id":2001,"sender_id":1001,"content":"你好呀","content_type":1,"created_at":1705838500}
'
```

### 3.4 搜索文档

```bash
curl -X POST "localhost:9200/messages/_search" -H 'Content-Type: application/json' -d'
{
    "query": {
        "bool": {
            "must": [
                {
                    "match": {
                        "content": "测试"
                    }
                }
            ],
            "filter": [
                {
                    "term": {
                        "conversation_id": 2001
                    }
                }
            ]
        }
    },
    "highlight": {
        "fields": {
            "content": {}
        }
    },
    "from": 0,
    "size": 20,
    "sort": [
        {
            "created_at": {
                "order": "desc"
            }
        }
    ]
}'
```

### 3.5 删除文档

```bash
# 根据ID删除
curl -X DELETE "localhost:9200/messages/_doc/3001"

# 根据查询删除
curl -X POST "localhost:9200/messages/_delete_by_query" -H 'Content-Type: application/json' -d'
{
    "query": {
        "term": {
            "message_id": 3001
        }
    }
}'
```

---

## 四、Go 客户端实现

### 4.1 初始化客户端

```go
package es

import (
    "github.com/elastic/go-elasticsearch/v8"
    "log"
)

func NewESClient(addresses []string) (*elasticsearch.Client, error) {
    cfg := elasticsearch.Config{
        Addresses: addresses,
    }
    
    client, err := elasticsearch.NewClient(cfg)
    if err != nil {
        return nil, err
    }
    
    // 测试连接
    res, err := client.Info()
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()
    
    log.Println("Elasticsearch connected")
    return client, nil
}
```

### 4.2 索引消息

```go
func (e *ESClient) IndexMessage(ctx context.Context, msg *Message) error {
    doc := map[string]interface{}{
        "message_id":      msg.MessageId,
        "conversation_id": msg.ConversationId,
        "sender_id":       msg.SenderId,
        "content":         msg.Content,
        "content_type":    msg.ContentType,
        "created_at":      msg.CreatedAt,
    }
    
    body, err := json.Marshal(doc)
    if err != nil {
        return err
    }
    
    res, err := e.client.Index(
        "messages",
        bytes.NewReader(body),
        e.client.Index.WithDocumentID(fmt.Sprintf("%d", msg.MessageId)),
        e.client.Index.WithContext(ctx),
    )
    if err != nil {
        return err
    }
    defer res.Body.Close()
    
    if res.IsError() {
        return fmt.Errorf("index error: %s", res.String())
    }
    
    return nil
}
```

### 4.3 批量索引

```go
func (e *ESClient) BulkIndexMessages(ctx context.Context, messages []*Message) error {
    var buf bytes.Buffer
    
    for _, msg := range messages {
        // 仅索引文本消息
        if msg.ContentType != 1 {
            continue
        }
        
        meta := map[string]interface{}{
            "index": map[string]interface{}{
                "_index": "messages",
                "_id":    fmt.Sprintf("%d", msg.MessageId),
            },
        }
        
        doc := map[string]interface{}{
            "message_id":      msg.MessageId,
            "conversation_id": msg.ConversationId,
            "sender_id":       msg.SenderId,
            "content":         msg.Content,
            "content_type":    msg.ContentType,
            "created_at":      msg.CreatedAt,
        }
        
        metaJSON, _ := json.Marshal(meta)
        docJSON, _ := json.Marshal(doc)
        
        buf.Write(metaJSON)
        buf.WriteByte('\n')
        buf.Write(docJSON)
        buf.WriteByte('\n')
    }
    
    res, err := e.client.Bulk(
        bytes.NewReader(buf.Bytes()),
        e.client.Bulk.WithContext(ctx),
    )
    if err != nil {
        return err
    }
    defer res.Body.Close()
    
    if res.IsError() {
        return fmt.Errorf("bulk error: %s", res.String())
    }
    
    return nil
}
```

### 4.4 搜索消息

```go
type SearchRequest struct {
    Keyword        string
    ConversationId int64
    Page           int
    PageSize       int
}

type SearchResponse struct {
    Messages []*SearchMessage
    Total    int64
}

type SearchMessage struct {
    MessageId      int64
    ConversationId int64
    SenderId       int64
    Content        string
    Highlight      string
    CreatedAt      int64
}

func (e *ESClient) SearchMessages(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
    // 构建查询
    query := map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []interface{}{
                    map[string]interface{}{
                        "match": map[string]interface{}{
                            "content": req.Keyword,
                        },
                    },
                },
            },
        },
        "highlight": map[string]interface{}{
            "fields": map[string]interface{}{
                "content": map[string]interface{}{
                    "pre_tags":  []string{"<em>"},
                    "post_tags": []string{"</em>"},
                },
            },
        },
        "from": (req.Page - 1) * req.PageSize,
        "size": req.PageSize,
        "sort": []interface{}{
            map[string]interface{}{
                "created_at": map[string]interface{}{
                    "order": "desc",
                },
            },
        },
    }
    
    // 按会话过滤
    if req.ConversationId > 0 {
        query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = []interface{}{
            map[string]interface{}{
                "term": map[string]interface{}{
                    "conversation_id": req.ConversationId,
                },
            },
        }
    }
    
    body, _ := json.Marshal(query)
    
    // 执行搜索
    res, err := e.client.Search(
        e.client.Search.WithContext(ctx),
        e.client.Search.WithIndex("messages"),
        e.client.Search.WithBody(bytes.NewReader(body)),
    )
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()
    
    if res.IsError() {
        return nil, fmt.Errorf("search error: %s", res.String())
    }
    
    // 解析响应
    var result map[string]interface{}
    if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    hits := result["hits"].(map[string]interface{})
    total := int64(hits["total"].(map[string]interface{})["value"].(float64))
    
    messages := make([]*SearchMessage, 0)
    for _, hit := range hits["hits"].([]interface{}) {
        h := hit.(map[string]interface{})
        source := h["_source"].(map[string]interface{})
        
        msg := &SearchMessage{
            MessageId:      int64(source["message_id"].(float64)),
            ConversationId: int64(source["conversation_id"].(float64)),
            SenderId:       int64(source["sender_id"].(float64)),
            Content:        source["content"].(string),
            CreatedAt:      int64(source["created_at"].(float64)),
        }
        
        // 高亮
        if highlight, ok := h["highlight"].(map[string]interface{}); ok {
            if content, ok := highlight["content"].([]interface{}); ok && len(content) > 0 {
                msg.Highlight = content[0].(string)
            }
        }
        
        messages = append(messages, msg)
    }
    
    return &SearchResponse{
        Messages: messages,
        Total:    total,
    }, nil
}
```

### 4.5 在 Search RPC 中使用

```go
// internal/svc/servicecontext.go
type ServiceContext struct {
    Config   config.Config
    ESClient *es.ESClient
}

func NewServiceContext(c config.Config) *ServiceContext {
    esClient, err := es.NewESClient(c.Elasticsearch.Addresses)
    if err != nil {
        panic(err)
    }
    
    return &ServiceContext{
        Config:   c,
        ESClient: esClient,
    }
}

// internal/logic/searchmessageslogic.go
func (l *SearchMessagesLogic) SearchMessages(req *search.SearchMessagesRequest) (*search.SearchMessagesResponse, error) {
    resp, err := l.svcCtx.ESClient.SearchMessages(l.ctx, &es.SearchRequest{
        Keyword:        req.Keyword,
        ConversationId: req.ConversationId,
        Page:           int(req.Page),
        PageSize:       int(req.PageSize),
    })
    if err != nil {
        return nil, err
    }
    
    messages := make([]*search.SearchMessageInfo, 0, len(resp.Messages))
    for _, msg := range resp.Messages {
        messages = append(messages, &search.SearchMessageInfo{
            MessageId:      msg.MessageId,
            ConversationId: msg.ConversationId,
            SenderId:       msg.SenderId,
            Content:        msg.Content,
            Highlight:      msg.Highlight,
            CreatedAt:      msg.CreatedAt,
        })
    }
    
    return &search.SearchMessagesResponse{
        Messages: messages,
        Total:    resp.Total,
    }, nil
}
```

---

## 五、性能优化

### 5.1 索引优化

**按月创建索引**:

```bash
# 当月索引
messages-2026-01

# 下月索引
messages-2026-02
```

**优点**:
- 旧索引可以关闭或删除
- 减少单个索引大小
- 提高查询性能

**别名机制**:

```bash
# 创建别名
curl -X POST "localhost:9200/_aliases" -H 'Content-Type: application/json' -d'
{
    "actions": [
        {"add": {"index": "messages-2026-01", "alias": "messages"}}
    ]
}'

# 查询使用别名
GET /messages/_search
```

### 5.2 查询优化

**使用过滤器**:

```json
{
    "query": {
        "bool": {
            "must": [
                {"match": {"content": "关键词"}}
            ],
            "filter": [
                {"term": {"conversation_id": 2001}},
                {"range": {"created_at": {"gte": 1705752000}}}
            ]
        }
    }
}
```

**使用 `_source` 过滤**:

```json
{
    "query": {...},
    "_source": ["message_id", "content", "created_at"]
}
```

### 5.3 批量操作

- 使用 Bulk API 批量索引
- 批量大小：1000 - 5000 条
- 并发请求数：2 - 4

### 5.4 缓存策略

- 热点查询结果缓存到 Redis
- TTL: 5 分钟
- Key: `search:{keyword}:{conversation_id}:{page}`

---

## 六、数据备份和恢复

### 6.1 快照备份

```bash
# 创建快照仓库
curl -X PUT "localhost:9200/_snapshot/beehive_backup" -H 'Content-Type: application/json' -d'
{
    "type": "fs",
    "settings": {
        "location": "/backup/elasticsearch"
    }
}'

# 创建快照
curl -X PUT "localhost:9200/_snapshot/beehive_backup/snapshot_20260121"

# 恢复快照
curl -X POST "localhost:9200/_snapshot/beehive_backup/snapshot_20260121/_restore"
```

### 6.2 定时备份

```bash
# 每天凌晨 3:00 备份
0 3 * * * curl -X PUT "localhost:9200/_snapshot/beehive_backup/snapshot_$(date +\%Y\%m\%d)"
```

---

## 七、监控和告警

### 7.1 集群健康

```bash
# 查看集群健康
curl -X GET "localhost:9200/_cluster/health"

# 响应
{
    "cluster_name": "beehive-es",
    "status": "green",  // green: 正常, yellow: 警告, red: 错误
    "number_of_nodes": 3,
    "active_shards": 9
}
```

### 7.2 索引统计

```bash
# 查看索引大小
curl -X GET "localhost:9200/_cat/indices?v"

# 响应
health status index           pri rep docs.count store.size
green  open   messages-2026-01  3   1    1000000     1.2gb
```

### 7.3 慢查询日志

```bash
# 开启慢查询日志
curl -X PUT "localhost:9200/messages/_settings" -H 'Content-Type: application/json' -d'
{
    "index.search.slowlog.threshold.query.warn": "10s",
    "index.search.slowlog.threshold.query.info": "5s"
}'
```

### 7.4 Kibana 可视化

访问：http://localhost:5601

**功能**:
- 索引管理
- 查询测试
- 日志分析
- 数据可视化

---

## 八、常见问题

### Q1: IK 分词器安装失败？

A: 
1. 检查版本是否匹配（ES 8.11.0 对应 IK 8.11.0）
2. 检查插件目录权限
3. 重启 Elasticsearch

### Q2: 搜索结果不准确？

A:
1. 检查分词器配置
2. 使用 `_analyze` API 测试分词
3. 调整 `minimum_should_match` 参数

### Q3: 搜索速度慢？

A:
1. 使用过滤器（filter）而不是查询（query）
2. 减少返回字段（_source 过滤）
3. 增加分片数
4. 优化查询语句

### Q4: 索引占用空间大？

A:
1. 按月创建索引
2. 删除旧索引
3. 使用 `_forcemerge` 合并段
4. 压缩存储（best_compression）

---

## 九、初始化脚本

**文件路径**: `/opt/Beehive/scripts/init_es.sh`

```bash
#!/bin/bash

ES_HOST="localhost:9200"

echo "Initializing Elasticsearch..."

# 1. 检查 IK 分词器
echo "Checking IK analyzer..."
curl -s "$ES_HOST/_cat/plugins" | grep analysis-ik || {
    echo "IK analyzer not found! Please install it first."
    exit 1
}

# 2. 创建索引
echo "Creating messages index..."
curl -X PUT "$ES_HOST/messages" -H 'Content-Type: application/json' -d'
{
    "settings": {
        "number_of_shards": 3,
        "number_of_replicas": 1,
        "refresh_interval": "5s",
        "analysis": {
            "analyzer": {
                "ik_smart_analyzer": {
                    "type": "custom",
                    "tokenizer": "ik_smart"
                },
                "ik_max_word_analyzer": {
                    "type": "custom",
                    "tokenizer": "ik_max_word"
                }
            }
        }
    },
    "mappings": {
        "properties": {
            "message_id": {"type": "long"},
            "conversation_id": {"type": "long"},
            "sender_id": {"type": "long"},
            "content": {
                "type": "text",
                "analyzer": "ik_max_word_analyzer",
                "search_analyzer": "ik_smart_analyzer"
            },
            "content_type": {"type": "integer"},
            "created_at": {"type": "date", "format": "epoch_second"}
        }
    }
}'

echo ""
echo "Elasticsearch initialized successfully!"
```

执行：

```bash
chmod +x scripts/init_es.sh
./scripts/init_es.sh
```

---

## 十、最佳实践

1. **索引策略**: 按月创建索引，使用别名查询
2. **分词器**: 索引用 `ik_max_word`，搜索用 `ik_smart`
3. **批量操作**: 使用 Bulk API 批量索引
4. **过滤优先**: 使用 filter 而不是 query
5. **缓存**: 热点查询结果缓存到 Redis
6. **监控**: 定期检查集群健康和索引大小
7. **备份**: 每天定时快照备份
