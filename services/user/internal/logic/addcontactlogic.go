package logic

import (
	"context"
	"errors"

	"github.com/HappyLadySauce/Beehive/services/user/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/user/pb"
	"github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type AddContactLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddContactLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddContactLogic {
	return &AddContactLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AddContactLogic) AddContact(in *pb.AddContactRequest) (*pb.AddContactResponse, error) {
	if in.GetOwnerId() == "" || in.GetContactUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "owner_id and contact_user_id are required")
	}
	if in.GetOwnerId() == in.GetContactUserId() {
		return nil, status.Error(codes.InvalidArgument, "cannot add self as contact")
	}
	// 校验对方用户存在
	_, err := l.svcCtx.UserMod.FindByID(in.GetContactUserId())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "contact user not found")
		}
		return nil, status.Errorf(codes.Internal, "query user failed: %v", err)
	}
	err = l.svcCtx.ContactMod.Add(in.GetOwnerId(), in.GetContactUserId())
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return &pb.AddContactResponse{}, nil
		}
		return nil, status.Errorf(codes.Internal, "add contact failed: %v", err)
	}
	return &pb.AddContactResponse{}, nil
}
