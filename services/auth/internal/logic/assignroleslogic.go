package logic

import (
	"context"
	"errors"

	"github.com/HappyLadySauce/Beehive/services/auth/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/auth/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAssignRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignRolesLogic {
	return &AssignRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RBAC（可选，对内/管理用途）：为用户设置角色
func (l *AssignRolesLogic) AssignRoles(in *pb.AssignRolesRequest) (*pb.AssignRolesResponse, error) {
	if in.GetUserId() == "" {
		return nil, errors.New("user_id is empty")
	}

	if err := l.svcCtx.RBACMod.ReplaceUserRoles(l.ctx, in.GetUserId(), in.GetRoles()); err != nil {
		return nil, err
	}

	return &pb.AssignRolesResponse{}, nil
}
