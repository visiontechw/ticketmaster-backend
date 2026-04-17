package eventsusecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/visiontechw/ticketmaster/internal/app/dto"
	"github.com/visiontechw/ticketmaster/internal/domain/repositories"
)

type GetEventByIdUseCase struct {
	eventRepository repositories.EventRepository
}

func NewGetEventByIdUseCases(repo repositories.EventRepository) *GetEventByIdUseCase {
	return &GetEventByIdUseCase{
		eventRepository: repo,
	}
}

func (uc *GetEventByIdUseCase) Execute(ctx context.Context, eventId uuid.UUID) (dto.EventResponse, error) {
	event, err := uc.eventRepository.FindByID(ctx, eventId)
	if err != nil {
		return dto.EventResponse{}, err
	}

	return dto.ToResponse(event), nil

}
