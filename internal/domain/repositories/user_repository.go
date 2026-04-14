package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/visiontechw/ticketmaster/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
