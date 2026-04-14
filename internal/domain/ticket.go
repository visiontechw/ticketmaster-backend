package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	Base
	OrderID       uuid.UUID
	TicketTypeID  uuid.UUID
	OwnerDocument string
	QRCode        string
	HashCode      string
	IsUsed        bool
	CancelledAt   *time.Time
}

func (t *Ticket) Validate() error {
	if t.OrderID == uuid.Nil {
		return errors.New("ticket must be linked to an order")
	}

	if t.TicketTypeID == uuid.Nil {
		return errors.New("ticket must have a ticket type")
	}

	if strings.TrimSpace(t.OwnerDocument) == "" {
		return errors.New("owner document is required")
	}

	if strings.TrimSpace(t.QRCode) == "" || strings.TrimSpace(t.HashCode) == "" {
		return errors.New("ticket identifiers (QRCode/HashCode) are required")
	}

	return nil
}

func NewTicket(orderID, ticketTypeID uuid.UUID, ownerDocument, qrCode, hashCode string) (*Ticket, error) {
	ticket := &Ticket{
		Base:          NewBase(),
		OrderID:       orderID,
		TicketTypeID:  ticketTypeID,
		OwnerDocument: strings.TrimSpace(ownerDocument),
		QRCode:        qrCode,
		HashCode:      hashCode,
		IsUsed:        false,
	}

	if err := ticket.Validate(); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (t *Ticket) CheckIn() error {
	if t.IsUsed {
		return errors.New("ticket has already been used")
	}

	t.IsUsed = true
	t.UpdatedAt = time.Now()

	return nil
}

func (t *Ticket) Cancel() error {
	if t.IsUsed {
		return errors.New("cannot cancel a ticket that has already been used")
	}
	now := time.Now()
	t.CancelledAt = &now
	t.UpdatedAt = now
	return nil
}

func (t *Ticket) String() string {
	return fmt.Sprintf("Ticket{ID:%s Order:%s Owner:%s Used:%v}",
		t.ID, t.OrderID, t.OwnerDocument, t.IsUsed)
}
