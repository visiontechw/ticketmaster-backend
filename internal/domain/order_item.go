package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	Base
	OrderID      uuid.UUID
	TicketTypeID uuid.UUID
	Quantity     int
	UnitPrice    float64
}

func (oi *OrderItem) Validate() error {
	if oi.OrderID == uuid.Nil {
		return errors.New("order_id is required")
	}

	if oi.TicketTypeID == uuid.Nil {
		return errors.New("ticket_type_id is required")
	}

	if oi.Quantity <= 0 {
		return errors.New("quantity must be at least 1")
	}

	if oi.UnitPrice < 0 {
		return errors.New("unit price cannot be negative")
	}

	return nil
}

func NewOrderItem(orderID, ticketTypeID uuid.UUID, quantity int, unitPrice float64) (*OrderItem, error) {
	oi := &OrderItem{
		Base:         NewBase(),
		OrderID:      orderID,
		TicketTypeID: ticketTypeID,
		Quantity:     quantity,
		UnitPrice:    unitPrice,
	}

	if err := oi.Validate(); err != nil {
		return nil, err
	}

	return oi, nil
}

func (oi *OrderItem) UpdateQuantity(newQuantity int) error {
	if newQuantity <= 0 {
		return errors.New("quantity must be at least 1")
	}
	oi.Quantity = newQuantity
	oi.UpdatedAt = time.Now()
	return nil
}
