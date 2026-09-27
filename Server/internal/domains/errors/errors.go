package errors

import "errors"

// Define Business logic errors
var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrRefreshTokenReuse    = errors.New("refresh token reuse detected")
)
