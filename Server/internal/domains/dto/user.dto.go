package dto

type RegisterReq struct {
	Email                string `json:"email" binding:"required,email"`
	FullName             string `json:"fullName" binding:"required"`
	Password             string `json:"password" binding:"required,min=6"`
	ConfirmationPassword string `json:"confirmationPassword" binding:"required"`
}

type UserRes struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
}
