package service

import (
	v1 "github.com/HappyLadySauce/Beehive/pkg/api/proto/user/v1"
	userv1 "github.com/HappyLadySauce/Beehive/pkg/api/user/v1"
)

// toProtoUserInfo 将 User 模型转换为 proto UserInfo
func toProtoUserInfo(user *userv1.User) *v1.UserInfo {
	return &v1.UserInfo{
		Id:          user.ID,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Email:       user.Email,
		Description: user.Description,
		Level:       int32(user.Level),
		Status:      user.Status,
	}
}
