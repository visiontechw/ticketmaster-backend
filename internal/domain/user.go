package domain

import (
	"errors"
	"net/mail"
	"strings"
	"time"
)

type User struct {
	Base
	Name     string
	Email    string
	Password string
	Active   bool
}

func (u *User) validate() error {
	if strings.TrimSpace(u.Name) == "" {
		return errors.New("name is required")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.New("invalid email format")
	}

	if strings.TrimSpace(u.Password) == "" {
		return errors.New("Password is required")
	}

	if len(u.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}

func NewUser(name, email, passwordHash string) (*User, error) {
	user := &User{
		Base:     NewBase(),
		Name:     strings.TrimSpace(name),
		Email:    strings.ToLower(strings.TrimSpace(email)),
		Password: strings.TrimSpace(passwordHash),
		Active:   true,
	}

	if err := user.validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) Activate() {
	u.Active = true
	u.UpdatedAt = time.Now()
}

func (u *User) Deactivate() {
	u.Active = false
	u.UpdatedAt = time.Now()
}
