package userusecases

import (
	"context"
	"errors"

	"github.com/visiontechw/ticketmaster/internal/app/dto"
	"github.com/visiontechw/ticketmaster/internal/domain/repositories"
	"github.com/visiontechw/ticketmaster/internal/domain/services"
)

type LoginUseCase struct {
	userRepository repositories.UserRepository
	passwordHasher services.PasswordHasher
	tokenManager   services.TokenManager
}

func NewLoginUseCase(
	repo repositories.UserRepository,
	hasher services.PasswordHasher,
	tokenManager services.TokenManager) *LoginUseCase {
	return &LoginUseCase{
		userRepository: repo,
		passwordHasher: hasher,
		tokenManager:   tokenManager,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input dto.LoginInput) (dto.LoginOutput, error) {
	// 1. Busca o usuário pelo email
	user, err := uc.userRepository.FindByEmail(ctx, input.Email)
	if err != nil {
		return dto.LoginOutput{}, err
	}

	// Segurança: se o usuário não existe, retornamos erro genérico para evitar enumeração de contas
	if user == nil {
		return dto.LoginOutput{}, errors.New("invalid email or password")
	}

	if !user.Active {
		return dto.LoginOutput{}, errors.New("user account is deactivated")
	}

	// 2. Compara a senha enviada com o Hash salvo
	if !uc.passwordHasher.Compare(user.Password, input.Password) {
		return dto.LoginOutput{}, errors.New("invalid email or password")
	}

	// 3. Verifica se o usuário está ativo

	// 4. Gera o Token (JWT ou outro)
	token, err := uc.tokenManager.Generate(user.ID, user.Email)
	if err != nil {
		return dto.LoginOutput{}, errors.New("failed to generate access token")
	}

	return dto.LoginOutput{
		Token: token,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
