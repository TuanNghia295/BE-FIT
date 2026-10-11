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

type RegisterUsecase struct {
	userRepo repositories.UserRepository
}

func NewRegisterUsecase(userRepo repositories.UserRepository) *RegisterUsecase {
	return &RegisterUsecase{userRepo: userRepo}
}

func (u *RegisterUsecase) Execute(ctx context.Context, email string, fullName, password string) (*entities.Users, error) {
	existingUser, err := u.userRepo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.New("Email already exists")
	}
	if err != nil && !errors.Is(err, domainErrors.ErrUserNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &entities.Users{
		Email:    email,
		FullName: fullName,
		Password: string(hashedPassword),
	}
	if err := u.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
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
	account, err := u.userRepo.FindByEmail(ctx, email)
	if errors.Is(err, domainErrors.ErrUserNotFound) || (err == nil && account == nil) {
		return nil, domainErrors.ErrUnauthorized
	}
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
		return nil, domainErrors.ErrUnauthorized
	}

	accessToken, err := u.tokenGen.GenerateAccessToken(account.ID)
	if err != nil {
		return nil, err
	}

	familyID := uuid.New()
	refreshToken, err := u.tokenGen.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	refreshRecord := &entities.RefreshToken{
		ID:          uuid.New(),
		TokenHash:   u.tokenGen.HashRefreshToken(refreshToken),
		TokenFamily: familyID,
		UserID:      account.ID,
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
		UserID:           account.ID,
	}, nil
}
