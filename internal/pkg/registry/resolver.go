package registry

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"google.golang.org/grpc/resolver"
	"k8s.io/klog/v2"
)

const (
	// Scheme etcd 解析器方案名
	Scheme = "etcd"
)

var (
	// roundRobinIndex 轮询索引（每个服务名独立）
	roundRobinIndex = make(map[string]int)
	roundRobinMutex sync.Mutex
)

// etcdResolver etcd 服务发现解析器
type etcdResolver struct {
	target      resolver.Target
	cc          resolver.ClientConn
	registry    *Registry
	serviceName string
	ctx         context.Context
	cancel      context.CancelFunc
}

// etcdResolverBuilder etcd 解析器构建器
type etcdResolverBuilder struct {
	registry *Registry
}

// NewResolverBuilder 创建 etcd 解析器构建器
func NewResolverBuilder(registry *Registry) resolver.Builder {
	return &etcdResolverBuilder{
		registry: registry,
	}
}

// Build 构建解析器
func (b *etcdResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// 从 target 获取服务名
	// 格式: etcd://beehive-auth 或 etcd:///beehive-auth
	serviceName := target.Endpoint()
	if serviceName == "" {
		// 尝试从 URL Path 获取
		serviceName = target.URL.Path
	}
	// 移除前导斜杠
	if strings.HasPrefix(serviceName, "/") {
		serviceName = serviceName[1:]
	}

	ctx, cancel := context.WithCancel(context.Background())
	r := &etcdResolver{
		target:      target,
		cc:          cc,
		registry:    b.registry,
		serviceName: serviceName,
		ctx:         ctx,
		cancel:      cancel,
	}

	// 初始解析
	if err := r.resolve(); err != nil {
		return nil, err
	}

	// 监听服务变化
	go r.watch()

	return r, nil
}

// Scheme 返回解析器方案名
func (b *etcdResolverBuilder) Scheme() string {
	return Scheme
}

// resolve 解析服务地址
func (r *etcdResolver) resolve() error {
	services, err := r.registry.Discover(r.serviceName)
	if err != nil {
		return fmt.Errorf("failed to discover service %s: %w", r.serviceName, err)
	}

	if len(services) == 0 {
		klog.Warningf("No instances found for service: %s", r.serviceName)
		return nil
	}

	// 构建所有可用服务地址列表（gRPC 会使用负载均衡器选择）
	var addresses []resolver.Address
	for _, svc := range services {
		addresses = append(addresses, resolver.Address{
			Addr: svc.GetAddress(),
		})
	}

	klog.V(4).Infof("Resolved service %s to %d instances", r.serviceName, len(addresses))

	// 更新 gRPC 客户端连接（gRPC 会使用默认的轮询负载均衡器）
	state := resolver.State{
		Addresses: addresses,
	}

	if err := r.cc.UpdateState(state); err != nil {
		return fmt.Errorf("failed to update resolver state: %w", err)
	}

	return nil
}

// selectService 使用轮询策略选择服务
func (r *etcdResolver) selectService(services []*ServiceInfo) *ServiceInfo {
	roundRobinMutex.Lock()
	defer roundRobinMutex.Unlock()

	if len(services) == 0 {
		return nil
	}

	index := roundRobinIndex[r.serviceName]
	service := services[index%len(services)]
	roundRobinIndex[r.serviceName] = (index + 1) % len(services)

	return service
}

// watch 监听服务变化
func (r *etcdResolver) watch() {
	err := r.registry.Watch(r.serviceName, func(services []*ServiceInfo) {
		klog.Infof("Service instances changed for %s: %d instances", r.serviceName, len(services))
		// 重置轮询索引
		roundRobinMutex.Lock()
		roundRobinIndex[r.serviceName] = 0
		roundRobinMutex.Unlock()
		// 重新解析
		if err := r.resolve(); err != nil {
			klog.Errorf("Failed to resolve after watch: %v", err)
		}
	})
	if err != nil {
		klog.Errorf("Failed to watch service %s: %v", r.serviceName, err)
	}
}

// ResolveNow 立即重新解析
func (r *etcdResolver) ResolveNow(resolver.ResolveNowOptions) {
	if err := r.resolve(); err != nil {
		klog.Errorf("Failed to resolve now: %v", err)
	}
}

// Close 关闭解析器
func (r *etcdResolver) Close() {
	r.cancel()
}
