package dto

type RegisterUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"fullName" binding:"required,min=2,max=100"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogOutUserRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
