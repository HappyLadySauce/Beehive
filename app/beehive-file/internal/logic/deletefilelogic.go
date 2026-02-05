package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-file/file"
	"github.com/HappyLadySauce/Beehive/app/beehive-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除文件
func (l *DeleteFileLogic) DeleteFile(in *file.DeleteFileRequest) (*file.CommonResponse, error) {
	// todo: add your logic here and delete this line

	return &file.CommonResponse{}, nil
}
