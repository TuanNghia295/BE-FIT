package user

import (
	"context"

	"github.com/TuanNghia295/BE-FIT/internal/domains/entities"
	"github.com/TuanNghia295/BE-FIT/internal/domains/repositories"
	"github.com/google/uuid"
)

type MeUsecase struct {
	userRepo repositories.UserRepository
}

func NewMeUseCase(userRepo repositories.UserRepository) *MeUsecase {
	return &MeUsecase{userRepo: userRepo}
}

func (u *MeUsecase) Execute(ctx context.Context, userID uuid.UUID) (*entities.Users, error) {
	return u.userRepo.FindByID(ctx, userID)
}
