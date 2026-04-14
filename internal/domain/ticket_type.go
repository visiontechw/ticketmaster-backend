package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TicketType struct {
	Base
	EventId uuid.UUID
	Name    string //VIP or Standart or anywhere
	Price   float64

	SalesStart time.Time
	SalesEnd   time.Time

	Capacity int
}

func (t *TicketType) validate() error {
	if strings.TrimSpace(t.Name) == "" {
		return errors.New("ticket type name is required")
	}

	if t.Price < 0 {
		return errors.New("price cannot be negative")
	}

	if t.Capacity <= 0 {
		return errors.New("capacity must be greater than zero")
	}

	if t.SalesEnd.Before(t.SalesStart) {
		return errors.New("sales end date cannot be before sales start date")
	}

	return nil
}

func NewTicketType(eventId uuid.UUID, name string, price float64, capacity int, start, end time.Time) (*TicketType, error) {
	tt := &TicketType{
		Base:       NewBase(),
		EventId:    eventId,
		Name:       strings.TrimSpace(name),
		Price:      price,
		Capacity:   capacity,
		SalesStart: start,
		SalesEnd:   end,
	}

	if err := tt.validate(); err != nil {
		return nil, err
	}

	return tt, nil
}

func (t *TicketType) IsSalesOpen() bool {
	now := time.Now()
	return now.After(t.SalesStart) && now.Before(t.SalesEnd)
}

func (t *TicketType) IsSoldOut(soldCount int) bool {
	return soldCount >= t.Capacity
}

func (t *TicketType) ExtendSales(newEnd time.Time) error {
	if newEnd.Before(t.SalesStart) {
		return errors.New("new end date cannot be before sales start")
	}
	t.SalesEnd = newEnd
	t.UpdatedAt = time.Now()
	return nil
}
