package usecase

import (
	"github.com/TuanNghia295/BE-FIT/internal/domains/repositories"
	"github.com/TuanNghia295/BE-FIT/internal/usecase/user"
)

type UseCase struct {
	RegisterUser *user.RegisterUsecase
	LoginUser    *user.LoginUsecase
	RefreshUser  *user.RefreshUsecase
	MeUser       *user.MeUsecase
}

func NewUseCase(
	userRepo repositories.UserRepository,
	refreshRepo repositories.RefreshTokenRepository,
	tokenService user.LoginTokenService,
) *UseCase {
	return &UseCase{
		RegisterUser: user.NewRegisterUsecase(userRepo),
		LoginUser:    user.NewLoginUseCase(userRepo, refreshRepo, tokenService),
		RefreshUser:  user.NewRefreshUseCase(refreshRepo, tokenService),
		MeUser:       user.NewMeUseCase(userRepo),
	}
}
