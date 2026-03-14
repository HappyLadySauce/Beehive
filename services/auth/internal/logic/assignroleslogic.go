package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/auth/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/auth/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AssignRolesLogic 负责 RBAC 下为用户设置/替换角色的业务逻辑，供管理或内部用途。
type AssignRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewAssignRolesLogic 构造一个分配角色逻辑实例。
func NewAssignRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignRolesLogic {
	return &AssignRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AssignRoles 为指定用户设置角色列表（全量替换）。
// - 要求 user_id 非空；
// - 使用 ReplaceUserRoles 将用户角色替换为请求中的 roles 列表。
func (l *AssignRolesLogic) AssignRoles(in *pb.AssignRolesRequest) (*pb.AssignRolesResponse, error) {
	// 1. user_id 校验。
	if in.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// 2. 替换角色并返回。
	if err := l.svcCtx.RBACMod.ReplaceUserRoles(l.ctx, in.GetUserId(), in.GetRoles()); err != nil {
		return nil, status.Errorf(codes.Internal, "assign roles failed: %v", err)
	}

	return &pb.AssignRolesResponse{}, nil
}
