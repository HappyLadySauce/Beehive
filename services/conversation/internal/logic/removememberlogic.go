package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/conversation/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/conversation/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveMemberLogic {
	return &RemoveMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveMemberLogic) RemoveMember(in *pb_conversationpb.RemoveMemberRequest) (*pb_conversationpb.RemoveMemberResponse, error) {
	// todo: add your logic here and delete this line

	return &pb_conversationpb.RemoveMemberResponse{}, nil
}
