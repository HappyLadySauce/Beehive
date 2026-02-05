package libnet

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"io"
	"net"
	"testing"
	"time"
)

// mockConn 模拟网络连接用于测试
type mockConn struct {
	readBuf  *bytes.Buffer
	writeBuf *bytes.Buffer
}

func newMockConn() *mockConn {
	return &mockConn{
		readBuf:  new(bytes.Buffer),
		writeBuf: new(bytes.Buffer),
	}
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	return m.readBuf.Read(b)
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	return m.writeBuf.Write(b)
}

func (m *mockConn) Close() error {
	return nil
}

func (m *mockConn) LocalAddr() net.Addr {
	return nil
}

func (m *mockConn) RemoteAddr() net.Addr {
	return nil
}

func (m *mockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

// TestSendReceive 测试发送和接收的完整流程
func TestSendReceive(t *testing.T) {
	// 创建测试消息
	msg := &Message{
		Header: Header{
			Version:   1,
			Status:    0,
			Cmd:       100,
			ServiceId: 200,
			Seq:       12345,
		},
		Data: []byte("Hello, Beehive!"),
	}

	// 创建模拟连接
	conn := newMockConn()
	codec := &BeehiveCodec{conn: conn}

	// 发送消息
	if err := codec.Send(msg); err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	// 将写入缓冲区的数据复制到读取缓冲区
	conn.readBuf = bytes.NewBuffer(conn.writeBuf.Bytes())

	// 接收消息
	receivedMsg, err := codec.Receive()
	if err != nil {
		t.Fatalf("Receive failed: %v", err)
	}

	// 验证 Header
	if receivedMsg.Header.Version != msg.Header.Version {
		t.Errorf("Version mismatch: got %d, want %d", receivedMsg.Header.Version, msg.Header.Version)
	}
	if receivedMsg.Header.Status != msg.Header.Status {
		t.Errorf("Status mismatch: got %d, want %d", receivedMsg.Header.Status, msg.Header.Status)
	}
	if receivedMsg.Header.Cmd != msg.Header.Cmd {
		t.Errorf("Cmd mismatch: got %d, want %d", receivedMsg.Header.Cmd, msg.Header.Cmd)
	}
	if receivedMsg.Header.ServiceId != msg.Header.ServiceId {
		t.Errorf("ServiceId mismatch: got %d, want %d", receivedMsg.Header.ServiceId, msg.Header.ServiceId)
	}
	if receivedMsg.Header.Seq != msg.Header.Seq {
		t.Errorf("Seq mismatch: got %d, want %d", receivedMsg.Header.Seq, msg.Header.Seq)
	}

	// 验证 Data
	if !bytes.Equal(receivedMsg.Data, msg.Data) {
		t.Errorf("Data mismatch: got %s, want %s", receivedMsg.Data, msg.Data)
	}
}

// TestSendReceiveEmptyBody 测试空 Body 的情况
func TestSendReceiveEmptyBody(t *testing.T) {
	msg := &Message{
		Header: Header{
			Version:   1,
			Status:    0,
			Cmd:       100,
			ServiceId: 200,
			Seq:       12345,
		},
		Data: nil,
	}

	conn := newMockConn()
	codec := &BeehiveCodec{conn: conn}

	if err := codec.Send(msg); err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	conn.readBuf = bytes.NewBuffer(conn.writeBuf.Bytes())

	receivedMsg, err := codec.Receive()
	if err != nil {
		t.Fatalf("Receive failed: %v", err)
	}

	if len(receivedMsg.Data) != 0 {
		t.Errorf("Expected empty Data, got %d bytes", len(receivedMsg.Data))
	}
}

// TestSendReceiveLargeBody 测试大 Body 的情况
func TestSendReceiveLargeBody(t *testing.T) {
	// 创建接近最大限制的数据
	largeData := make([]byte, MaxBodySize-100)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	msg := &Message{
		Header: Header{
			Version:   1,
			Status:    0,
			Cmd:       100,
			ServiceId: 200,
			Seq:       12345,
		},
		Data: largeData,
	}

	conn := newMockConn()
	codec := &BeehiveCodec{conn: conn}

	if err := codec.Send(msg); err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	conn.readBuf = bytes.NewBuffer(conn.writeBuf.Bytes())

	receivedMsg, err := codec.Receive()
	if err != nil {
		t.Fatalf("Receive failed: %v", err)
	}

	if !bytes.Equal(receivedMsg.Data, msg.Data) {
		t.Errorf("Large data mismatch")
	}
}

// TestSendFrameTooLarge 测试超大帧的错误处理
func TestSendFrameTooLarge(t *testing.T) {
	// 创建超过最大限制的数据
	tooLargeData := make([]byte, MaxBodySize+1)

	msg := &Message{
		Header: Header{
			Version:   1,
			Status:    0,
			Cmd:       100,
			ServiceId: 200,
			Seq:       12345,
		},
		Data: tooLargeData,
	}

	conn := newMockConn()
	codec := &BeehiveCodec{conn: conn}

	err := codec.Send(msg)
	if err != ErrFrameTooLarge {
		t.Errorf("Expected ErrFrameTooLarge, got %v", err)
	}
}

// TestReceiveChecksumFailed 测试校验和失败的情况
func TestReceiveChecksumFailed(t *testing.T) {
	conn := newMockConn()

	// 手动构造一个校验和错误的包
	payloadSize := uint32(FieldSizeChecksum + HeaderSize + 10)
	
	// 写入长度前缀
	lenBuf := make([]byte, FieldSizeLen)
	binary.BigEndian.PutUint32(lenBuf, payloadSize)
	conn.readBuf.Write(lenBuf)

	// 写入错误的校验和
	checksumBuf := make([]byte, FieldSizeChecksum)
	binary.BigEndian.PutUint32(checksumBuf, 0xDEADBEEF) // 错误的校验和
	conn.readBuf.Write(checksumBuf)

	// 写入 Header
	header := make([]byte, HeaderSize)
	header[0] = 1 // Version
	header[1] = 0 // Status
	binary.BigEndian.PutUint16(header[2:], 100)  // Cmd
	binary.BigEndian.PutUint16(header[4:], 200)  // ServiceId
	binary.BigEndian.PutUint32(header[6:], 12345) // Seq
	conn.readBuf.Write(header)

	// 写入 Body
	body := []byte("test data!")
	conn.readBuf.Write(body)

	codec := &BeehiveCodec{conn: conn}

	_, err := codec.Receive()
	if err != ErrChecksumFailed {
		t.Errorf("Expected ErrChecksumFailed, got %v", err)
	}
}

// TestReceiveFrameTooSmall 测试帧太小的情况
func TestReceiveFrameTooSmall(t *testing.T) {
	conn := newMockConn()

	// 写入一个太小的长度
	lenBuf := make([]byte, FieldSizeLen)
	binary.BigEndian.PutUint32(lenBuf, 5) // 小于最小要求
	conn.readBuf.Write(lenBuf)

	codec := &BeehiveCodec{conn: conn}

	_, err := codec.Receive()
	if err != ErrFrameTooSmall {
		t.Errorf("Expected ErrFrameTooSmall, got %v", err)
	}
}

// TestReceiveFrameTooLarge 测试帧太大的情况
func TestReceiveFrameTooLarge(t *testing.T) {
	conn := newMockConn()

	// 写入一个太大的长度
	lenBuf := make([]byte, FieldSizeLen)
	binary.BigEndian.PutUint32(lenBuf, MaxFrameSize+1000)
	conn.readBuf.Write(lenBuf)

	codec := &BeehiveCodec{conn: conn}

	_, err := codec.Receive()
	if err != ErrFrameTooLarge {
		t.Errorf("Expected ErrFrameTooLarge, got %v", err)
	}
}

// TestReceiveIncompleteFrame 测试不完整帧的情况
func TestReceiveIncompleteFrame(t *testing.T) {
	conn := newMockConn()

	// 写入长度前缀
	payloadSize := uint32(FieldSizeChecksum + HeaderSize + 10)
	lenBuf := make([]byte, FieldSizeLen)
	binary.BigEndian.PutUint32(lenBuf, payloadSize)
	conn.readBuf.Write(lenBuf)

	// 只写入部分数据
	conn.readBuf.Write([]byte{1, 2, 3})

	codec := &BeehiveCodec{conn: conn}

	_, err := codec.Receive()
	if err != io.EOF && err != io.ErrUnexpectedEOF {
		t.Errorf("Expected EOF or UnexpectedEOF, got %v", err)
	}
}

// TestMessageHelpers 测试 Message 的辅助方法
func TestMessageHelpers(t *testing.T) {
	msg := &Message{
		Header: Header{
			Version:   1,
			Status:    0,
			Cmd:       100,
			ServiceId: 200,
			Seq:       12345,
		},
		Data: []byte("test"),
	}

	expectedPayloadSize := FieldSizeChecksum + HeaderSize + len(msg.Data)
	if msg.PayloadSize() != expectedPayloadSize {
		t.Errorf("PayloadSize mismatch: got %d, want %d", msg.PayloadSize(), expectedPayloadSize)
	}

	expectedFrameSize := FrameOverhead + len(msg.Data)
	if msg.FrameSize() != expectedFrameSize {
		t.Errorf("FrameSize mismatch: got %d, want %d", msg.FrameSize(), expectedFrameSize)
	}
}

// TestChecksumCorrectness 测试校验和的正确性
func TestChecksumCorrectness(t *testing.T) {
	data := []byte("test data for checksum")
	expectedChecksum := crc32.ChecksumIEEE(data)

	// 验证 CRC32 计算一致性
	checksum := crc32.ChecksumIEEE(data)
	if checksum != expectedChecksum {
		t.Errorf("Checksum mismatch: got %d, want %d", checksum, expectedChecksum)
	}
}

// BenchmarkSend 性能测试：发送消息
func BenchmarkSend(b *testing.B) {
	msg := &Message{
		Header: Header{
			Version:   1,
			Status:    0,
			Cmd:       100,
			ServiceId: 200,
			Seq:       12345,
		},
		Data: make([]byte, 1024),
	}

	conn := newMockConn()
	codec := &BeehiveCodec{conn: conn}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn.writeBuf.Reset()
		if err := codec.Send(msg); err != nil {
			b.Fatalf("Send failed: %v", err)
		}
	}
}

// BenchmarkReceive 性能测试：接收消息
func BenchmarkReceive(b *testing.B) {
	msg := &Message{
		Header: Header{
			Version:   1,
			Status:    0,
			Cmd:       100,
			ServiceId: 200,
			Seq:       12345,
		},
		Data: make([]byte, 1024),
	}

	conn := newMockConn()
	codec := &BeehiveCodec{conn: conn}

	// 预先发送消息
	if err := codec.Send(msg); err != nil {
		b.Fatalf("Send failed: %v", err)
	}
	wireData := conn.writeBuf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn.readBuf = bytes.NewBuffer(wireData)
		if _, err := codec.Receive(); err != nil {
			b.Fatalf("Receive failed: %v", err)
		}
	}
}
