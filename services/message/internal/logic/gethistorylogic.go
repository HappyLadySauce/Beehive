package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/message/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/message/pb"

	"github.com/zeromicro/go-zero/core/logx"
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
	// todo: add your logic here and delete this line

	return &pb.GetHistoryResponse{}, nil
}
