package svc

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive/services/auth/internal/config"
	"github.com/HappyLadySauce/Beehive/services/auth/internal/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ServiceContext struct {
	Config config.Config

	DB      *gorm.DB
	Redis   *redis.Client
	RBACMod *model.RBACModel
	UserMod *model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 使用 Silent 避免登录失败（用户不存在）时 GORM 打印 "record not found" 到控制台
	db, err := gorm.Open(postgres.Open(c.PostgresDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
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

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	return &ServiceContext{
		Config:  c,
		DB:      db,
		Redis:   rdb,
		RBACMod: model.NewRBACModel(db),
		UserMod: model.NewUserModel(db),
	}
}
