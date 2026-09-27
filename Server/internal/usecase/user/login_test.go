package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/TuanNghia295/BE-FIT/internal/domains/entities"
	domainErrors "github.com/TuanNghia295/BE-FIT/internal/domains/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type testUserRepository struct {
	user *entities.Users
	err  error
}

func (r *testUserRepository) Create(_ context.Context, user *entities.Users) error {
	r.user = user
	return nil
}

func (r *testUserRepository) FindByEmail(context.Context, string) (*entities.Users, error) {
	return r.user, r.err
}

type testRefreshTokenRepository struct {
	token         *entities.RefreshToken
	created       *entities.RefreshToken
	familyRevoked bool
	rotated       bool
	rotationErr   error
}

func (r *testRefreshTokenRepository) Create(_ context.Context, token *entities.RefreshToken) error {
	r.created = token
	return nil
}

func (r *testRefreshTokenRepository) FindByTokenHash(context.Context, string) (*entities.RefreshToken, error) {
	if r.token == nil {
		return nil, domainErrors.ErrRefreshTokenNotFound
	}
	return r.token, nil
}

func (r *testRefreshTokenRepository) Rotate(_ context.Context, _ string, replacement *entities.RefreshToken) error {
	if r.rotationErr != nil {
		return r.rotationErr
	}
	r.rotated = true
	r.created = replacement
	return nil
}

func (r *testRefreshTokenRepository) Revoke(context.Context, uuid.UUID) error {
	return nil
}

func (r *testRefreshTokenRepository) RevokeFamily(context.Context, uuid.UUID) error {
	r.familyRevoked = true
	return nil
}

type testTokenService struct{}

func (testTokenService) GenerateAccessToken(uuid.UUID) (string, error) {
	return "access-token", nil
}

func (testTokenService) GenerateRefreshToken() (string, error) {
	return "refresh-token", nil
}

func (testTokenService) HashRefreshToken(token string) string {
	return "hash:" + token
}

func TestLoginUsecaseCreatesHashedRefreshToken(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}

	account := &entities.Users{ID: uuid.New(), Email: "person@example.com", Password: string(passwordHash)}
	users := &testUserRepository{user: account}
	refreshTokens := &testRefreshTokenRepository{}
	usecase := NewLoginUseCase(users, refreshTokens, testTokenService{})

	result, err := usecase.Execute(context.Background(), account.Email, "password123")
	if err != nil {
		t.Fatalf("Execute returned an error: %v", err)
	}
	if result.AccessToken != "access-token" || result.RefreshToken != "refresh-token" {
		t.Fatalf("unexpected token result: %#v", result)
	}
	if refreshTokens.created == nil || refreshTokens.created.TokenHash != "hash:refresh-token" {
		t.Fatalf("expected only the hashed refresh token to be persisted")
	}
	if refreshTokens.created.UserID != account.ID || refreshTokens.created.TokenFamily == uuid.Nil {
		t.Fatalf("refresh token is missing user or family association")
	}
	if !refreshTokens.created.ExpiresAt.After(time.Now()) {
		t.Fatalf("refresh token expiry should be in the future")
	}
}

func TestLoginUsecaseRejectsInvalidPassword(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}

	account := &entities.Users{ID: uuid.New(), Email: "person@example.com", Password: string(passwordHash)}
	refreshTokens := &testRefreshTokenRepository{}
	usecase := NewLoginUseCase(&testUserRepository{user: account}, refreshTokens, testTokenService{})

	_, err = usecase.Execute(context.Background(), account.Email, "wrong-password")
	if !errors.Is(err, domainErrors.ErrUnauthorized) {
		t.Fatalf("expected unauthorized error, got %v", err)
	}
	if refreshTokens.created != nil {
		t.Fatal("invalid credentials must not create a refresh token")
	}
}

func TestRefreshUsecaseRotatesTokenWithinFamily(t *testing.T) {
	familyID := uuid.New()
	accountID := uuid.New()
	current := &entities.RefreshToken{
		ID:          uuid.New(),
		TokenHash:   "hash:old-refresh-token",
		TokenFamily: familyID,
		UserID:      accountID,
		ExpiresAt:   time.Now().Add(time.Hour),
	}
	refreshTokens := &testRefreshTokenRepository{token: current}
	usecase := NewRefreshUseCase(refreshTokens, testTokenService{})

	result, err := usecase.Execute(context.Background(), "old-refresh-token")
	if err != nil {
		t.Fatalf("Execute returned an error: %v", err)
	}
	if !refreshTokens.rotated || refreshTokens.created.TokenFamily != familyID {
		t.Fatal("expected replacement token to rotate within the same family")
	}
	if result.RefreshToken != "refresh-token" || result.UserID != accountID {
		t.Fatalf("unexpected refresh result: %#v", result)
	}
}

func TestRefreshUsecaseRevokesFamilyOnReuse(t *testing.T) {
	now := time.Now()
	current := &entities.RefreshToken{
		ID:          uuid.New(),
		TokenHash:   "hash:old-refresh-token",
		TokenFamily: uuid.New(),
		UserID:      uuid.New(),
		RevokeAt:    &now,
		ExpiresAt:   time.Now().Add(time.Hour),
	}
	refreshTokens := &testRefreshTokenRepository{token: current}
	usecase := NewRefreshUseCase(refreshTokens, testTokenService{})

	_, err := usecase.Execute(context.Background(), "old-refresh-token")
	if !errors.Is(err, domainErrors.ErrUnauthorized) {
		t.Fatalf("expected unauthorized error, got %v", err)
	}
	if !refreshTokens.familyRevoked {
		t.Fatal("reused refresh token should revoke its token family")
	}
}
