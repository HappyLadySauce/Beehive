package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/message/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/message/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetLastMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLastMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLastMessagesLogic {
	return &GetLastMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLastMessagesLogic) GetLastMessages(in *pb.GetLastMessagesRequest) (*pb.GetLastMessagesResponse, error) {
	ids := in.GetConversationIds()
	if len(ids) == 0 {
		return &pb.GetLastMessagesResponse{LastMessages: nil}, nil
	}
	m, err := l.svcCtx.Msg.GetLastByConversations(ids)
	if err != nil {
		l.Errorf("get last messages failed: %v", err)
		return nil, status.Errorf(codes.Internal, "get last messages failed: %v", err)
	}
	lastMessages := make(map[string]*pb.MessageRecord)
	for k, msg := range m {
		lastMessages[k] = &pb.MessageRecord{
			ServerMsgId:    msg.ServerMsgID,
			ConversationId: msg.ConversationID,
			FromUserId:     msg.FromUserID,
			ToUserId:       msg.ToUserID,
			Body:           &pb.MessageBody{Type: msg.BodyType, Text: msg.BodyText},
			ServerTime:     msg.ServerTime,
		}
	}
	return &pb.GetLastMessagesResponse{LastMessages: lastMessages}, nil
}
