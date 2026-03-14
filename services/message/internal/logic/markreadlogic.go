package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/message/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/message/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type MarkReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkReadLogic {
	return &MarkReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MarkReadLogic) MarkRead(in *pb.MarkReadRequest) (*pb.MarkReadResponse, error) {
	if in.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if in.GetConversationId() == "" {
		return nil, status.Error(codes.InvalidArgument, "conversation_id is required")
	}
	if in.GetServerMsgId() == "" {
		return nil, status.Error(codes.InvalidArgument, "server_msg_id is required")
	}
	msg, err := l.svcCtx.Msg.GetByServerMsgID(in.GetConversationId(), in.GetServerMsgId())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "message not found")
		}
		l.Errorf("get message by server_msg_id failed: %v", err)
		return nil, status.Errorf(codes.Internal, "get message failed: %v", err)
	}
	if err := l.svcCtx.Read.UpsertLastRead(in.GetUserId(), in.GetConversationId(), msg.ServerTime); err != nil {
		l.Errorf("upsert last read failed: %v", err)
		return nil, status.Errorf(codes.Internal, "mark read failed: %v", err)
	}
	return &pb.MarkReadResponse{}, nil
}
