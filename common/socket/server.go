package socket

import (
	"crypto/tls"
	"errors"
	"net"
	"strings"
	"io"
	"time"

	"github.com/HappyLadySauce/Beehive/common/libnet"
	"github.com/zeromicro/go-zero/core/logx"
)

// server 服务器
type Server struct {
	Name		string
	Manager		*libnet.Manager
	Listener	net.Listener
	Protocol	libnet.Protocol
	SendChanSize	int
}

// 创建服务器
func NewServer(name string, l net.Listener, p libnet.Protocol, sendChanSize int) *Server {
	return &Server{
		Name: name,
		Manager: libnet.NewManager(name),
		Listener: l,
		Protocol: p,
		SendChanSize: sendChanSize,
	}
}

// 接受会话
func (s *Server) Accept() (*libnet.Session, error) {
	// 接受会话
	var tempDelay time.Duration
	// 循环接受会话
	for {
		// 接受会话
		conn, err := s.Listener.Accept()
		// 如果接受会话失败
		if err != nil {
			var ne net.Error
			// 如果接受会话失败是超时错误
			if errors.As(err, &ne) && ne.Timeout() {
				// 如果临时延迟为0，则设置为5毫秒
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					// 如果临时延迟大于1秒，则设置为1秒
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					// 如果临时延迟大于1秒，则设置为1秒
					tempDelay = max
				}
				// 记录错误日志
				logx.Errorf("[Accept] accept error: %v, tempDelay: %v", err, tempDelay)
				// 睡眠临时延迟
				time.Sleep(tempDelay)
				continue
			}
			// 如果接受会话失败是关闭网络连接错误
			if strings.Contains(err.Error(), "use of closed network connection") {
				// 返回EOF错误
				return nil, io.EOF
			}
			// 返回错误
			return nil, err
		}
		// 创建会话并返回 注册会话在登录阶段之后进行
		return libnet.NewSession(s.Manager, s.Protocol.NewCodec(conn), s.SendChanSize), nil
	}
}

func (s *Server) Close() {
	s.Listener.Close()
	s.Manager.Close()
}

func NewServe(name, address string, protocol libnet.Protocol, sendChanSize int) (*Server, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewServer(name, listener, protocol, sendChanSize), nil
}

func NewTlsServer(name, address string, config *tls.Config, protocol libnet.Protocol, sendChanSize int) (*Server, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	listener, err := tls.Listen("tcp", addr.String(), config)
	if err != nil {
		return nil, err
	}
	return NewServer(name, listener, protocol, sendChanSize), nil
}
