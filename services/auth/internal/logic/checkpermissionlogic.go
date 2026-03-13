package logic

import (
	"context"
	"errors"

	"github.com/HappyLadySauce/Beehive/services/auth/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/auth/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckPermissionLogic {
	return &CheckPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RBAC：检查用户是否具备某个权限
func (l *CheckPermissionLogic) CheckPermission(in *pb.CheckPermissionRequest) (*pb.CheckPermissionResponse, error) {
	if in.GetUserId() == "" {
		return nil, errors.New("user_id is empty")
	}
	if in.GetPermission() == "" {
		return nil, errors.New("permission is empty")
	}

	perms, err := l.svcCtx.RBACMod.GetUserPermissions(in.GetUserId())
	if err != nil {
		return nil, err
	}

	allowed := false
	for _, p := range perms {
		if p == in.GetPermission() {
			allowed = true
			break
		}
	}

	return &pb.CheckPermissionResponse{
		Allowed: allowed,
	}, nil
}
