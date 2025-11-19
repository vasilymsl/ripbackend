package ds

import "time"

type User struct {
	ID          int    `gorm:"primaryKey"`
	Username    string `gorm:"uniqueIndex;not null"`
	Password    string `gorm:"not null"`
	Email       string `gorm:"uniqueIndex"`
	FullName    string `gorm:"column:full_name"`
	Phone       string
	IsModerator bool      `gorm:"column:is_moderator;default:false"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// TableName переопределяет имя таблицы для GORM
func (User) TableName() string {
	return "users"
}
