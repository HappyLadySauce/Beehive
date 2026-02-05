package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-file/file"
	"github.com/HappyLadySauce/Beehive/app/beehive-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckFileExistsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckFileExistsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckFileExistsLogic {
	return &CheckFileExistsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 检查文件是否存在（去重）
func (l *CheckFileExistsLogic) CheckFileExists(in *file.CheckFileExistsRequest) (*file.CheckFileExistsResponse, error) {
	// todo: add your logic here and delete this line

	return &file.CheckFileExistsResponse{}, nil
}
