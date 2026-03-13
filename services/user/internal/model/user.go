package model

import (
	"time"

	"gorm.io/gorm"
)

// User 对应 users 表，主要给 AuthService 等使用，这里只建最小字段。
type User struct {
	ID           string    `gorm:"column:id;type:uuid;primaryKey"`
	Username     string    `gorm:"column:username;type:text;uniqueIndex;not null"`
	PasswordHash string    `gorm:"column:password_hash;type:text;not null"`
	Status       string    `gorm:"column:status;type:text;not null;default:'normal'"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamptz;not null;default:now()"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`
}

func (User) TableName() string {
	return "users"
}

// UserProfile 对应 user_profiles 表，给 UserService 读写 profile 使用。
type UserProfile struct {
	UserID    string    `gorm:"column:user_id;type:uuid;primaryKey"`
	Nickname  string    `gorm:"column:nickname;type:text;not null;default:''"`
	AvatarURL string    `gorm:"column:avatar_url;type:text;not null;default:''"`
	Bio       string    `gorm:"column:bio;type:text;not null;default:''"`
	Status    string    `gorm:"column:status;type:text;not null;default:'normal'"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}

type UserProfileModel struct {
	db *gorm.DB
}

func NewUserProfileModel(db *gorm.DB) *UserProfileModel {
	return &UserProfileModel{db: db}
}

func (m *UserProfileModel) FindByID(id string) (*UserProfile, error) {
	var p UserProfile
	if err := m.db.Where("user_id = ?", id).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (m *UserProfileModel) FindByIDs(ids []string) ([]*UserProfile, error) {
	if len(ids) == 0 {
		return []*UserProfile{}, nil
	}
	var profiles []*UserProfile
	if err := m.db.Where("user_id IN ?", ids).Find(&profiles).Error; err != nil {
		return nil, err
	}
	return profiles, nil
}

func (m *UserProfileModel) UpdateProfile(p *UserProfile) error {
	return m.db.Model(&UserProfile{}).
		Where("user_id = ?", p.UserID).
		Updates(map[string]interface{}{
			"nickname":   p.Nickname,
			"avatar_url": p.AvatarURL,
			"bio":        p.Bio,
			"status":     p.Status,
			"updated_at": time.Now(),
		}).Error
}

