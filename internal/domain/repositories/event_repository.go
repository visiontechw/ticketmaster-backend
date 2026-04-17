package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/visiontechw/ticketmaster/internal/domain"
)

type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) (*domain.Event, error)
	Update(ctx context.Context, event *domain.Event) (*domain.Event, error)
	FindByID(ctx context.Context, eventId uuid.UUID) (*domain.Event, error)
	GetEventsByStatusPaginated(ctx context.Context, active bool, page, limit int) ([]*domain.Event, error)
	Delete(ctx context.Context, eventId uuid.UUID, userId uuid.UUID) error
}
