package domain

import "errors"

var (
	// Erros de Usuário
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailAlreadyUsed = errors.New("email already in use")
	ErrInvalidPassword  = errors.New("invalid password")

	// Erros de Pedido (Order)
	ErrOrderNotFound   = errors.New("order not found")
	ErrOrderNotPending = errors.New("order is not in pending status")
	ErrEmptyOrder      = errors.New("order must have at least one item")

	// Erros de Evento
	ErrEventFull     = errors.New("event has no more capacity")
	ErrEventPastDate = errors.New("event date cannot be in the past")
)
