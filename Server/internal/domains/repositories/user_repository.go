package repositories

import (
	"context"

	"github.com/TuanNghia295/BE-FIT/internal/domains/entities"
	"github.com/google/uuid"
)

// This is abstraction. interface is the method collection
type UserRepository interface {
	Create(ctx context.Context, user *entities.Users) error
	FindByEmail(ctx context.Context, email string) (*entities.Users, error)
	FindByID(ctx context.Context, userID uuid.UUID) (*entities.Users, error)
}
