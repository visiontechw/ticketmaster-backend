package userusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/visiontechw/ticketmaster/internal/app/dto"
	"github.com/visiontechw/ticketmaster/internal/domain/repositories"
)

type GetUserByIdUseCase struct {
	userRepository repositories.UserRepository
}

func NewGetUserByIdUseCase(
	repo repositories.UserRepository) *GetUserByIdUseCase {
	return &GetUserByIdUseCase{
		userRepository: repo,
	}
}

func (uc *GetUserByIdUseCase) Execute(ctx context.Context) (dto.CreateUserOutput, error) {

	ownerID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		return dto.CreateUserOutput{}, errors.New("unauthorized: userid not found")
	}

	user, err := uc.userRepository.FindByID(ctx, ownerID)
	if err != nil {
		return dto.CreateUserOutput{}, err
	}

	if user == nil {
		return dto.CreateUserOutput{}, errors.New("user not found")
	}

	return dto.CreateUserOutput{
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
