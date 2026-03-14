// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/HappyLadySauce/Beehive/services/auth/authservice"
	"github.com/HappyLadySauce/Beehive/services/conversation/conversationservice"
	"github.com/HappyLadySauce/Beehive/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive/services/gateway/internal/push"
	"github.com/HappyLadySauce/Beehive/services/gateway/internal/ratelimit"
	"github.com/HappyLadySauce/Beehive/services/gateway/internal/ws"
	"github.com/HappyLadySauce/Beehive/services/message/messageservice"
	"github.com/HappyLadySauce/Beehive/services/presence/presenceservice"
	"github.com/HappyLadySauce/Beehive/services/user/userservice"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	Hub    *ws.Hub

	AuthSvc         authservice.AuthService
	PresenceSvc      presenceservice.PresenceService
	UserSvc          userservice.UserService
	ConversationSvc  conversationservice.ConversationService  // 可选，未配置时为 nil
	MessageSvc       messageservice.MessageService             // 可选，未配置时为 nil
	PushConsumer     *push.Consumer                             // 可选，未配置 RabbitMQ 时为 nil
	MessageSendLimit *ratelimit.MessageSendLimiter              // 可选，未配置 Redis/限流时为 nil
}

func NewServiceContext(c config.Config) *ServiceContext {
	authCli := zrpc.MustNewClient(c.AuthRpcConf)
	presenceCli := zrpc.MustNewClient(c.PresenceRpcConf)
	userCli := zrpc.MustNewClient(c.UserRpcConf)

	ctx := &ServiceContext{
		Config:      c,
		Hub:         ws.NewHub(c.GatewayID),
		AuthSvc:     authservice.NewAuthService(authCli),
		PresenceSvc: presenceservice.NewPresenceService(presenceCli),
		UserSvc:     userservice.NewUserService(userCli),
	}
	if c.AuthRpcConfigured() {
		ctx.AuthSvc = authservice.NewAuthService(authCli)
	}
	if c.PresenceRpcConfigured() {
		ctx.PresenceSvc = presenceservice.NewPresenceService(presenceCli)
	}
	if c.UserRpcConfigured() {
		ctx.UserSvc = userservice.NewUserService(userCli)
	}
	if c.ConversationRpcConfigured() {
		convCli := zrpc.MustNewClient(c.ConversationRpcConf)
		ctx.ConversationSvc = conversationservice.NewConversationService(convCli)
	}
	if c.MessageRpcConfigured() {
		msgCli := zrpc.MustNewClient(c.MessageRpcConf)
		ctx.MessageSvc = messageservice.NewMessageService(msgCli)
	}
	if c.PushConsumerConfigured() && ctx.ConversationSvc != nil && ctx.PresenceSvc != nil {
		consumer, err := push.NewConsumer(c, ctx.Hub, ctx.ConversationSvc, ctx.PresenceSvc)
		if err != nil {
			panic("push consumer: " + err.Error())
		}
		ctx.PushConsumer = consumer
	}
	if c.RateLimitConfigured() {
		rdb := redis.NewClient(&redis.Options{
			Addr:     c.RedisAddr,
			Password: c.RedisPassword,
			DB:       c.RedisDB,
		})
		ctx.MessageSendLimit = ratelimit.NewMessageSendLimiter(rdb, c.RateLimitMessageSendPerMinute)
	}
	return ctx
}