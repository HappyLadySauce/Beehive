package client

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"k8s.io/klog/v2"

	pb "github.com/HappyLadySauce/Beehive/api/proto/user/v1"
	"github.com/HappyLadySauce/Beehive/internal/beehive-auth/config"
	"github.com/HappyLadySauce/Beehive/internal/pkg/registry"
)

// Client User Service gRPC 客户端
type Client struct {
	conn        *grpc.ClientConn
	userService pb.UserServiceClient
	registry    *registry.Registry
}

// NewClient 创建新的 User Service 客户端
func NewClient(cfg *config.Config) (*Client, error) {
	var conn *grpc.ClientConn
	var err error
	var serviceRegistry *registry.Registry

	addr := cfg.Services.UserServiceAddr

	// 检查是否使用 etcd 服务发现
	useEtcdDiscovery := strings.HasPrefix(addr, "etcd://")
	if useEtcdDiscovery {
		// 需要 etcd 配置
		if len(cfg.Etcd.Endpoints) == 0 || cfg.Etcd.Endpoints[0] == "" {
			return nil, fmt.Errorf("etcd endpoints not configured but etcd:// address specified")
		}

		// 创建 etcd 注册中心
		serviceRegistry, err = registry.NewRegistry(
			cfg.Etcd.Endpoints,
			cfg.Etcd.DialTimeout,
			cfg.Etcd.Username,
			cfg.Etcd.Password,
			cfg.Etcd.Prefix,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create etcd registry: %w", err)
		}

		// 注册 etcd 解析器
		resolver.Register(registry.NewResolverBuilder(serviceRegistry))

		// 提取服务名
		serviceName := strings.TrimPrefix(addr, "etcd://")
		target := fmt.Sprintf("%s://%s", registry.Scheme, serviceName)

		// 使用 etcd 服务发现连接
		conn, err = grpc.NewClient(
			target,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			serviceRegistry.Close()
			return nil, fmt.Errorf("failed to connect to User Service via etcd: %w", err)
		}
		klog.Infof("Connected to User Service via etcd service discovery: %s", serviceName)
	} else {
		// 使用直接地址连接
		conn, err = grpc.NewClient(
			addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to User Service: %w", err)
		}
		klog.Infof("Connected to User Service at %s", addr)
	}

	client := &Client{
		conn:        conn,
		userService: pb.NewUserServiceClient(conn),
		registry:    serviceRegistry,
	}

	return client, nil
}

// Close 关闭连接
func (c *Client) Close() error {
	err := c.conn.Close()
	if c.registry != nil {
		if closeErr := c.registry.Close(); closeErr != nil {
			klog.Warningf("Failed to close etcd registry: %v", closeErr)
		}
	}
	return err
}

// GetUserByID 获取用户信息（包含密码哈希和盐值）
func (c *Client) GetUserByID(ctx context.Context, userID string) (*pb.GetUserByIDResponse, error) {
	req := &pb.GetUserByIDRequest{
		Id: userID,
	}
	return c.userService.GetUserByID(ctx, req)
}

// GetUser 获取用户信息（不包含敏感信息）
func (c *Client) GetUser(ctx context.Context, userID string) (*pb.GetUserResponse, error) {
	req := &pb.GetUserRequest{
		Id: userID,
	}
	return c.userService.GetUser(ctx, req)
}
