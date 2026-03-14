package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/user/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/user/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// GetUserByUsernameLogic 按用户名查询用户资料（返回与 GetUser 相同的 User，id 为 10 位）
type GetUserByUsernameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByUsernameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByUsernameLogic {
	return &GetUserByUsernameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserByUsernameLogic) GetUserByUsername(in *pb.GetUserByUsernameRequest) (*pb.GetUserResponse, error) {
	if in.GetUsername() == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	u, err := l.svcCtx.UserMod.FindByUsername(in.GetUsername())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "query user failed: %v", err)
	}
	// 复用 GetUser 逻辑：按 id 查 profile 并返回
	return NewGetUserLogic(l.ctx, l.svcCtx).GetUser(&pb.GetUserRequest{Id: u.ID})
}
