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

type CreateContactRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateContactRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateContactRequestLogic {
	return &CreateContactRequestLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *CreateContactRequestLogic) CreateContactRequest(in *pb.CreateContactRequestRequest) (*pb.CreateContactRequestResponse, error) {
	if in.GetFromUserId() == "" || in.GetToUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "from_user_id and to_user_id are required")
	}
	if in.GetFromUserId() == in.GetToUserId() {
		return nil, status.Error(codes.InvalidArgument, "cannot send request to self")
	}
	// 校验对方用户存在
	_, err := l.svcCtx.UserMod.FindByID(in.GetToUserId())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "to user not found")
		}
		return nil, status.Errorf(codes.Internal, "query user failed: %v", err)
	}
	// 已是好友则不再发申请
	ok, err := l.svcCtx.ContactMod.Exists(in.GetFromUserId(), in.GetToUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "check contact failed: %v", err)
	}
	if ok {
		return nil, status.Error(codes.AlreadyExists, "already contacts")
	}
	req, err := l.svcCtx.ContactRequestMod.CreateOrReapply(in.GetFromUserId(), in.GetToUserId(), in.GetMessage())
	if err != nil {
		l.Errorf("create contact request failed: %v", err)
		return nil, status.Errorf(codes.Internal, "create contact request failed: %v", err)
	}
	return &pb.CreateContactRequestResponse{RequestId: req.ID}, nil
}
