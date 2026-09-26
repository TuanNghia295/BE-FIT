package http

// use to handle and process HTTP request

import (
	"net/http"

	"github.com/TuanNghia295/BE-FIT/internal/delivery/http/dto"
	"github.com/TuanNghia295/BE-FIT/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	registerUsecase *user.RegisterUsecase
}

func NewUserHandler(registerUscase *user.RegisterUsecase) *UserHandler {
	return &UserHandler{
		registerUsecase: registerUscase,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterUserRequest

	// ShouldBindJSON will parse JSON from FrontEnd -> Map JSON to DTO -> Validate binding tags -> Return errors if validate fail
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	user, err := h.registerUsecase.Execute(req.Email, req.FullName, req.Password)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, user)
}
