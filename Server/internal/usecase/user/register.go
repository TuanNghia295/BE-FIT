package user

// Usecase use to write business logic
import (
	"context"
	"errors"

	"github.com/TuanNghia295/BE-FIT/internal/domains/entities"
	domainErrors "github.com/TuanNghia295/BE-FIT/internal/domains/errors"
	"github.com/TuanNghia295/BE-FIT/internal/domains/repositories"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUsecase struct {
	userRepo repositories.UserRepository
}

func NewRegisterUsecase(userRepo repositories.UserRepository) *RegisterUsecase {
	return &RegisterUsecase{
		userRepo: userRepo,
	}
}

func (u *RegisterUsecase) Execute(ctx context.Context, email string, fullName, password string) (*entities.Users, error) {
	// check email exist
	existingEmail, err := u.userRepo.FindByEmail(ctx, email)
	if err == nil && existingEmail != nil {
		return nil, errors.New("Email already exists")
	}
	if err != nil && !errors.Is(err, domainErrors.ErrUserNotFound) {
		return nil, err
	}

	// Hash password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := &entities.Users{
		Email:    email,
		FullName: fullName,
		Password: string(hashedPass),
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
