package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-file/file"
	"github.com/HappyLadySauce/Beehive/app/beehive-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type InitMultipartUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInitMultipartUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitMultipartUploadLogic {
	return &InitMultipartUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 初始化分片上传
func (l *InitMultipartUploadLogic) InitMultipartUpload(in *file.InitMultipartUploadRequest) (*file.InitMultipartUploadResponse, error) {
	// todo: add your logic here and delete this line

	return &file.InitMultipartUploadResponse{}, nil
}
