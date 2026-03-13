package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/HappyLadySauce/Beehive/services/auth/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/auth/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginRequest) (*pb.LoginResponse, error) {
	if in.GetUsername() == "" || in.GetPassword() == "" {
		return nil, errors.New("username or password is empty")
	}

	user, err := l.svcCtx.UserMod.FindByUsername(in.GetUsername())
	if err != nil {
		return nil, fmt.Errorf("user not found or db error: %w", err)
	}

	if user.Status != "normal" {
		return nil, errors.New("user status is not normal")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.GetPassword())); err != nil {
		return nil, errors.New("invalid credentials")
	}

	roles, err := l.svcCtx.RBACMod.GetUserRoles(user.ID)
	if err != nil {
		return nil, fmt.Errorf("get user roles failed: %w", err)
	}

	accessTTL := tokenTTLSeconds(l.svcCtx.Config.AccessTokenTTLSeconds, 3600)
	refreshTTL := tokenTTLSeconds(l.svcCtx.Config.RefreshTokenTTLSeconds, 30*24*3600)

	accessToken := uuid.NewString()
	refreshToken := uuid.NewString()

	if err := storeToken(l.ctx, l.svcCtx.Redis, accessToken, user.ID, roles, time.Duration(accessTTL)*time.Second); err != nil {
		return nil, fmt.Errorf("store access token failed: %w", err)
	}
	if err := storeToken(l.ctx, l.svcCtx.Redis, refreshToken, user.ID, roles, time.Duration(refreshTTL)*time.Second); err != nil {
		return nil, fmt.Errorf("store refresh token failed: %w", err)
	}

	return &pb.LoginResponse{
		UserId:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessTTL),
	}, nil
}
