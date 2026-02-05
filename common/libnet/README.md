# Beehive 网络协议库

## 协议格式

### 网络传输格式

```
[LenPrefix(4)][Checksum(4)][Version(1)][Status(1)][Cmd(2)][ServiceId(2)][Seq(4)][Data...]
 ← 传输层 →    ←─────────────────── 应用层（Header + Body）───────────────────→
```

### 字段说明

| 字段 | 大小 | 说明 |
|------|------|------|
| LenPrefix | 4字节 | 后续数据长度（Checksum + Header + Body） |
| Checksum | 4字节 | CRC32 校验和（校验 Header + Body） |
| Version | 1字节 | 协议版本 |
| Status | 1字节 | 状态码 |
| Cmd | 2字节 | 命令类型 |
| ServiceId | 2字节 | 服务ID |
| Seq | 4字节 | 序列号 |
| Data | 变长 | 消息体数据（最大 4KB） |

### 常量定义

```go
FieldSizeLen       = 4    // 长度前缀字段
FieldSizeChecksum  = 4    // 校验和字段
HeaderSize         = 10   // 协议头部大小（不含长度前缀和校验和）
FrameOverhead      = 18   // 完整包的固定开销
MaxBodySize        = 4096 // 最大数据长度 4KB
MaxFrameSize       = 4114 // 完整包最大长度
```

## 使用示例

### 创建连接和编解码器

```go
// 创建协议
protocol := libnet.NewBeehiveProtocol()

// 从网络连接创建编解码器
codec := protocol.NewCodec(conn)
defer codec.Close()
```

### 发送消息

```go
msg := &libnet.Message{
    Header: libnet.Header{
        Version:   1,
        Status:    0,
        Cmd:       100,
        ServiceId: 200,
        Seq:       12345,
    },
    Data: []byte("Hello, Beehive!"),
}

if err := codec.Send(msg); err != nil {
    log.Printf("发送失败: %v", err)
}
```

### 接收消息

```go
msg, err := codec.Receive()
if err != nil {
    log.Printf("接收失败: %v", err)
    return
}

log.Printf("收到消息: Version=%d, Cmd=%d, Data=%s",
    msg.Header.Version, msg.Header.Cmd, msg.Data)
```

### 设置超时

```go
// 设置读取超时
codec.SetReadDeadline(time.Now().Add(5 * time.Second))

// 设置写入超时
codec.SetWriteDeadline(time.Now().Add(5 * time.Second))
```

## 设计特点

### 1. 清晰分层

- **传输层**：长度前缀（LenPrefix），用于 TCP 流的分帧
- **应用层**：校验和 + 协议头 + 数据体

### 2. 性能优化

- 使用 `sync.Pool` 缓冲池减少内存分配
- 零拷贝设计，减少不必要的数据复制

### 3. 数据完整性

- CRC32 校验和确保数据在传输过程中不被损坏
- 长度校验防止非法数据包

### 4. 易于使用

- `Header` 结构不包含传输层的 `Len` 字段，职责清晰
- 提供 `FrameSize()` 和 `PayloadSize()` 辅助方法
- 常量命名清晰，易于理解

## 错误处理

| 错误 | 说明 |
|------|------|
| `ErrFrameTooLarge` | 帧大小超过最大限制 |
| `ErrFrameTooSmall` | 帧大小小于最小要求 |
| `ErrChecksumFailed` | 校验和验证失败 |

## 性能指标

基于 Benchmark 测试：

- 发送 1KB 消息：~1.5μs/op
- 接收 1KB 消息：~2.0μs/op
- 内存分配：使用缓冲池后显著减少

## 协议演进

如果未来需要修改协议：

1. **去掉校验和**：将 `FieldSizeChecksum` 设为 0，跳过校验逻辑
2. **增加字段**：在 `Header` 末尾添加新字段，更新 `HeaderSize` 和偏移量
3. **版本控制**：使用 `Version` 字段区分不同版本的协议

## 测试

运行单元测试：

```bash
cd common/libnet
go test -v
```

运行性能测试：

```bash
go test -bench=. -benchmem
```
