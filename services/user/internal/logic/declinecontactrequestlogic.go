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

type DeclineContactRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeclineContactRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeclineContactRequestLogic {
	return &DeclineContactRequestLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DeclineContactRequestLogic) DeclineContactRequest(in *pb.DeclineContactRequestRequest) (*pb.DeclineContactRequestResponse, error) {
	if in.GetUserId() == "" || in.GetRequestId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and request_id are required")
	}
	err := l.svcCtx.ContactRequestMod.Decline(in.GetRequestId(), in.GetUserId())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "request not found or not pending")
		}
		l.Errorf("decline contact request failed: %v", err)
		return nil, status.Errorf(codes.Internal, "decline contact request failed: %v", err)
	}
	return &pb.DeclineContactRequestResponse{}, nil
}
