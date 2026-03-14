package svc

import (
	"time"

	"github.com/HappyLadySauce/Beehive/services/conversation/internal/config"
	"github.com/HappyLadySauce/Beehive/services/conversation/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Conv   *model.ConversationModel
	JoinReq *model.GroupJoinRequestModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	if c.PostgresDSN == "" {
		panic("conversation service requires PostgresDSN")
	}
	db, err := gorm.Open(postgres.Open(c.PostgresDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetConnMaxLifetime(time.Hour)
	}
	return &ServiceContext{
		Config:  c,
		DB:      db,
		Conv:    model.NewConversationModel(db),
		JoinReq: model.NewGroupJoinRequestModel(db),
	}
}
