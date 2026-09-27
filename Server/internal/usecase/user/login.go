package user

import (
	"context"
	"errors"
	"time"

	"github.com/TuanNghia295/BE-FIT/internal/domains/entities"
	domainErrors "github.com/TuanNghia295/BE-FIT/internal/domains/errors"
	"github.com/TuanNghia295/BE-FIT/internal/domains/repositories"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginTokenService interface {
	GenerateAccessToken(userID uuid.UUID) (string, error)
	GenerateRefreshToken() (string, error)
	HashRefreshToken(token string) string
}

type LoginUsecase struct {
	userRepo    repositories.UserRepository
	refreshRepo repositories.RefreshTokenRepository
	tokenGen    LoginTokenService
}

type LoginResult struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	TokenType        string    `json:"token_type"`
	ExpiresIn        int       `json:"expires_in"`
	RefreshExpiresAt time.Time `json:"-"`
	UserID           uuid.UUID `json:"user_id"`
}

func NewLoginUseCase(
	userRepo repositories.UserRepository,
	refreshRepo repositories.RefreshTokenRepository,
	tokenGen LoginTokenService,
) *LoginUsecase {
	return &LoginUsecase{
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
		tokenGen:    tokenGen,
	}
}

func (u *LoginUsecase) Execute(ctx context.Context, email string, password string) (*LoginResult, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if errors.Is(err, domainErrors.ErrUserNotFound) || (err == nil && user == nil) {
		return nil, domainErrors.ErrUnauthorized
	}
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, domainErrors.ErrUnauthorized
	}

	// create access token
	accessToken, err := u.tokenGen.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	// create family token
	familyID := uuid.New()

	// create refresh token
	refreshToken, err := u.tokenGen.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	refreshRecord := &entities.RefreshToken{
		ID:          uuid.New(),
		TokenHash:   u.tokenGen.HashRefreshToken(refreshToken),
		TokenFamily: familyID,
		UserID:      user.ID,
		ExpiresAt:   expiresAt,
	}
	if err := u.refreshRepo.Create(ctx, refreshRecord); err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        15 * 60,
		RefreshExpiresAt: expiresAt,
		UserID:           user.ID,
	}, nil
}
