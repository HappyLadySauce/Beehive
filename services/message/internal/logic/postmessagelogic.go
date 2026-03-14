package logic

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive/services/message/internal/model"
	"github.com/HappyLadySauce/Beehive/services/message/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/message/pb"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPostMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostMessageLogic {
	return &PostMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PostMessageLogic) PostMessage(in *pb.PostMessageRequest) (*pb.PostMessageResponse, error) {
	if in.GetConversationId() == "" {
		return nil, status.Error(codes.InvalidArgument, "conversation_id is required")
	}
	if in.GetFromUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "from_user_id is required")
	}
	body := in.GetBody()
	if body == nil {
		return nil, status.Error(codes.InvalidArgument, "body is required")
	}
	bodyType := body.GetType()
	if bodyType == "" {
		bodyType = "text"
	}
	serverMsgID := uuid.Must(uuid.NewUUID()).String()
	serverTime := time.Now().Unix()
	// toUserId 仅在点对点消息时使用；群聊/广播时应为 NULL，而不是空串，避免 uuid 列解析错误
	var toUserIDPtr *string
	if in.GetToUserId() != "" {
		v := in.GetToUserId()
		toUserIDPtr = &v
	}
	msg := &model.Message{
		ID:             uuid.Must(uuid.NewUUID()).String(),
		ServerMsgID:    serverMsgID,
		ClientMsgID:    in.GetClientMsgId(),
		ConversationID: in.GetConversationId(),
		FromUserID:     in.GetFromUserId(),
		ToUserID:       toUserIDPtr,
		BodyType:       bodyType,
		BodyText:       body.GetText(),
		ServerTime:     serverTime,
	}
	if err := l.svcCtx.Msg.Create(msg); err != nil {
		l.Errorf("create message failed: %v", err)
		return nil, status.Errorf(codes.Internal, "create message failed: %v", err)
	}
	event := map[string]interface{}{
		"serverMsgId":    serverMsgID,
		"clientMsgId":   in.GetClientMsgId(),
		"conversationId": in.GetConversationId(),
		"fromUserId":    in.GetFromUserId(),
		"toUserId":      in.GetToUserId(),
		"body":           map[string]string{"type": bodyType, "text": body.GetText()},
		"serverTime":     serverTime,
	}
	if err := l.svcCtx.MQ.PublishJSON(event); err != nil {
		l.Errorf("publish message.created failed: %v", err)
	}
	return &pb.PostMessageResponse{
		ServerMsgId:    serverMsgID,
		ConversationId: in.GetConversationId(),
		ServerTime:     serverTime,
	}, nil
}