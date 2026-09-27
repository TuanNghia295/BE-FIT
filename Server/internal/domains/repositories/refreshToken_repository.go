package repositories

import (
	"context"

	"github.com/TuanNghia295/BE-FIT/internal/domains/entities"
	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entities.RefreshToken) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*entities.RefreshToken, error)
	Rotate(ctx context.Context, tokenHash string, replacement *entities.RefreshToken) error
	Revoke(ctx context.Context, tokenID uuid.UUID) error
	RevokeFamily(ctx context.Context, tokenFamily uuid.UUID) error
}
