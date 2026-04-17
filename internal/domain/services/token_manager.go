package services

import (
	"github.com/google/uuid"
	"github.com/visiontechw/ticketmaster/internal/infra/auth"
)

type TokenManager interface {
	Generate(userID uuid.UUID, email string) (string, error)
	Validate(tokenString string) (*auth.UserClaims, error)
}
