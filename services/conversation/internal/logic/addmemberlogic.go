package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/conversation/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/conversation/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMemberLogic {
	return &AddMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddMemberLogic) AddMember(in *pb_conversationpb.AddMemberRequest) (*pb_conversationpb.AddMemberResponse, error) {
	// todo: add your logic here and delete this line

	return &pb_conversationpb.AddMemberResponse{}, nil
}
