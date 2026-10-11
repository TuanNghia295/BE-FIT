package postgres

// Used to write functions related to database queries
import (
	"context"
	"errors"

	"github.com/TuanNghia295/BE-FIT/internal/domains/entities"
	domainErrors "github.com/TuanNghia295/BE-FIT/internal/domains/errors"
	"github.com/TuanNghia295/BE-FIT/internal/domains/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	db *gorm.DB // Database handle. Ready to read/write database
	// db *gorm.DB: GORM’s database object is a handle used to execute queries
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.Users) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entities.Users, error) {
	var user entities.Users

	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, userID uuid.UUID) (*entities.Users, error) {
	var user entities.Users

	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(
	ctx context.Context,
	userID uuid.UUID,
	updatedUser *entities.Users,
	updatedFitness *entities.FitnessProfile,
) (*entities.Users, error) {
	var user entities.Users

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&entities.Users{}).
			Where("id = ?", userID).
			Updates(map[string]interface{}{
				"email":    updatedUser.Email,
				"fullName": updatedUser.FullName,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return domainErrors.ErrUserNotFound
		}

		updatedFitness.UserID = userID
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "userId"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"gender",
				"height",
				"weight",
				"level",
				"updatedAt",
			}),
		}).Create(updatedFitness).Error; err != nil {
			return err
		}

		return tx.Where("id = ?", userID).First(&user).Error
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}
