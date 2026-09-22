package entities

import "time"

type RefreshToken struct {
	id          int64
	tokenHash   string
	tokenFamily int64
	userId      int64
	revokeAt    time.Time
	expiredAt   time.Time
	replacedBy  int64
	createdAt   time.Time
	updatedAt   time.Time
}
