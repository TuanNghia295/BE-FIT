package user

// Usecase use to write business logic
import (
	"errors"

	"github.com/TuanNghia295/BE-FIT/internal/domains/entities"
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

func (u *RegisterUsecase) Execute(email string, fullName, password string) (*entities.Users, error) {
	// check email exist
	existingEmail, _ := u.userRepo.FindByEmail(email)

	if existingEmail != nil {
		return nil, errors.New("Email already exists")
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

	err = u.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
