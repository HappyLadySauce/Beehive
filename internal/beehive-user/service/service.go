package service

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"

	"github.com/HappyLadySauce/Beehive/internal/beehive-user/store"
	"github.com/HappyLadySauce/Beehive/internal/pkg/utils/id"
	v1 "github.com/HappyLadySauce/Beehive/pkg/api/proto/user/v1"
	userv1 "github.com/HappyLadySauce/Beehive/pkg/api/user/v1"
	"github.com/HappyLadySauce/Beehive/pkg/utils/passwd"
)

// Service User Service 实现
type Service struct {
	v1.UnimplementedUserServiceServer
	store       *store.Store
	idGenerator *id.Generator
}

// NewService 创建新的 User Service
func NewService(s *store.Store) *Service {
	return &Service{
		store:       s,
		idGenerator: id.NewGenerator(s.DB()),
	}
}

// Register 用户注册
func (s *Service) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	// 1. 验证请求参数
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if req.Nickname == "" {
		return nil, status.Error(codes.InvalidArgument, "nickname is required")
	}
	if len(req.Password) < 6 {
		return nil, status.Error(codes.InvalidArgument, "password must be at least 6 characters")
	}

	// 2. 检查邮箱是否已存在（排除已删除的用户）
	var existingUser userv1.User
	if err := s.store.DB().Where("email = ? AND deleted_at IS NULL", req.Email).First(&existingUser).Error; err == nil {
		klog.Warningf("Registration failed: email %s already exists", req.Email)
		return nil, status.Error(codes.AlreadyExists, "email already exists")
	}

	// 3. 生成盐值和密码哈希
	salt, err := passwd.GenerateSalt()
	if err != nil {
		klog.Errorf("Failed to generate salt: %v", err)
		return nil, status.Error(codes.Internal, "failed to generate salt")
	}

	passwordHash, err := passwd.HashPassword(req.Password, salt)
	if err != nil {
		klog.Errorf("Failed to hash password: %v", err)
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	// 4. 生成用户ID（10位数字，从1000000000开始递增）
	userID, err := s.idGenerator.GenerateUserID(ctx)
	if err != nil {
		klog.Errorf("Failed to generate user ID: %v", err)
		return nil, status.Error(codes.Internal, "failed to generate user ID")
	}

	// 5. 创建用户记录
	now := time.Now()
	user := &userv1.User{
		ID:           userID,
		Nickname:     req.Nickname,
		Avatar:       req.Avatar,
		Email:        req.Email,
		Description:  req.Description,
		Salt:         salt,
		PasswordHash: passwordHash,
		Level:        0, // 普通用户
		Status:       "offline",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.store.DB().Create(user).Error; err != nil {
		klog.Errorf("Failed to create user: %v", err)
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	klog.Infof("User registered successfully: id=%s, email=%s", userID, req.Email)

	return &v1.RegisterResponse{
		Id: userID,
	}, nil
}

// GetUser 根据用户ID获取用户信息（不包含敏感信息）
func (s *Service) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	var user userv1.User
	if err := s.store.DB().Where("id = ? AND deleted_at IS NULL", req.Id).First(&user).Error; err != nil {
		klog.Warningf("User not found: id=%s, error=%v", req.Id, err)
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &v1.GetUserResponse{
		User: toProtoUserInfo(&user),
	}, nil
}

// GetUserByID 根据用户ID获取用户信息（供 Auth Service 调用，包含敏感信息）
func (s *Service) GetUserByID(ctx context.Context, req *v1.GetUserByIDRequest) (*v1.GetUserByIDResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	var user userv1.User
	if err := s.store.DB().Where("id = ? AND deleted_at IS NULL", req.Id).First(&user).Error; err != nil {
		klog.Warningf("User not found: id=%s, error=%v", req.Id, err)
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &v1.GetUserByIDResponse{
		User:         toProtoUserInfo(&user),
		Salt:         user.Salt,
		PasswordHash: user.PasswordHash,
	}, nil
}

// UpdateUser 更新用户信息
func (s *Service) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	// 1. 检查用户是否存在
	var user userv1.User
	if err := s.store.DB().Where("id = ? AND deleted_at IS NULL", req.Id).First(&user).Error; err != nil {
		klog.Warningf("User not found: id=%s, error=%v", req.Id, err)
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// 2. 构建更新字段
	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	// 如果没有要更新的字段，直接返回当前用户信息
	if len(updates) == 0 {
		return &v1.UpdateUserResponse{
			User: toProtoUserInfo(&user),
		}, nil
	}

	// 3. 更新用户信息
	updates["updated_at"] = time.Now()
	if err := s.store.DB().Model(&user).Updates(updates).Error; err != nil {
		klog.Errorf("Failed to update user: id=%s, error=%v", req.Id, err)
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	// 4. 重新查询更新后的用户信息
	if err := s.store.DB().Where("id = ?", req.Id).First(&user).Error; err != nil {
		klog.Errorf("Failed to fetch updated user: id=%s, error=%v", req.Id, err)
		return nil, status.Error(codes.Internal, "failed to fetch updated user")
	}

	klog.Infof("User updated successfully: id=%s", req.Id)

	return &v1.UpdateUserResponse{
		User: toProtoUserInfo(&user),
	}, nil
}
