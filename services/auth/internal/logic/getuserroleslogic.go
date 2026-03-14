package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/auth/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/auth/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetUserRolesLogic 负责 RBAC 下查询用户系统级角色的业务逻辑。
type GetUserRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewGetUserRolesLogic 构造一个查询用户角色逻辑实例。
func NewGetUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserRolesLogic {
	return &GetUserRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserRoles 查询指定用户的系统级角色列表。
// - 要求 user_id 非空；
// - 调用 RBACMod.GetUserRoles 获取角色并返回。
func (l *GetUserRolesLogic) GetUserRoles(in *pb.GetUserRolesRequest) (*pb.GetUserRolesResponse, error) {
	// 1. user_id 校验。
	if in.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// 2. 查角色并返回。
	roles, err := l.svcCtx.RBACMod.GetUserRoles(in.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get user roles failed: %v", err)
	}

	return &pb.GetUserRolesResponse{
		Roles: roles,
	}, nil
}
