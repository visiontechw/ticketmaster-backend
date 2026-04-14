package auth

import "golang.org/x/crypto/bcrypt"

type BcryptService struct{}

func NewBcryptService() *BcryptService {
	return &BcryptService{}
}

func (s *BcryptService) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *BcryptService) Compare(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
