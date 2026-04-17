package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/visiontechw/ticketmaster/internal/domain"
)

type ListEventsInput struct {
	Active bool
	Page   int
	Limit  int
}

type CreateEventRequest struct {
	Name        string    `json:"name" validate:"required,min=3,max=255"`
	Description string    `json:"description" validate:"max=1000"`
	OccursAt    time.Time `json:"occurs_at" validate:"required,gt=now"`
	Location    string    `json:"location" validate:"required"`
	Capacity    int64     `json:"capacity" validate:"required,min=1"`
	EventTypeId uuid.UUID `json:"event_type_id" validate:"required"`
}

type UpdateEventRequest struct {
	Name        string    `json:"name" validate:"required,min=3,max=255"`
	Description string    `json:"description" validate:"max=1000"`
	OccursAt    time.Time `json:"occurs_at" validate:"required,gt=now"`
	Location    string    `json:"location" validate:"required"`
	Capacity    int64     `json:"capacity" validate:"required,min=1"`
	EventTypeId uuid.UUID `json:"event_type_id" validate:"required"`
}

type EventResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OccursAt    time.Time `json:"occurs_at"`
	Location    string    `json:"location"`
	Capacity    int64     `json:"capacity"`
	OwnerID     string    `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	EventTypeId uuid.UUID `json:"event_type_id"`
}

func (r *CreateEventRequest) ToEntity(ownerID uuid.UUID) *domain.Event {
	return &domain.Event{
		Name:        r.Name,
		Description: r.Description,
		OccursAt:    r.OccursAt,
		Location:    r.Location,
		Capacity:    r.Capacity,
		EventTypeId: r.EventTypeId,
		OwnerId:     ownerID,
	}
}

func ToResponse(event *domain.Event) EventResponse {
	return EventResponse{
		ID:          event.ID.String(),
		Name:        event.Name,
		Description: event.Description,
		OccursAt:    event.OccursAt,
		Location:    event.Location,
		Capacity:    event.Capacity,
		OwnerID:     event.OwnerId.String(),
		CreatedAt:   event.CreatedAt,
		EventTypeId: event.EventTypeId,
	}
}

func (i *ListEventsInput) Normalize() {
	if i.Page < 1 {
		i.Page = 1
	}

	if i.Limit < 1 {
		i.Limit = 10 // Padrão recomendado: 10-25 itens
	} else if i.Limit > 100 {
		i.Limit = 100
	}
}
