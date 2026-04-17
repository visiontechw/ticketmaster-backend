package eventsusecases

import (
	"context"
	"fmt"

	"github.com/visiontechw/ticketmaster/internal/app/dto"
	"github.com/visiontechw/ticketmaster/internal/domain/repositories"
)

type GetEventsUseCase struct {
	eventRepository repositories.EventRepository
}

func NewGetEventsUseCases(repo repositories.EventRepository) *GetEventsUseCase {
	return &GetEventsUseCase{
		eventRepository: repo,
	}
}

func (uc *GetEventsUseCase) Execute(ctx context.Context, input dto.ListEventsInput) ([]dto.EventResponse, error) {
	input.Normalize()
	events, err := uc.eventRepository.GetEventsByStatusPaginated(ctx, input.Active, input.Page, input.Limit)
	if err != nil {
		return nil, fmt.Errorf("usecase list events: %w", err)
	}

	responses := make([]dto.EventResponse, len(events))

	for i, event := range events {
		responses[i] = dto.ToResponse(event)
	}

	return responses, nil
}
