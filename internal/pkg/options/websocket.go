package options

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
)

// WebSocketOptions WebSocket 配置
type WebSocketOptions struct {
	ReadTimeout  time.Duration `json:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout" mapstructure:"write_timeout"`
	PingInterval time.Duration `json:"ping_interval" mapstructure:"ping_interval"`
}

// NewWebSocketOptions 创建 WebSocket 配置选项
func NewWebSocketOptions() *WebSocketOptions {
	return &WebSocketOptions{
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 10 * time.Second,
		PingInterval: 30 * time.Second,
	}
}

// Validate 验证 WebSocket 配置
func (o *WebSocketOptions) Validate() []error {
	var errs []error
	if o.ReadTimeout <= 0 {
		errs = append(errs, fmt.Errorf("websocket.read_timeout must be greater than 0"))
	}
	if o.WriteTimeout <= 0 {
		errs = append(errs, fmt.Errorf("websocket.write_timeout must be greater than 0"))
	}
	if o.PingInterval <= 0 {
		errs = append(errs, fmt.Errorf("websocket.ping_interval must be greater than 0"))
	}
	return errs
}

// AddFlags 添加命令行标志
func (o *WebSocketOptions) AddFlags(fs *pflag.FlagSet) {
	fs.DurationVar(&o.ReadTimeout, "websocket.read-timeout", o.ReadTimeout, "WebSocket read timeout")
	fs.DurationVar(&o.WriteTimeout, "websocket.write-timeout", o.WriteTimeout, "WebSocket write timeout")
	fs.DurationVar(&o.PingInterval, "websocket.ping-interval", o.PingInterval, "WebSocket ping interval")
}
