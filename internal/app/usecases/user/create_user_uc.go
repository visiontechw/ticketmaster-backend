package userusecases

import (
	"context"

	"github.com/visiontechw/ticketmaster/internal/app/dto"
	"github.com/visiontechw/ticketmaster/internal/domain"
	"github.com/visiontechw/ticketmaster/internal/domain/repositories"
	"github.com/visiontechw/ticketmaster/internal/domain/services"
)

type CreateUserUseCase struct {
	userRepository repositories.UserRepository
	passwordHasher services.PasswordHasher
}

func NewCreateUserUseCases(repo repositories.UserRepository, hasher services.PasswordHasher) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository: repo,
		passwordHasher: hasher,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input dto.CreateUserInput) (dto.CreateUserOutput, error) {
	existingUser, _ := uc.userRepository.FindByEmail(ctx, input.Email)

	if existingUser != nil {
		return dto.CreateUserOutput{}, domain.ErrEmailAlreadyUsed
	}

	passwordHash, err := uc.passwordHasher.Hash(input.Password)
	if err != nil {
		return dto.CreateUserOutput{}, err
	}

	user, err := domain.NewUser(input.Name, input.Email, passwordHash)
	if err != nil {
		return dto.CreateUserOutput{}, err
	}

	if err := uc.userRepository.Create(ctx, user); err != nil {
		return dto.CreateUserOutput{}, err
	}

	return dto.CreateUserOutput{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
