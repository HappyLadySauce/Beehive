package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-file/file"
	"github.com/HappyLadySauce/Beehive/app/beehive-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompleteMultipartUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCompleteMultipartUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteMultipartUploadLogic {
	return &CompleteMultipartUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 完成分片上传
func (l *CompleteMultipartUploadLogic) CompleteMultipartUpload(in *file.CompleteMultipartUploadRequest) (*file.UploadFileResponse, error) {
	// todo: add your logic here and delete this line

	return &file.UploadFileResponse{}, nil
}
