package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/visiontechw/ticketmaster/internal/domain"
)

func TestOrder_StateTransitions(t *testing.T) {
	is := assert.New(t)

	// Setup básico
	userId := uuid.New()
	eventId := uuid.New()
	typeticketId := uuid.New()

	items := []domain.OrderItem{{Quantity: 1, TicketTypeID: typeticketId, UnitPrice: 100.0}}
	order, _ := domain.NewOrder(userId, eventId, items)

	t.Run("Deve impedir pagamento de pedido cancelado", func(t *testing.T) {
		order.Cancel()
		err := order.MarkAsPaid()
		is.Error(err)
		is.Equal(domain.StatusCanceled, order.Status)
	})

	t.Run("Deve permitir reembolso apenas de pedidos pagos", func(t *testing.T) {
		newOrder, _ := domain.NewOrder(uuid.New(), uuid.New(), items)

		// Tentar reembolsar pendente deve falhar
		is.Error(newOrder.Refund())

		// Pagar e depois reembolsar deve funcionar
		newOrder.MarkAsPaid()
		is.NoError(newOrder.Refund())
		is.Equal(domain.StatusRefunded, newOrder.Status)
	})

	t.Run("Deve validar se o cálculo do total está correto na criação", func(t *testing.T) {
		itemsMult := []domain.OrderItem{
			{Quantity: 2, UnitPrice: 50.0}, // 100
			{Quantity: 1, UnitPrice: 30.0}, // 30
		}
		o, err := domain.NewOrder(uuid.New(), uuid.New(), itemsMult)
		is.NoError(err)
		is.Equal(130.0, o.TotalAmount)
	})
}
func TestOrder_RemoveItem_AutoCancel(t *testing.T) {
	is := assert.New(t)

	ticketID := uuid.New()
	items := []domain.OrderItem{
		{TicketTypeID: ticketID, Quantity: 1, UnitPrice: 100.0},
	}

	order, _ := domain.NewOrder(uuid.New(), uuid.New(), items)

	t.Run("Deve cancelar o pedido automaticamente ao remover o último item", func(t *testing.T) {
		err := order.RemoveItem(ticketID)

		is.NoError(err)
		is.Len(order.Items, 0)
		is.Equal(domain.StatusCanceled, order.Status, "O status deve mudar para canceled")
	})
}
