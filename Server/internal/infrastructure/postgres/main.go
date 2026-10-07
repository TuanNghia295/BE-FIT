package postgres

import (
	"github.com/TuanNghia295/BE-FIT/internal/domains/repositories"
	"gorm.io/gorm"
)

type ConfigPostgreRepo struct {
	UserRepository         repositories.UserRepository
	RefreshTokenRepository repositories.RefreshTokenRepository
}

func ConfigRepo(db *gorm.DB) ConfigPostgreRepo {
	return ConfigPostgreRepo{
		UserRepository:         NewUserRepository(db),
		RefreshTokenRepository: NewRefreshTokenRepository(db),
	}
}
