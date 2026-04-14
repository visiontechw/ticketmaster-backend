package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/visiontechw/ticketmaster/internal/domain"
)

func validTicket(t *testing.T) *domain.Ticket {
	t.Helper()
	ticket, err := domain.NewTicket(
		uuid.New(),
		uuid.New(),
		"123.456.789-00",
		"QR-ABC123",
		"HASH-XYZ999",
	)
	require.NoError(t, err)
	return ticket
}

func cancelledTicket(t *testing.T) *domain.Ticket {
	t.Helper()
	ticket := validTicket(t)
	require.NoError(t, ticket.Cancel())
	return ticket
}

func usedTicket(t *testing.T) *domain.Ticket {
	t.Helper()
	ticket := validTicket(t)
	require.NoError(t, ticket.CheckIn())
	return ticket
}

func TestNewTicket_Create(t *testing.T) {
	orderID := uuid.New()
	typeID := uuid.New()

	ticket, err := domain.NewTicket(orderID, typeID, "123.456.789-00", "QR-001", "HASH-001")

	require.NoError(t, err)
	assert.Equal(t, orderID, ticket.OrderID)
	assert.Equal(t, typeID, ticket.TicketTypeID)
	assert.Equal(t, "123.456.789-00", ticket.OwnerDocument)
	assert.Equal(t, "QR-001", ticket.QRCode)
	assert.Equal(t, "HASH-001", ticket.HashCode)
	assert.False(t, ticket.IsUsed)
	assert.NotEqual(t, uuid.Nil, ticket.ID)

}

func TestNewTicket_CheckIn(t *testing.T) {
	ticket := validTicket(t)

	err := ticket.CheckIn()

	assert.NoError(t, err)
	assert.True(t, ticket.IsUsed)
}

func TestCheckIn_AlreadyUsed(t *testing.T) {
	ticket := usedTicket(t)

	err := ticket.CheckIn()

	assert.ErrorContains(t, err, "already been used")
}

func TestCancel_Success(t *testing.T) {
	ticket := validTicket(t)

	err := ticket.Cancel()

	require.NoError(t, err)
	assert.NotNil(t, ticket.CancelledAt)
}

func TestCancel_AlreadyUsed(t *testing.T) {
	ticket := usedTicket(t)

	err := ticket.Cancel()

	assert.ErrorContains(t, err, "already been used")
}
