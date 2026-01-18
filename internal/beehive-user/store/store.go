package store

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"k8s.io/klog/v2"

	"github.com/HappyLadySauce/Beehive/internal/beehive-user/config"
	"github.com/HappyLadySauce/Beehive/internal/pkg/utils/id"
	v1 "github.com/HappyLadySauce/Beehive/pkg/api/user/v1"
)

// Store 数据库存储
type Store struct {
	db *gorm.DB
}

// NewStore 创建新的数据库存储
func NewStore(cfg *config.Config) (*Store, error) {
	// 构建 DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.Port, cfg.Database.SSLMode,
	)

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	store := &Store{
		db: db,
	}

	// 执行数据库迁移
	if err := store.AutoMigrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	klog.Info("Database connection established successfully")
	return store, nil
}

// AutoMigrate 自动迁移数据库表
func (s *Store) AutoMigrate() error {
	// 初始化ID序列
	if err := id.InitSequences(s.db); err != nil {
		return fmt.Errorf("failed to initialize ID sequences: %w", err)
	}

	// 迁移数据库表
	return s.db.AutoMigrate(
		&v1.User{},
	)
}

// Close 关闭数据库连接
func (s *Store) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// DB 返回 gorm.DB 实例
func (s *Store) DB() *gorm.DB {
	return s.db
}
