package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/message/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/message/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHistoryLogic {
	return &GetHistoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetHistoryLogic) GetHistory(in *pb.GetHistoryRequest) (*pb.GetHistoryResponse, error) {
	if in.GetConversationId() == "" {
		return nil, status.Error(codes.InvalidArgument, "conversation_id is required")
	}
	limit := int(in.GetLimit())
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	list, err := l.svcCtx.Msg.GetHistory(in.GetConversationId(), in.GetBeforeTime(), limit+1)
	if err != nil {
		l.Errorf("get history failed: %v", err)
		return nil, status.Errorf(codes.Internal, "get history failed: %v", err)
	}
	hasMore := len(list) > limit
	if hasMore {
		list = list[:limit]
	}
	items := make([]*pb.MessageRecord, 0, len(list))
	for _, m := range list {
		toUserID := ""
		if m.ToUserID != nil {
			toUserID = *m.ToUserID
		}
		items = append(items, &pb.MessageRecord{
			ServerMsgId:    m.ServerMsgID,
			ConversationId: m.ConversationID,
			FromUserId:     m.FromUserID,
			ToUserId:       toUserID,
			Body:           &pb.MessageBody{Type: m.BodyType, Text: m.BodyText},
			ServerTime:     m.ServerTime,
		})
	}
	return &pb.GetHistoryResponse{Items: items, HasMore: hasMore}, nil
}
