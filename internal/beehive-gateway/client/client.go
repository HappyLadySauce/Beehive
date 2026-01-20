package client

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"k8s.io/klog/v2"

	"github.com/HappyLadySauce/Beehive/internal/beehive-gateway/config"
	"github.com/HappyLadySauce/Beehive/internal/pkg/registry"
	authpb "github.com/HappyLadySauce/Beehive/pkg/api/proto/auth/v1"
	messagepb "github.com/HappyLadySauce/Beehive/pkg/api/proto/message/v1"
	presencepb "github.com/HappyLadySauce/Beehive/pkg/api/proto/presence/v1"
	userpb "github.com/HappyLadySauce/Beehive/pkg/api/proto/user/v1"
)

// Client gRPC 客户端管理器
type Client struct {
	authConn        *grpc.ClientConn
	userConn        *grpc.ClientConn
	messageConn     *grpc.ClientConn
	presenceConn    *grpc.ClientConn
	authService     authpb.AuthServiceClient
	userService     userpb.UserServiceClient
	messageService  messagepb.MessageServiceClient
	presenceService presencepb.PresenceServiceClient
	registry        *registry.Registry
}

// NewClient 创建新的 gRPC 客户端管理器
func NewClient(cfg *config.Config) (*Client, error) {
	client := &Client{}

	// 检查是否需要 etcd 服务发现
	useEtcdDiscovery := strings.HasPrefix(cfg.Services.AuthServiceAddr, "etcd://") ||
		strings.HasPrefix(cfg.Services.UserServiceAddr, "etcd://") ||
		strings.HasPrefix(cfg.Services.MessageServiceAddr, "etcd://") ||
		strings.HasPrefix(cfg.Services.PresenceServiceAddr, "etcd://")

	if useEtcdDiscovery {
		if len(cfg.Etcd.Endpoints) == 0 || cfg.Etcd.Endpoints[0] == "" {
			return nil, fmt.Errorf("etcd endpoints not configured but etcd:// address specified")
		}

		// 创建 etcd 注册中心
		serviceRegistry, err := registry.NewRegistry(
			cfg.Etcd.Endpoints,
			cfg.Etcd.DialTimeout,
			cfg.Etcd.Username,
			cfg.Etcd.Password,
			cfg.Etcd.Prefix,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create etcd registry: %w", err)
		}
		client.registry = serviceRegistry

		// 注册 etcd 解析器
		resolver.Register(registry.NewResolverBuilder(serviceRegistry))
	}

	// 连接 Auth Service
	authConn, err := connectToService(cfg.Services.AuthServiceAddr, "Auth Service")
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to Auth Service: %w", err)
	}
	client.authConn = authConn
	client.authService = authpb.NewAuthServiceClient(authConn)

	// 连接 User Service
	userConn, err := connectToService(cfg.Services.UserServiceAddr, "User Service")
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to User Service: %w", err)
	}
	client.userConn = userConn
	client.userService = userpb.NewUserServiceClient(userConn)

	// 连接 Message Service
	messageConn, err := connectToService(cfg.Services.MessageServiceAddr, "Message Service")
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to Message Service: %w", err)
	}
	client.messageConn = messageConn
	client.messageService = messagepb.NewMessageServiceClient(messageConn)

	// 连接 Presence Service
	presenceConn, err := connectToService(cfg.Services.PresenceServiceAddr, "Presence Service")
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to Presence Service: %w", err)
	}
	client.presenceConn = presenceConn
	client.presenceService = presencepb.NewPresenceServiceClient(presenceConn)

	klog.Info("All gRPC clients connected successfully")
	return client, nil
}

// connectToService 连接到指定的服务
func connectToService(addr, serviceName string) (*grpc.ClientConn, error) {
	var target string
	if strings.HasPrefix(addr, "etcd://") {
		serviceName := strings.TrimPrefix(addr, "etcd://")
		target = fmt.Sprintf("%s://%s", registry.Scheme, serviceName)
	} else {
		target = addr
	}

	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(addr, "etcd://") {
		klog.Infof("Connected to %s via etcd service discovery: %s", serviceName, strings.TrimPrefix(addr, "etcd://"))
	} else {
		klog.Infof("Connected to %s at %s", serviceName, addr)
	}

	return conn, nil
}

// Close 关闭所有连接
func (c *Client) Close() error {
	var errs []error

	if c.authConn != nil {
		if err := c.authConn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close Auth Service connection: %w", err))
		}
	}

	if c.userConn != nil {
		if err := c.userConn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close User Service connection: %w", err))
		}
	}

	if c.messageConn != nil {
		if err := c.messageConn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close Message Service connection: %w", err))
		}
	}

	if c.presenceConn != nil {
		if err := c.presenceConn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close Presence Service connection: %w", err))
		}
	}

	if c.registry != nil {
		if err := c.registry.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close etcd registry: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %v", errs)
	}

	return nil
}

// AuthService 返回 Auth Service 客户端
func (c *Client) AuthService() authpb.AuthServiceClient {
	return c.authService
}

// UserService 返回 User Service 客户端
func (c *Client) UserService() userpb.UserServiceClient {
	return c.userService
}

// MessageService 返回 Message Service 客户端
func (c *Client) MessageService() messagepb.MessageServiceClient {
	return c.messageService
}

// PresenceService 返回 Presence Service 客户端
func (c *Client) PresenceService() presencepb.PresenceServiceClient {
	return c.presenceService
}

// ValidateToken 验证 Token（便捷方法）
func (c *Client) ValidateToken(ctx context.Context, token string) (*authpb.ValidateTokenResponse, error) {
	req := &authpb.ValidateTokenRequest{
		Token: token,
	}
	return c.authService.ValidateToken(ctx, req)
}
