package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/visiontechw/ticketmaster/internal/domain"
)

func TestNewUser_Created(t *testing.T) {
	is := assert.New(t)
	t.Run("Deve criar um usuário com campos normalizados", func(t *testing.T) {
		user, err := domain.NewUser("  teste_user  ", "TEST@MAIL.COM", "!s3nh66!")

		require.NoError(t, err)
		is.NotNil(user)
		is.Equal("teste_user", user.Name)
		is.Equal("test@mail.com", user.Email)
		is.True(user.Active)
		is.NotEmpty(user.ID)
	})

}

func TestUser_Lifecycle(t *testing.T) {
	is := assert.New(t)

	// Setup inicial
	user, err := domain.NewUser("teste_user", "test@mail.com", "!s3nh66!")
	require.NoError(t, err)

	t.Run("Deve desativar o usuário e atualizar o timestamp", func(t *testing.T) {
		initialUpdate := user.UpdatedAt

		// Pequeno delay para garantir que o tempo mude (opcional em testes de unidade)
		time.Sleep(time.Millisecond)

		user.Deactivate()

		is.False(user.Active, "O status Active deve ser false")
		is.True(user.UpdatedAt.After(initialUpdate), "UpdatedAt deve ser maior que o tempo de criação")
	})

	t.Run("Deve reativar o usuário", func(t *testing.T) {
		user.Activate()
		is.True(user.Active, "O status Active deve voltar para true")
	})
}

//Para os cenários de erro, a melhor abordagem é o Table-Driven Testing.
//Isso permite testar todas as validações do seu método validate()

func TestNewUser_ValidationErrors(t *testing.T) {
	is := assert.New(t)

	tests := []struct {
		testName    string
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			testName:    "Erro: Nome vazio",
			name:        "   ",
			email:       "test@mail.com",
			password:    "!s3nh66!",
			expectedErr: "name is required",
		},
		{
			testName:    "Erro: Email inválido (sem @)",
			name:        "User Test",
			email:       "email_invalido.com",
			password:    "!s3nh66!",
			expectedErr: "invalid email format",
		},
		{
			testName:    "Erro: Email vazio",
			name:        "User Test",
			email:       "",
			password:    "!s3nh66!",
			expectedErr: "invalid email format",
		},
		{
			testName:    "Erro: Senha vazia",
			name:        "User Test",
			email:       "test@mail.com",
			password:    "    ",
			expectedErr: "Password is required", // Note o 'P' maiúsculo conforme seu código
		},
		{
			testName:    "Erro: Senha muito curta",
			name:        "User Test",
			email:       "test@mail.com",
			password:    "1234567",
			expectedErr: "password must be at least 8 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			user, err := domain.NewUser(tt.name, tt.email, tt.password)

			// Asserções para erro
			is.Error(err, "Deveria retornar erro para o cenário: %s", tt.testName)
			is.Nil(user, "O usuário deve ser nil quando houver erro de validação")
			is.Contains(err.Error(), tt.expectedErr, "A mensagem de erro deve ser a esperada")
		})
	}
}
