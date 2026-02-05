package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-file/file"
	"github.com/HappyLadySauce/Beehive/app/beehive-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 上传文件
func (l *UploadFileLogic) UploadFile(in *file.UploadFileRequest) (*file.UploadFileResponse, error) {
	// todo: add your logic here and delete this line

	return &file.UploadFileResponse{}, nil
}
