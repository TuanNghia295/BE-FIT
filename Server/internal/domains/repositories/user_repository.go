package repositories

import (
	"github.com/TuanNghia295/BE-FIT/internal/domains/entities"
)

// This is abstraction. interface is the method collection
type UserRepository interface {
	Create(user *entities.Users) error
	FindByEmail(email string) (*entities.Users, error)
}
