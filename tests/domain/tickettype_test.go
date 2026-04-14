package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/visiontechw/ticketmaster/internal/domain"
)

func TestNewTicketType_Create(t *testing.T) {
	is := assert.New(t)
	ticket, err := domain.NewTicketType(uuid.New(), "Standard", 100.5, 50, time.Now(),
		time.Now().Add(time.Hour))

	is.NoError(err)
	is.Equal("Standard", ticket.Name)
	is.Equal(100.5, ticket.Price)
	is.Equal(50, ticket.Capacity)

}

func TestNewTicketType_Errors(t *testing.T) {
	is := assert.New(t)
	eventID := uuid.New()
	now := time.Now()

	tests := []struct {
		name        string
		ticketName  string
		price       float64
		capacity    int
		start       time.Time
		end         time.Time
		expectedErr string
	}{
		{"Preço negativo", "VIP", -10.0, 100, now, now.Add(time.Hour), "price cannot be negative"},
		{"Capacidade zero", "VIP", 50.0, 0, now, now.Add(time.Hour), "capacity must be greater than zero"},
		{"Data invertida", "VIP", 50.0, 100, now.Add(time.Hour), now, "sales end date cannot be before sales start date"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ticket, err := domain.NewTicketType(eventID, tt.ticketName, tt.price, tt.capacity, tt.start, tt.end)
			is.Nil(ticket)
			is.Error(err)
			is.Equal(tt.expectedErr, err.Error())
		})
	}
}

func TestTicketType_SalesPeriod(t *testing.T) {
	is := assert.New(t)

	t.Run("Deve retornar true quando as vendas estão abertas", func(t *testing.T) {
		start := time.Now().Add(-1 * time.Hour)
		end := time.Now().Add(1 * time.Hour)
		ticket, _ := domain.NewTicketType(uuid.New(), "Early Bird", 100, 50, start, end)

		is.True(ticket.IsSalesOpen())
		is.False(ticket.IsSoldOut(10))
	})

	t.Run("Deve identificar Sold Out corretamente", func(t *testing.T) {
		ticket, _ := domain.NewTicketType(uuid.New(), "Standard", 100, 50, time.Now(), time.Now().Add(time.Hour))

		is.True(ticket.IsSoldOut(50))
		is.True(ticket.IsSoldOut(51))
		is.False(ticket.IsSoldOut(49))
	})

}
