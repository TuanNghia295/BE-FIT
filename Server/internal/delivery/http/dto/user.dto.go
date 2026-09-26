package dto

type RegisterUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"fullName" binding:"required,min=2,max=100"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}
