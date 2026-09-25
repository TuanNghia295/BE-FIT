package http

// use to handle and process HTTP request

import (
	"net/http"

	"github.com/TuanNghia295/BE-FIT/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	registerUsecase *user.RegisterUsecase
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
}

func NewUserHandler(registerUscase *user.RegisterUsecase) *UserHandler {
	return &UserHandler{
		registerUsecase: registerUscase,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	user, err := h.registerUsecase.Execute(req.Email, req.FullName, req.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}
