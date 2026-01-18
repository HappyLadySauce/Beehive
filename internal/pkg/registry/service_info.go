package registry

import (
	"encoding/json"
	"fmt"
)

// ServiceInfo 服务信息
type ServiceInfo struct {
	ServiceName string            `json:"service_name"` // 服务名称
	Address     string            `json:"address"`      // 服务地址
	Port        int               `json:"port"`         // 服务端口
	InstanceID  string            `json:"instance_id"`  // 实例ID
	Metadata    map[string]string `json:"metadata"`     // 元数据
}

// ToJSON 将服务信息转换为 JSON
func (s *ServiceInfo) ToJSON() (string, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON 从 JSON 解析服务信息
func (s *ServiceInfo) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), s)
}

// GetAddress 获取完整地址（address:port）
func (s *ServiceInfo) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.Address, s.Port)
}
