package id

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"k8s.io/klog/v2"
)

const (
	// UserIDSequenceName 用户ID序列名称
	UserIDSequenceName = "user_id_seq"
	// GroupIDSequenceName 群ID序列名称
	GroupIDSequenceName = "group_id_seq"
	// UserIDStartValue 用户ID起始值（10位数字）
	UserIDStartValue = 1000000000
	// GroupIDStartValue 群ID起始值（9位数字）
	GroupIDStartValue = 100000000
)

// Generator ID生成器
type Generator struct {
	db *gorm.DB
}

// NewGenerator 创建新的ID生成器
func NewGenerator(db *gorm.DB) *Generator {
	return &Generator{
		db: db,
	}
}

// InitSequences 初始化数据库序列（如果不存在则创建）
func InitSequences(db *gorm.DB) error {
	// 创建用户ID序列
	userSeqSQL := fmt.Sprintf(
		"CREATE SEQUENCE IF NOT EXISTS %s START WITH %d",
		UserIDSequenceName, UserIDStartValue,
	)
	if err := db.Exec(userSeqSQL).Error; err != nil {
		return fmt.Errorf("failed to create user_id_seq: %w", err)
	}

	// 创建群ID序列
	groupSeqSQL := fmt.Sprintf(
		"CREATE SEQUENCE IF NOT EXISTS %s START WITH %d",
		GroupIDSequenceName, GroupIDStartValue,
	)
	if err := db.Exec(groupSeqSQL).Error; err != nil {
		return fmt.Errorf("failed to create group_id_seq: %w", err)
	}

	klog.Info("ID sequences initialized successfully")
	return nil
}

// GenerateUserID 生成10位用户ID（从1000000000开始递增）
func (g *Generator) GenerateUserID(ctx context.Context) (string, error) {
	var nextID int64
	// 使用 fmt.Sprintf 直接拼接序列名称，因为 PostgreSQL 序列名称不能使用参数绑定
	sql := fmt.Sprintf("SELECT nextval('%s')", UserIDSequenceName)
	if err := g.db.WithContext(ctx).Raw(sql).Scan(&nextID).Error; err != nil {
		return "", fmt.Errorf("failed to generate user ID: %w", err)
	}

	// 验证ID范围（10位数字：1000000000 - 9999999999）
	if nextID < UserIDStartValue || nextID > 9999999999 {
		return "", fmt.Errorf("user ID out of range: %d", nextID)
	}

	userID := fmt.Sprintf("%d", nextID)
	klog.V(4).Infof("Generated user ID: %s", userID)
	return userID, nil
}

// GenerateGroupID 生成9位群ID（从100000000开始递增）
func (g *Generator) GenerateGroupID(ctx context.Context) (string, error) {
	var nextID int64
	// 使用 fmt.Sprintf 直接拼接序列名称，因为 PostgreSQL 序列名称不能使用参数绑定
	sql := fmt.Sprintf("SELECT nextval('%s')", GroupIDSequenceName)
	if err := g.db.WithContext(ctx).Raw(sql).Scan(&nextID).Error; err != nil {
		return "", fmt.Errorf("failed to generate group ID: %w", err)
	}

	// 验证ID范围（9位数字：100000000 - 999999999）
	if nextID < GroupIDStartValue || nextID > 999999999 {
		return "", fmt.Errorf("group ID out of range: %d", nextID)
	}

	groupID := fmt.Sprintf("%d", nextID)
	klog.V(4).Infof("Generated group ID: %s", groupID)
	return groupID, nil
}
