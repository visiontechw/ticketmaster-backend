package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	// Status iniciais
	StatusPending OrderStatus = "pending" // Aguardando pagamento
	// Status finais de sucesso
	StatusPaid      OrderStatus = "paid"      // Pagamento confirmado
	StatusCompleted OrderStatus = "completed" // Evento já ocorreu e pedido foi finalizado
	// Status de falha/cancelamento
	StatusCanceled OrderStatus = "canceled" // Cancelado pelo usuário ou sistema
	StatusFailed   OrderStatus = "failed"   // Erro no processamento do pagamento
	StatusRefunded OrderStatus = "refunded" // Valor devolvido ao cliente
)

type Order struct {
	Base
	UserID      uuid.UUID
	EventID     uuid.UUID
	TotalAmount float64
	Status      OrderStatus
	Items       []OrderItem
}

func (o *Order) Validate() error {
	if o.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}

	if o.EventID == uuid.Nil {
		return errors.New("event_id is required")
	}

	if len(o.Items) == 0 {
		return errors.New("order must have at least one item")
	}

	// Validação de consistência financeira
	var calculatedTotal float64
	for _, item := range o.Items {
		calculatedTotal += item.SubTotal()
	}

	// Verificamos se o total informado bate com a soma dos itens
	// Nota: Em sistemas reais, costuma-se usar uma pequena margem de tolerância para float
	if o.TotalAmount != calculatedTotal {
		return errors.New("total amount does not match the sum of items")
	}

	if o.Status == "" {
		return errors.New("order status is required")
	}

	return nil
}

func (o *Order) CalculateTotal() float64 {
	var total float64
	for _, item := range o.Items {
		total += item.SubTotal()
	}
	return total
}

func NewOrder(userID, eventID uuid.UUID, items []OrderItem) (*Order, error) {
	order := &Order{
		Base:    NewBase(),
		UserID:  userID,
		EventID: eventID,
		Status:  StatusPending, // Todo pedido nasce como pendente
		Items:   items,
	}
	order.TotalAmount = order.CalculateTotal()

	if err := order.Validate(); err != nil {
		return nil, err
	}

	return order, nil
}

func (oi *OrderItem) SubTotal() float64 {
	return float64(oi.Quantity) * oi.UnitPrice
}

func (o *Order) MarkAsPaid() error {
	if o.Status != StatusPending {
		return errors.New("only pending orders can be marked as paid")
	}
	o.Status = StatusPaid
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) Cancel() error {
	if o.Status == StatusCompleted || o.Status == StatusPaid {
		return errors.New("cannot cancel a finished or paid order")
	}
	o.Status = StatusCanceled
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) Refund() error {
	if o.Status != StatusPaid {
		return errors.New("only paid orders can be refunded")
	}
	o.Status = StatusRefunded
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) RemoveItem(ticketTypeID uuid.UUID) error {
	if o.Status != StatusPending {
		return errors.New("cannot remove items from an order that is not pending")
	}

	index := -1
	for i, item := range o.Items {
		if item.TicketTypeID == ticketTypeID {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("item not found in order")
	}

	o.Items = append(o.Items[:index], o.Items[index+1:]...)

	if len(o.Items) == 0 {
		return o.Cancel()
	}

	o.TotalAmount = o.CalculateTotal()
	o.UpdatedAt = time.Now()

	return nil
}
