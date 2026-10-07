package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/TuanNghia295/BE-FIT/internal/delivery/http/middleware"
	"github.com/TuanNghia295/BE-FIT/internal/domains/entities"
	domainErrors "github.com/TuanNghia295/BE-FIT/internal/domains/errors"
	"github.com/TuanNghia295/BE-FIT/internal/domains/repositories"
	userUsecase "github.com/TuanNghia295/BE-FIT/internal/usecase/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type handlerUserRepo struct{ account *entities.Users }

func (r *handlerUserRepo) Create(context.Context, *entities.Users) error { return nil }
func (r *handlerUserRepo) FindByEmail(context.Context, string) (*entities.Users, error) {
	return r.account, nil
}
func (r *handlerUserRepo) FindByID(_ context.Context, userID uuid.UUID) (*entities.Users, error) {
	if r.account == nil || r.account.ID != userID {
		return nil, domainErrors.ErrUserNotFound
	}
	return r.account, nil
}

type handlerRefreshRepo struct {
	token *entities.RefreshToken
}

func (r *handlerRefreshRepo) Create(_ context.Context, token *entities.RefreshToken) error {
	r.token = token
	return nil
}
func (r *handlerRefreshRepo) FindByTokenHash(_ context.Context, tokenHash string) (*entities.RefreshToken, error) {
	if r.token == nil || r.token.TokenHash != tokenHash {
		return nil, domainErrors.ErrRefreshTokenNotFound
	}
	return r.token, nil
}
func (r *handlerRefreshRepo) Rotate(_ context.Context, oldHash string, replacement *entities.RefreshToken) error {
	if r.token == nil || r.token.TokenHash != oldHash || r.token.RevokeAt != nil {
		return domainErrors.ErrRefreshTokenReuse
	}
	now := time.Now()
	r.token.RevokeAt = &now
	r.token.ReplacedByID = &replacement.ID
	r.token = replacement
	return nil
}
func (r *handlerRefreshRepo) Revoke(context.Context, uuid.UUID) error { return nil }
func (r *handlerRefreshRepo) RevokeFamily(context.Context, uuid.UUID) error {
	if r.token != nil {
		now := time.Now()
		r.token.RevokeAt = &now
	}
	return nil
}

type handlerTokenService struct{}

func (handlerTokenService) GenerateAccessToken(uuid.UUID) (string, error) { return "access-token", nil }
func (handlerTokenService) GenerateRefreshToken() (string, error)         { return "raw-refresh-token", nil }
func (handlerTokenService) HashRefreshToken(token string) string          { return "hashed:" + token }

func newLoginTestHandler(t *testing.T, cookieSecure bool) *UserHandler {
	t.Helper()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}

	userRepo := &handlerUserRepo{account: &entities.Users{
		ID:       uuid.New(),
		Email:    "person@example.com",
		Password: string(passwordHash),
	}}
	refreshRepo := &handlerRefreshRepo{}
	tokenService := handlerTokenService{}
	login := userUsecase.NewLoginUseCase(userRepo, refreshRepo, tokenService)
	refresh := userUsecase.NewRefreshUseCase(refreshRepo, tokenService)
	return NewUserHandler(nil, login, refresh, nil, cookieSecure)
}

func TestLoginWebReceivesHttpOnlyRefreshCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/login", newLoginTestHandler(t, true).Login)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"email":"person@example.com","password":"password123"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Client-Type", "web")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	cookies := response.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != "refresh_token" || !cookies[0].HttpOnly || !cookies[0].Secure {
		t.Fatalf("expected secure HttpOnly refresh cookie, got %#v", cookies)
	}
	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if _, exists := body["refresh_token"]; exists {
		t.Fatal("web response must not expose refresh token in JSON")
	}
}

func TestLoginMobileReceivesRefreshTokenInJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/login", newLoginTestHandler(t, false).Login)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"email":"person@example.com","password":"password123"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Client-Type", "mobile")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["refresh_token"] != "raw-refresh-token" {
		t.Fatalf("expected raw refresh token in mobile response, got %v", body["refresh_token"])
	}
	if cookies := response.Result().Cookies(); len(cookies) != 0 {
		t.Fatalf("mobile response should not set cookies: %#v", cookies)
	}
	if expiry, ok := body["expires_in"].(float64); !ok || expiry != 900 || time.Until(time.Now().Add(time.Duration(expiry)*time.Second)) <= 0 {
		t.Fatalf("unexpected access-token expiry: %v", body["expires_in"])
	}
}

func TestWebRefreshUsesLoginCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	userRepo := &handlerUserRepo{account: &entities.Users{
		ID:       uuid.New(),
		Email:    "web-refresh@example.com",
		Password: string(passwordHash),
	}}
	refreshRepo := &handlerRefreshRepo{}
	tokenService := handlerTokenService{}
	handler := NewUserHandler(
		nil,
		userUsecase.NewLoginUseCase(userRepo, refreshRepo, tokenService),
		userUsecase.NewRefreshUseCase(refreshRepo, tokenService),
		userUsecase.NewMeUseCase(userRepo),
		true,
	)
	router := gin.New()
	router.POST("/login", handler.Login)
	router.POST("/refresh", handler.Refresh)

	loginRequest := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"email":"web-refresh@example.com","password":"password123"}`))
	loginRequest.Header.Set("Content-Type", "application/json")
	loginRequest.Header.Set("X-Client-Type", "web")
	loginResponse := httptest.NewRecorder()
	router.ServeHTTP(loginResponse, loginRequest)
	if loginResponse.Code != http.StatusOK {
		t.Fatalf("expected login status 200, got %d: %s", loginResponse.Code, loginResponse.Body.String())
	}
	cookies := loginResponse.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Path != "/" || !cookies[0].HttpOnly {
		t.Fatalf("expected HttpOnly refresh cookie scoped to /, got %#v", cookies)
	}

	refreshRequest := httptest.NewRequest(http.MethodPost, "/refresh", nil)
	refreshRequest.Header.Set("X-Client-Type", "web")
	refreshRequest.AddCookie(cookies[0])
	refreshResponse := httptest.NewRecorder()
	router.ServeHTTP(refreshResponse, refreshRequest)
	if refreshResponse.Code != http.StatusOK {
		t.Fatalf("expected refresh status 200, got %d: %s", refreshResponse.Code, refreshResponse.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(refreshResponse.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if _, exists := body["refresh_token"]; exists {
		t.Fatal("web refresh response must not expose refresh token in JSON")
	}
}

func TestMeReturnsAuthenticatedUserWithoutPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	account := &entities.Users{
		ID:       userID,
		Email:    "me@example.com",
		FullName: "Me User",
		Password: "must-not-be-returned",
	}
	userRepo := &handlerUserRepo{account: account}
	handler := NewUserHandler(
		nil,
		nil,
		nil,
		userUsecase.NewMeUseCase(userRepo),
		false,
	)
	router := gin.New()
	router.GET("/me", func(c *gin.Context) {
		c.Set(middleware.AuthenticatedUserIDKey, userID)
		c.Next()
	}, handler.Me)

	request := httptest.NewRequest(http.MethodGet, "/me", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["email"] != account.Email || body["fullName"] != account.FullName {
		t.Fatalf("unexpected user response: %#v", body)
	}
	if _, exists := body["password"]; exists {
		t.Fatal("user response must not expose password")
	}
}

var _ repositories.UserRepository = (*handlerUserRepo)(nil)
var _ repositories.RefreshTokenRepository = (*handlerRefreshRepo)(nil)
