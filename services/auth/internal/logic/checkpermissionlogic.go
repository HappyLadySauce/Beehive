package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/auth/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/auth/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CheckPermissionLogic 负责 RBAC 下检查用户是否具备某权限的业务逻辑。
type CheckPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewCheckPermissionLogic 构造一个权限检查逻辑实例。
func NewCheckPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckPermissionLogic {
	return &CheckPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CheckPermission 检查指定用户是否具备目标权限。
// - 要求 user_id、permission 均非空；
// - 拉取用户权限列表，判断是否包含目标权限，返回 Allowed。
func (l *CheckPermissionLogic) CheckPermission(in *pb.CheckPermissionRequest) (*pb.CheckPermissionResponse, error) {
	// 1. 参数校验。
	if in.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if in.GetPermission() == "" {
		return nil, status.Error(codes.InvalidArgument, "permission is required")
	}

	// 2. 拉取用户权限列表。
	perms, err := l.svcCtx.RBACMod.GetUserPermissions(in.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get user permissions failed: %v", err)
	}

	// 3. 判断是否包含目标权限并返回。
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
