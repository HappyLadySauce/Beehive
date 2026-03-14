// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/adminapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUsersLogic {
	return &ListUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUsersLogic) ListUsers(req *types.ListUsersReq) (resp *types.ListUsersResp, err error) {
	page, pageSize := req.Page, req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	return &types.ListUsersResp{
		Code:    0,
		Message: "ok",
		Data: types.ListUsersData{
			Items:    []types.UserItem{},
			Page:     page,
			PageSize: pageSize,
			Total:    0,
		},
	}, nil
}
