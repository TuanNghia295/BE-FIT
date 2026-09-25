package postgres

// Used to write functions related to database queries
import (
	"github.com/TuanNghia295/BE-FIT/internal/domains/entities"
	"github.com/TuanNghia295/BE-FIT/internal/domains/repositories"
	"github.com/jinzhu/gorm"
)

type userRepository struct {
	db *gorm.DB // Database handle. Ready to read/write database
	// db *gorm.DB: GORM’s database object is a handle used to execute queries
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entities.Users) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*entities.Users, error) {
	var user entities.Users

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
