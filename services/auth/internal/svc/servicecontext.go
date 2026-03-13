package svc

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive/services/auth/internal/config"
	"github.com/HappyLadySauce/Beehive/services/auth/internal/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config

	DB      *gorm.DB
	Redis   *redis.Client
	RBACMod *model.RBACModel
	UserMod *model.UserModel
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
