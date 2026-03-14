// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/config"
	"github.com/HappyLadySauce/Beehive/services/auth/authservice"
	"github.com/HappyLadySauce/Beehive/services/conversation/conversationservice"
	"github.com/HappyLadySauce/Beehive/services/message/messageservice"
	"github.com/HappyLadySauce/Beehive/services/presence/presenceservice"
	"github.com/HappyLadySauce/Beehive/services/user/userservice"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	AuthSvc         authservice.AuthService
	UserSvc         userservice.UserService
	PresenceSvc     presenceservice.PresenceService
	MessageSvc      messageservice.MessageService
	ConversationSvc conversationservice.ConversationService
}

func NewServiceContext(c config.Config) *ServiceContext {
	authCli := zrpc.MustNewClient(c.AuthRpc)
	userCli := zrpc.MustNewClient(c.UserRpc)
	presenceCli := zrpc.MustNewClient(c.PresenceRpc)
	messageCli := zrpc.MustNewClient(c.MessageRpc)
	conversationCli := zrpc.MustNewClient(c.ConversationRpc)

	return &ServiceContext{
		Config:          c,
		AuthSvc:         authservice.NewAuthService(authCli),
		UserSvc:         userservice.NewUserService(userCli),
		PresenceSvc:     presenceservice.NewPresenceService(presenceCli),
		MessageSvc:      messageservice.NewMessageService(messageCli),
		ConversationSvc: conversationservice.NewConversationService(conversationCli),
	}
}
