package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/message/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/message/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUnreadCountsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUnreadCountsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUnreadCountsLogic {
	return &GetUnreadCountsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUnreadCountsLogic) GetUnreadCounts(in *pb.GetUnreadCountsRequest) (*pb.GetUnreadCountsResponse, error) {
	if in.GetUserId() == "" {
		return &pb.GetUnreadCountsResponse{Counts: make(map[string]int32)}, nil
	}
	convIDs := in.GetConversationIds()
	if len(convIDs) == 0 {
		return &pb.GetUnreadCountsResponse{Counts: make(map[string]int32)}, nil
	}
	lastReadTimes, err := l.svcCtx.Read.GetLastReadTimes(in.GetUserId(), convIDs)
	if err != nil {
		l.Errorf("get last read times failed: %v", err)
		return nil, err
	}
	counts := make(map[string]int32, len(convIDs))
	for _, cid := range convIDs {
		lastRead := lastReadTimes[cid]
		n, err := l.svcCtx.Msg.CountUnread(cid, in.GetUserId(), lastRead)
		if err != nil {
			l.Errorf("count unread for conversation %s failed: %v", cid, err)
			counts[cid] = 0
			continue
		}
		counts[cid] = int32(n)
	}
	return &pb.GetUnreadCountsResponse{Counts: counts}, nil
}
