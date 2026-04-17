package eventsusecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/visiontechw/ticketmaster/internal/domain/repositories"
)

type DeleteEventUseCase struct {
	eventRepository repositories.EventRepository
}

func NewDeleteEventUseCase(repo repositories.EventRepository) *DeleteEventUseCase {
	return &DeleteEventUseCase{
		eventRepository: repo,
	}
}

func (uc *DeleteEventUseCase) Execute(ctx context.Context, eventId uuid.UUID, userId uuid.UUID) error {
	return uc.eventRepository.Delete(ctx, eventId, userId)
}