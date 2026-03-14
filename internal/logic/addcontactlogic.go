package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/user/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddContactLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddContactLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddContactLogic {
	return &AddContactLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 联系人
func (l *AddContactLogic) AddContact(in *pb.AddContactRequest) (*pb.AddContactResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.AddContactResponse{}, nil
}
