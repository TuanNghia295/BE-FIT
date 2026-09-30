package user

import (
	"context"
	"errors"
	"time"

	"github.com/TuanNghia295/BE-FIT/internal/domains/entities"
	domainErrors "github.com/TuanNghia295/BE-FIT/internal/domains/errors"
	"github.com/TuanNghia295/BE-FIT/internal/domains/repositories"
	"github.com/google/uuid"
)

type RefreshUsecase struct {
	refreshRepo repositories.RefreshTokenRepository
	tokenGen    LoginTokenService
}

func NewRefreshUseCase(refreshRepo repositories.RefreshTokenRepository, tokenGen LoginTokenService) *RefreshUsecase {
	return &RefreshUsecase{refreshRepo: refreshRepo, tokenGen: tokenGen}
}

func (u *RefreshUsecase) Execute(ctx context.Context, rawToken string) (*LoginResult, error) {
	if rawToken == "" {
		return nil, domainErrors.ErrUnauthorized
	}

	oldHash := u.tokenGen.HashRefreshToken(rawToken)
	current, err := u.refreshRepo.FindByTokenHash(ctx, oldHash)
	if err != nil {
		if errors.Is(err, domainErrors.ErrRefreshTokenNotFound) {
			return nil, domainErrors.ErrUnauthorized
		}
		return nil, err
	}

	if current.RevokeAt != nil {
		if err := u.refreshRepo.RevokeFamily(ctx, current.TokenFamily); err != nil {
			return nil, err
		}
		return nil, domainErrors.ErrUnauthorized
	}
	if !current.ExpiresAt.After(time.Now()) {
		if err := u.refreshRepo.Revoke(ctx, current.ID); err != nil {
			return nil, err
		}
		return nil, domainErrors.ErrUnauthorized
	}

	accessToken, err := u.tokenGen.GenerateAccessToken(current.UserID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.tokenGen.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	replacement := &entities.RefreshToken{
		ID:          uuid.New(),
		TokenHash:   u.tokenGen.HashRefreshToken(refreshToken),
		TokenFamily: current.TokenFamily,
		UserID:      current.UserID,
		ExpiresAt:   expiresAt,
	}
	if err := u.refreshRepo.Rotate(ctx, oldHash, replacement); err != nil {
		if errors.Is(err, domainErrors.ErrRefreshTokenNotFound) ||
			errors.Is(err, domainErrors.ErrRefreshTokenExpired) ||
			errors.Is(err, domainErrors.ErrRefreshTokenReuse) {
			return nil, domainErrors.ErrUnauthorized
		}
		return nil, err
	}

	return &LoginResult{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        15 * 60,
		RefreshExpiresAt: expiresAt,
		UserID:           current.UserID,
	}, nil
}

func (u *RefreshUsecase) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return domainErrors.ErrUnauthorized
	}
	// Hash token
	tokenHash := u.tokenGen.HashRefreshToken(refreshToken)
	// Find token hash from DB to validate
	token, err := u.refreshRepo.FindByTokenHash(ctx, tokenHash)
	if errors.Is(err, domainErrors.ErrRefreshTokenNotFound) {
		return domainErrors.ErrUnauthorized
	}
	if err != nil {
		return err
	}

	return u.refreshRepo.RevokeFamily(ctx, token.TokenFamily)
}
