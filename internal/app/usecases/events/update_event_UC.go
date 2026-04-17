package eventsusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/visiontechw/ticketmaster/internal/app/dto"
	"github.com/visiontechw/ticketmaster/internal/domain/repositories"
)

type UpdateEventUseCase struct {
	eventRepository repositories.EventRepository
}

func NewUpdateEventUseCase(repo repositories.EventRepository) *UpdateEventUseCase {
	return &UpdateEventUseCase{
		eventRepository: repo,
	}
}

func (uc *UpdateEventUseCase) Execute(ctx context.Context, eventId uuid.UUID, input dto.UpdateEventRequest) (dto.EventResponse, error) {
	// Get current event
	event, err := uc.eventRepository.FindByID(ctx, eventId)
	if err != nil {
		return dto.EventResponse{}, err
	}

	// Extract user_id from context (authorization)
	userId, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		return dto.EventResponse{}, errors.New("unauthorized: user not found")
	}

	// Check ownership
	if !event.CanEdit(userId) {
		return dto.EventResponse{}, errors.New("unauthorized: you cannot edit this event")
	}

	// Update fields
	event.Name = input.Name
	event.Description = input.Description
	event.OccursAt = input.OccursAt
	event.Location = input.Location
	event.Capacity = input.Capacity
	event.EventTypeId = input.EventTypeId

	// Save
	updated, err := uc.eventRepository.Update(ctx, event)
	if err != nil {
		return dto.EventResponse{}, err
	}

	return dto.ToResponse(updated), nil
}