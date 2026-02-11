package libnet

import (
	"encoding/binary"
	"hash/crc32"
	"io"
	"net"
	"sync"
	"time"
)

// 缓冲池，减少内存分配
var bufferPool = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, MaxFrameSize)
		return &buf
	},
}

// 获取缓冲区
func acquireBuffer() *[]byte {
	return bufferPool.Get().(*[]byte)
}

// 释放缓冲区
func releaseBuffer(buf *[]byte) {
	bufferPool.Put(buf)
}

// 自定义协议接口
type Protocol interface {
	NewCodec(conn net.Conn) Codec
}

// 自定义编码器接口
type Codec interface {
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	Receive() (*Message, error)
	Send(msg Message) error
	Close() error
}

// 蜂巢协议
type BeehiveProtocol struct {}

// 创建蜂巢协议
func NewBeehiveProtocol() *BeehiveProtocol {
	return &BeehiveProtocol{}
}

// 创建编码器
func (p *BeehiveProtocol) NewCodec(conn net.Conn) Codec {
	return &BeehiveCodec{conn: conn}
}

// 编码器
type BeehiveCodec struct {
	conn net.Conn
}

// SetReadDeadline 设置读取超时
func (c *BeehiveCodec) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline 设置写入超时
func (c *BeehiveCodec) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

// Receive 接收消息
func (c *BeehiveCodec) Receive() (*Message, error) {
	// 1. 读取长度前缀
	var lenBuf [FieldSizeLen]byte
	if _, err := io.ReadFull(c.conn, lenBuf[:]); err != nil {
		return nil, err
	}
	payloadSize := binary.BigEndian.Uint32(lenBuf[:])

	// 2. 校验长度
	if payloadSize < FieldSizeChecksum+HeaderSize {
		return nil, ErrFrameTooSmall
	}
	if payloadSize > MaxFrameSize-FieldSizeLen {
		return nil, ErrFrameTooLarge
	}

	// 3. 读取有效载荷（使用缓冲池）
	bufPtr := acquireBuffer()
	defer releaseBuffer(bufPtr)

	payload := (*bufPtr)[:payloadSize]
	if _, err := io.ReadFull(c.conn, payload); err != nil {
		return nil, err
	}

	// 4. 校验 checksum
	receivedChecksum := binary.BigEndian.Uint32(payload[:FieldSizeChecksum])
	data := payload[FieldSizeChecksum:]
	calculatedChecksum := crc32.ChecksumIEEE(data)
	if receivedChecksum != calculatedChecksum {
		return nil, ErrChecksumFailed
	}

	// 5. 解析 Header
	var msg Message
	msg.Version = data[0]
	msg.Status = data[1]
	msg.Cmd = binary.BigEndian.Uint16(data[2:4])
	msg.ServiceId = binary.BigEndian.Uint16(data[4:6])
	msg.Seq = binary.BigEndian.Uint32(data[6:10])

	// 6. 复制 Body（需要复制因为 buffer 会被归还池）
	bodySize := len(data) - HeaderSize
	if bodySize > 0 {
		msg.Data = make([]byte, bodySize)
		copy(msg.Data, data[HeaderSize:])
	}

	return &msg, nil
}

// Send 发送消息
func (c *BeehiveCodec) Send(msg Message) error {
	payloadSize := msg.PayloadSize()
	frameSize := msg.FrameSize()

	// 1. 校验大小
	if frameSize > MaxFrameSize {
		return ErrFrameTooLarge
	}

	// 2. 使用缓冲池
	bufPtr := acquireBuffer()
	defer releaseBuffer(bufPtr)

	buf := (*bufPtr)[:frameSize]

	// 3. 写入长度前缀
	binary.BigEndian.PutUint32(buf[OffsetLen:], uint32(payloadSize))

	// 4. 写入 Header
	buf[OffsetVersion] = msg.Version
	buf[OffsetStatus] = msg.Status
	binary.BigEndian.PutUint16(buf[OffsetCmd:], msg.Cmd)
	binary.BigEndian.PutUint16(buf[OffsetServiceId:], msg.ServiceId)
	binary.BigEndian.PutUint32(buf[OffsetSeq:], msg.Seq)

	// 5. 写入 Body
	copy(buf[OffsetBody:], msg.Data)

	// 6. 计算并写入 Checksum（校验 Header + Body）
	checksum := crc32.ChecksumIEEE(buf[OffsetVersion:frameSize])
	binary.BigEndian.PutUint32(buf[OffsetChecksum:], checksum)

	// 7. 发送
	_, err := c.conn.Write(buf)
	return err
}

func (c *BeehiveCodec) Close() error {
	return c.conn.Close()
}