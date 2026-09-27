package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/TuanNghia295/BE-FIT/internal/domains/entities"
	domainErrors "github.com/TuanNghia295/BE-FIT/internal/domains/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: db,
	}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *entities.RefreshToken) error {
	// using context to manage deadline, goroutine, cancellation of this request
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *RefreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*entities.RefreshToken, error) {
	var token entities.RefreshToken

	err := r.db.WithContext(ctx).Where(`"tokenHash" = ?`, tokenHash).First(&token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainErrors.ErrRefreshTokenNotFound
	}
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *RefreshTokenRepository) Rotate(ctx context.Context, tokenHash string, replacement *entities.RefreshToken) error {
	var rotationErr error
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var current entities.RefreshToken
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(`"tokenHash" = ?`, tokenHash).
			First(&current).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domainErrors.ErrRefreshTokenNotFound
		}
		if err != nil {
			return err
		}

		now := time.Now()
		if current.RevokeAt != nil {
			rotationErr = domainErrors.ErrRefreshTokenReuse
			return tx.Model(&entities.RefreshToken{}).
				Where(`"tokenFamily" = ?`, current.TokenFamily).
				Update("RevokeAt", now).Error
		}
		if !current.ExpiresAt.After(now) {
			rotationErr = domainErrors.ErrRefreshTokenExpired
			return tx.Model(&current).Update("RevokeAt", now).Error
		}

		replacement.ID = uuid.New()
		replacement.UserID = current.UserID
		replacement.TokenFamily = current.TokenFamily
		current.RevokeAt = &now
		current.ReplacedByID = &replacement.ID

		if err := tx.Save(&current).Error; err != nil {
			return err
		}
		return tx.Create(replacement).Error
	})
	if err != nil {
		return err
	}
	return rotationErr
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, tokenID uuid.UUID) error {
	now := time.Now()

	return r.db.WithContext(ctx).Model(&entities.RefreshToken{}).
		Where("id = ?", tokenID).
		Update("RevokeAt", now).Error
}

func (r *RefreshTokenRepository) RevokeFamily(ctx context.Context, tokenFamily uuid.UUID) error {
	now := time.Now()

	return r.db.WithContext(ctx).Model(&entities.RefreshToken{}).
		Where(`"tokenFamily" = ?`, tokenFamily).
		Update("RevokeAt", now).Error
}
