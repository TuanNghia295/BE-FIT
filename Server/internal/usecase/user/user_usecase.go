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

type UserProfile struct {
	User    *entities.Users
	Fitness *entities.FitnessProfile
}

func NewMeUseCase(userRepo repositories.UserRepository) *MeUsecase {
	return &MeUsecase{userRepo: userRepo}
}

func (u *MeUsecase) Execute(ctx context.Context, userID uuid.UUID) (*entities.Users, error) {
	return u.userRepo.FindByID(ctx, userID)
}

// UpdateProfile updates the user's profile information in the database.
// updatedUser is a pointer to an entities.Users and entities.FitnessProfile struct that contains the updated profile information.
func (u *MeUsecase) UpdateProfile(ctx context.Context, userID uuid.UUID, updatedUser UserProfile) (*entities.Users, error) {
	return u.userRepo.Update(ctx, userID, updatedUser.User, updatedUser.Fitness)
}
