package http

// use to handle and process HTTP request

import (
	"net/http"
	"strings"
	"time"

	"github.com/TuanNghia295/BE-FIT/internal/delivery/http/dto"
	domainErrors "github.com/TuanNghia295/BE-FIT/internal/domains/errors"
	"github.com/TuanNghia295/BE-FIT/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	registerUsecase *user.RegisterUsecase
	loginUsecase    *user.LoginUsecase
	refreshUsecase  *user.RefreshUsecase
	cookieSecure    bool
}

func NewUserHandler(
	registerUsecase *user.RegisterUsecase,
	loginUsecase *user.LoginUsecase,
	refreshUsecase *user.RefreshUsecase,
	cookieSecure bool,
) *UserHandler {
	return &UserHandler{
		registerUsecase: registerUsecase,
		loginUsecase:    loginUsecase,
		refreshUsecase:  refreshUsecase,
		cookieSecure:    cookieSecure,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterUserRequest

	// ShouldBindJSON will parse JSON from FrontEnd -> Map JSON to DTO -> Validate binding tags -> Return errors if validate fail
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	user, err := h.registerUsecase.Execute(c.Request.Context(), req.Email, req.FullName, req.Password)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	clientType, ok := authClientType(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Client-Type must be web or mobile"})
		return
	}

	var req dto.LoginUserRequest

	// ShouldBindJSON will parse JSON from FrontEnd -> Map JSON to DTO -> Validate binding tags -> Return errors if validate fail
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid login request"})
		return
	}

	result, err := h.loginUsecase.Execute(c.Request.Context(), req.Email, req.Password)

	if err != nil {
		c.Error(err)
		return
	}

	h.writeAuthResponse(c, clientType, result)
}

func (h *UserHandler) Refresh(c *gin.Context) {
	clientType, ok := authClientType(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Client-Type must be web or mobile"})
		return
	}

	var rawToken string
	if clientType == "web" {
		var err error
		rawToken, err = c.Cookie("refresh_token")
		if err != nil {
			c.Error(domainErrors.ErrUnauthorized)
			return
		}
	} else {
		var req dto.RefreshTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
			return
		}
		rawToken = req.RefreshToken
	}

	result, err := h.refreshUsecase.Execute(c.Request.Context(), rawToken)
	if err != nil {
		c.Error(err)
		return
	}

	h.writeAuthResponse(c, clientType, result)
}

func authClientType(c *gin.Context) (string, bool) {
	clientType := strings.ToLower(strings.TrimSpace(c.GetHeader("X-Client-Type")))
	return clientType, clientType == "web" || clientType == "mobile"
}

func (h *UserHandler) writeAuthResponse(c *gin.Context, clientType string, result *user.LoginResult) {
	if clientType == "web" {
		c.SetSameSite(http.SameSiteLaxMode)
		maxAge := int(time.Until(result.RefreshExpiresAt).Seconds())
		if maxAge < 0 {
			maxAge = 0
		}
		c.SetCookie("refresh_token", result.RefreshToken, maxAge, "/refresh", "", h.cookieSecure, true)
		c.JSON(http.StatusOK, gin.H{
			"access_token": result.AccessToken,
			"token_type":   result.TokenType,
			"expires_in":   result.ExpiresIn,
			"user_id":      result.UserID,
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
