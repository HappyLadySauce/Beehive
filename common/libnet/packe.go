package libnet

import "errors"

/*
自定义消息包结构

网络传输格式：
[LenPrefix(4)][Checksum(4)][Version(1)][Status(1)][Cmd(2)][ServiceId(2)][Seq(4)][Data...]
 ← 传输层 →    ←─────────────────── 应用层（Header + Body）───────────────────→

- LenPrefix（4字节）：后续数据长度（Checksum + Header + Body）
- Checksum（4字节）：CRC32 校验和（校验 Header + Body）
- Header（10字节）：Version(1) + Status(1) + Cmd(2) + ServiceId(2) + Seq(4)
- Body：变长数据
*/

const (
	// 字段大小定义
	FieldSizeLen       = 4 // 长度前缀字段
	FieldSizeChecksum  = 4 // 校验和字段
	FieldSizeVersion   = 1 // 版本字段
	FieldSizeStatus    = 1 // 状态字段
	FieldSizeCmd       = 2 // 命令字段
	FieldSizeServiceId = 2 // 服务ID字段
	FieldSizeSeq       = 4 // 序列号字段

	// 协议头部大小（不含长度前缀和校验和）
	HeaderSize = FieldSizeVersion + FieldSizeStatus + FieldSizeCmd +
		FieldSizeServiceId + FieldSizeSeq

	// 完整包的固定开销
	FrameOverhead = FieldSizeLen + FieldSizeChecksum + HeaderSize

	// 最大限制
	MaxBodySize  = 1 << 12                      // 4KB
	MaxFrameSize = FrameOverhead + MaxBodySize  // 完整包最大长度

	// 偏移量定义（针对完整的网络包缓冲区）
	OffsetLen       = 0
	OffsetChecksum  = OffsetLen + FieldSizeLen
	OffsetVersion   = OffsetChecksum + FieldSizeChecksum
	OffsetStatus    = OffsetVersion + FieldSizeVersion
	OffsetCmd       = OffsetStatus + FieldSizeStatus
	OffsetServiceId = OffsetCmd + FieldSizeCmd
	OffsetSeq       = OffsetServiceId + FieldSizeServiceId
	OffsetBody      = OffsetSeq + FieldSizeSeq
)

var (
	ErrFrameTooLarge  = errors.New("frame size exceeds maximum")
	ErrFrameTooSmall  = errors.New("frame size too small")
	ErrChecksumFailed = errors.New("checksum verification failed")
)

// Header 应用层协议头部（不包含传输层的长度前缀）
type Header struct {
	Version   uint8  // 版本
	Status    uint8  // 状态
	Cmd       uint16 // 命令
	ServiceId uint16 // 服务ID
	Seq       uint32 // 序列号
}

// Message 完整消息
type Message struct {
	Header // 包头 10字节
	Data   []byte // 数据
}

// FrameSize 计算完整帧的大小
func (m *Message) FrameSize() int {
	return FrameOverhead + len(m.Data)
}

// PayloadSize 计算有效载荷大小（Checksum + Header + Body）
func (m *Message) PayloadSize() int {
	return FieldSizeChecksum + HeaderSize + len(m.Data)
}
