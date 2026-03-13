package model

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Role 对应系统级角色表 roles。
type Role struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey"`
	Name        string    `gorm:"column:name;type:text;uniqueIndex;not null"` // 例如：user / admin / super_admin
	Description string    `gorm:"column:description;type:text;not null;default:''"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamptz;not null;default:now()"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`
}

func (Role) TableName() string {
	return "roles"
}

// Permission 对应权限表 permissions。
type Permission struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey"`
	Code        string    `gorm:"column:code;type:text;uniqueIndex;not null"` // 例如：admin.user.read
	Description string    `gorm:"column:description;type:text;not null;default:''"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamptz;not null;default:now()"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`
}

func (Permission) TableName() string {
	return "permissions"
}

// RolePermission 连接角色与权限，多对多关系。
type RolePermission struct {
	RoleID       string    `gorm:"column:role_id;type:uuid;primaryKey"`
	PermissionID string    `gorm:"column:permission_id;type:uuid;primaryKey"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamptz;not null;default:now()"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

// UserRole 连接用户与角色，多对多关系。
type UserRole struct {
	UserID    string    `gorm:"column:user_id;type:uuid;primaryKey"`
	RoleID    string    `gorm:"column:role_id;type:uuid;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null;default:now()"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

// RBACModel 提供 RBAC 相关的基础查询。
type RBACModel struct {
	db *gorm.DB
}

func NewRBACModel(db *gorm.DB) *RBACModel {
	return &RBACModel{db: db}
}

// GetUserRoles 查询用户的所有系统级角色名称。
func (m *RBACModel) GetUserRoles(userID string) ([]string, error) {
	var roleNames []string
	err := m.db.
		Table("user_roles ur").
		Joins("JOIN roles r ON ur.role_id = r.id").
		Where("ur.user_id = ?", userID).
		Pluck("r.name", &roleNames).Error
	if err != nil {
		return nil, err
	}
	return roleNames, nil
}

// GetUserPermissions 查询用户聚合后的权限编码集合。
func (m *RBACModel) GetUserPermissions(userID string) ([]string, error) {
	var perms []string
	err := m.db.
		Table("user_roles ur").
		Joins("JOIN role_permissions rp ON ur.role_id = rp.role_id").
		Joins("JOIN permissions p ON rp.permission_id = p.id").
		Where("ur.user_id = ?", userID).
		Pluck("DISTINCT p.code", &perms).Error
	if err != nil {
		return nil, err
	}
	return perms, nil
}

// ReplaceUserRoles 以覆盖方式重置某个用户的角色列表：
// 先删除该用户所有 user_roles，再按给定角色名写入。
func (m *RBACModel) ReplaceUserRoles(ctx context.Context, userID string, roleNames []string) error {
	tx := m.db.WithContext(ctx).Begin()
	if err := tx.Where("user_id = ?", userID).Delete(&UserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if len(roleNames) == 0 {
		return tx.Commit().Error
	}

	var roles []Role
	if err := tx.Where("name IN ?", roleNames).Find(&roles).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 确保所有传入的角色名称都存在，否则回滚并返回错误，避免静默丢失角色。
	found := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		found[r.Name] = struct{}{}
	}
	var missing []string
	for _, name := range roleNames {
		if _, ok := found[name]; !ok {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		tx.Rollback()
		return fmt.Errorf("roles not found: %v", missing)
	}

	now := time.Now()
	for _, r := range roles {
		ur := UserRole{
			UserID:    userID,
			RoleID:    r.ID,
			CreatedAt: now,
		}
		if err := tx.Create(&ur).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

