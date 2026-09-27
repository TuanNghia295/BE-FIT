package entities

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID           uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	TokenHash    string     `gorm:"column:tokenHash;unique;not null"`
	TokenFamily  uuid.UUID  `gorm:"column:tokenFamily;type:uuid;not null"`
	UserID       uuid.UUID  `gorm:"column:userId;type:uuid;not null"`
	RevokeAt     *time.Time `gorm:"column:revokeAt"`
	ExpiresAt    time.Time  `gorm:"column:expiresAt;not null"`
	ReplacedByID *uuid.UUID `gorm:"column:replacedById;type:uuid"`
	CreatedAt    time.Time  `gorm:"column:createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updatedAt"`
}

func (RefreshToken) TableName() string { return "refresh_token" }
