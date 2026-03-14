package svc

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/presence/internal/config"
	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.RedisAddr,
		Password: c.RedisPassword,
		DB:       c.RedisDB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config: c,
		Redis:  rdb,
	}
}
