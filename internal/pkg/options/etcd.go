package options

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

// EtcdOptions etcd 配置选项
type EtcdOptions struct {
	Endpoints   []string      `json:"endpoints" mapstructure:"endpoints"`       // etcd 地址列表
	DialTimeout time.Duration `json:"dial-timeout" mapstructure:"dial-timeout"` // 连接超时
	Username    string        `json:"username" mapstructure:"username"`         // 用户名（可选）
	Password    string        `json:"password" mapstructure:"password"`         // 密码（可选）
	Prefix      string        `json:"prefix" mapstructure:"prefix"`             // 服务注册前缀
}

// NewEtcdOptions 创建 etcd 配置选项
func NewEtcdOptions() *EtcdOptions {
	return &EtcdOptions{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
		Username:    "",
		Password:    "",
		Prefix:      "/beehive/services",
	}
}

// Validate 验证 etcd 配置
func (o *EtcdOptions) Validate() []error {
	var errs []error

	if len(o.Endpoints) == 0 {
		errs = append(errs, fmt.Errorf("etcd.endpoints cannot be empty"))
	}

	for _, endpoint := range o.Endpoints {
		if endpoint == "" {
			errs = append(errs, fmt.Errorf("etcd.endpoints contains empty endpoint"))
		}
	}

	if o.DialTimeout <= 0 {
		errs = append(errs, fmt.Errorf("etcd.dial-timeout must be greater than 0"))
	}

	if o.Prefix == "" {
		errs = append(errs, fmt.Errorf("etcd.prefix cannot be empty"))
	}

	return errs
}

// AddFlags 添加命令行标志
func (o *EtcdOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(&o.Endpoints, "etcd.endpoints", o.Endpoints, "Etcd server endpoints (comma-separated)")
	fs.DurationVar(&o.DialTimeout, "etcd.dial-timeout", o.DialTimeout, "Etcd dial timeout")
	fs.StringVar(&o.Username, "etcd.username", o.Username, "Etcd username (optional)")
	fs.StringVar(&o.Password, "etcd.password", o.Password, "Etcd password (optional)")
	fs.StringVar(&o.Prefix, "etcd.prefix", o.Prefix, "Service registration prefix in etcd")
}

// GetEndpointsString 获取逗号分隔的 endpoints 字符串（用于配置解析）
func (o *EtcdOptions) GetEndpointsString() string {
	return strings.Join(o.Endpoints, ",")
}

// SetEndpointsFromString 从字符串设置 endpoints（用于配置解析）
func (o *EtcdOptions) SetEndpointsFromString(s string) {
	if s == "" {
		o.Endpoints = []string{"localhost:2379"}
		return
	}
	o.Endpoints = strings.Split(s, ",")
	for i := range o.Endpoints {
		o.Endpoints[i] = strings.TrimSpace(o.Endpoints[i])
	}
}
