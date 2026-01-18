package registry

import (
	"context"
	"fmt"
	"path"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"k8s.io/klog/v2"
)

// Registry 服务注册中心
type Registry struct {
	client  *clientv3.Client
	prefix  string
	leaseID clientv3.LeaseID
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewRegistry 创建新的服务注册中心
func NewRegistry(endpoints []string, dialTimeout time.Duration, username, password, prefix string) (*Registry, error) {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	}

	if username != "" {
		cfg.Username = username
		cfg.Password = password
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Registry{
		client: client,
		prefix: prefix,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Close 关闭注册中心
func (r *Registry) Close() error {
	r.cancel()
	return r.client.Close()
}

// Register 注册服务到 etcd
func (r *Registry) Register(serviceInfo *ServiceInfo, ttl int64) error {
	// 1. 创建租约
	leaseResp, err := r.client.Grant(r.ctx, ttl)
	if err != nil {
		return fmt.Errorf("failed to create lease: %w", err)
	}
	r.leaseID = leaseResp.ID

	// 2. 构建 key
	key := r.getServiceKey(serviceInfo.ServiceName, serviceInfo.InstanceID)

	// 3. 序列化服务信息
	value, err := serviceInfo.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal service info: %w", err)
	}

	// 4. 注册服务
	_, err = r.client.Put(r.ctx, key, value, clientv3.WithLease(r.leaseID))
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	klog.Infof("Service registered: %s -> %s:%d", key, serviceInfo.Address, serviceInfo.Port)

	// 5. 启动租约续约
	go r.keepAlive(ttl)

	return nil
}

// Deregister 注销服务
func (r *Registry) Deregister(serviceName, instanceID string) error {
	key := r.getServiceKey(serviceName, instanceID)
	_, err := r.client.Delete(r.ctx, key)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}

	// 撤销租约
	if r.leaseID != 0 {
		_, _ = r.client.Revoke(r.ctx, r.leaseID)
	}

	klog.Infof("Service deregistered: %s", key)
	return nil
}

// Discover 发现服务，返回所有可用实例
func (r *Registry) Discover(serviceName string) ([]*ServiceInfo, error) {
	keyPrefix := r.getServicePrefix(serviceName)
	resp, err := r.client.Get(r.ctx, keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to discover service: %w", err)
	}

	var services []*ServiceInfo
	for _, kv := range resp.Kvs {
		var info ServiceInfo
		if err := info.FromJSON(string(kv.Value)); err != nil {
			klog.Warningf("Failed to unmarshal service info: %v", err)
			continue
		}
		services = append(services, &info)
	}

	return services, nil
}

// Watch 监听服务变化
func (r *Registry) Watch(serviceName string, callback func([]*ServiceInfo)) error {
	keyPrefix := r.getServicePrefix(serviceName)
	watchChan := r.client.Watch(r.ctx, keyPrefix, clientv3.WithPrefix())

	go func() {
		for range watchChan {
			services, err := r.Discover(serviceName)
			if err != nil {
				klog.Errorf("Failed to discover services after watch event: %v", err)
				continue
			}
			callback(services)
		}
	}()

	return nil
}

// keepAlive 保持租约活跃
func (r *Registry) keepAlive(ttl int64) {
	ch, err := r.client.KeepAlive(r.ctx, r.leaseID)
	if err != nil {
		klog.Errorf("Failed to keep alive lease: %v", err)
		return
	}

	for ka := range ch {
		if ka == nil {
			klog.Warning("Lease keep alive channel closed")
			return
		}
		klog.V(4).Infof("Lease kept alive: ID=%d, TTL=%d", ka.ID, ka.TTL)
	}
}

// getServiceKey 获取服务 key
func (r *Registry) getServiceKey(serviceName, instanceID string) string {
	return path.Join(r.prefix, serviceName, instanceID)
}

// getServicePrefix 获取服务前缀
func (r *Registry) getServicePrefix(serviceName string) string {
	return path.Join(r.prefix, serviceName) + "/"
}
