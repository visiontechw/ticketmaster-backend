package eventsusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/visiontechw/ticketmaster/internal/app/dto"
	"github.com/visiontechw/ticketmaster/internal/domain/repositories"
)

type CreateventUseCase struct {
	eventRepository repositories.EventRepository
}

func NewCreateventUseCases(repo repositories.EventRepository) *CreateventUseCase {
	return &CreateventUseCase{
		eventRepository: repo,
	}
}

func (uc *CreateventUseCase) Execute(ctx context.Context, input dto.CreateEventRequest) (dto.EventResponse, error) {

	ownerID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		return dto.EventResponse{}, errors.New("unauthorized: owner id not found")
	}

	event := input.ToEntity(ownerID)

	evt, err := uc.eventRepository.Create(ctx, event)

	if err != nil {
		return dto.EventResponse{}, errors.New("Erro ao criar evento")

	}

	return dto.ToResponse(evt), nil

}
