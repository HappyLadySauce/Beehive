package svc

import (
	"time"

	"github.com/HappyLadySauce/Beehive/services/message/internal/config"
	"github.com/HappyLadySauce/Beehive/services/message/internal/model"
	"github.com/HappyLadySauce/Beehive/services/message/internal/mq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ServiceContext struct {
	Config  config.Config
	DB      *gorm.DB
	Msg     *model.MessageModel
	MQ      *mq.Publisher
}

func NewServiceContext(c config.Config) *ServiceContext {
	if c.PostgresDSN == "" {
		panic("message service requires PostgresDSN")
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
	var pub *mq.Publisher
	if c.RabbitMQURL != "" {
		pub, err = mq.NewPublisher(c.RabbitMQURL, c.RabbitMQExchange, c.RabbitMQRouteKey)
		if err != nil {
			panic("rabbitmq: " + err.Error())
		}
	}
	return &ServiceContext{
		Config: c,
		DB:     db,
		Msg:    model.NewMessageModel(db),
		MQ:     pub,
	}
}
