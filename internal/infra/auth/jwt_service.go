package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	secretKey []byte
	issuer    string
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{
		secretKey: []byte(secret),
		issuer:    "ticketfacil-api",
	}
}

func (s *JWTService) Generate(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"iss": s.issuer,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *JWTService) Validate(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid token claims")
	}

	sub, _ := claims["sub"].(string)
	return uuid.Parse(sub)
}
