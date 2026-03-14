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

type AcceptContactRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAcceptContactRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AcceptContactRequestLogic {
	return &AcceptContactRequestLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AcceptContactRequestLogic) AcceptContactRequest(in *pb.AcceptContactRequestRequest) (*pb.AcceptContactRequestResponse, error) {
	if in.GetUserId() == "" || in.GetRequestId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and request_id are required")
	}
	err := l.svcCtx.ContactRequestMod.Accept(in.GetRequestId(), in.GetUserId())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "request not found or not pending")
		}
		l.Errorf("accept contact request failed: %v", err)
		return nil, status.Errorf(codes.Internal, "accept contact request failed: %v", err)
	}
	return &pb.AcceptContactRequestResponse{}, nil
}
