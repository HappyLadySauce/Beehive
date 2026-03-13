package logic

import (
	"context"
	"errors"

	"github.com/HappyLadySauce/Beehive/services/auth/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/auth/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserRolesLogic {
	return &GetUserRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RBAC：查询用户系统级角色
func (l *GetUserRolesLogic) GetUserRoles(in *pb.GetUserRolesRequest) (*pb.GetUserRolesResponse, error) {
	if in.GetUserId() == "" {
		return nil, errors.New("user_id is empty")
	}

	roles, err := l.svcCtx.RBACMod.GetUserRoles(in.GetUserId())
	if err != nil {
		return nil, err
	}

	return &pb.GetUserRolesResponse{
		Roles: roles,
	}, nil
}
