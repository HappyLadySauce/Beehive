package svc

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive/services/user/internal/config"
	"github.com/HappyLadySauce/Beehive/services/user/internal/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config

	DB             *gorm.DB
	Redis          *redis.Client
	UserMod            *model.UserModel
	UserProfileMod     *model.UserProfileModel
	ContactMod         *model.ContactModel
	ContactRequestMod  *model.ContactRequestModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(postgres.Open(c.PostgresDSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     c.RedisAddr,
		Password: c.RedisPassword,
		DB:       c.RedisDB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	// 简单连接池设置，可按需从配置扩展。
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	return &ServiceContext{
		Config:         c,
		DB:             db,
		Redis:          rdb,
		UserMod:           model.NewUserModel(db),
		UserProfileMod:    model.NewUserProfileModel(db),
		ContactMod:        model.NewContactModel(db),
		ContactRequestMod: model.NewContactRequestModel(db),
	}
}
