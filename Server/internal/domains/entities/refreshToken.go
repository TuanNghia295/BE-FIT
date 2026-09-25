package entities

import "time"

type RefreshToken struct {
	ID          string    `gorm:"column:id"`
	TokenHash   string    `gorm:"column:tokenHash"`
	TokenFamily int64     `gorm:"column:tokenFamily"`
	UserID      string    `gorm:"column:userId"`
	RevokeAt    time.Time `gorm:"column:revokeAt"`
	ExpiredAt   time.Time `gorm:"column:expiredAt"`
	ReplacedBy  string    `gorm:"column:replacedBy"`
	CreatedAt   time.Time `gorm:"column:createdAt"`
	UpdatedAt   time.Time `gorm:"column:updatedAt"`
}

func (RefreshToken) TableName() string {
	return "RefreshToken"
}
