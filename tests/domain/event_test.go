package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/visiontechw/ticketmaster/internal/domain"
)

func futureDate() time.Time {
	return time.Date(2026, 4, 20, 10, 0, 0, 0, time.UTC)
}
func TestNewEvent_Success(t *testing.T) {
	is := assert.New(t)
	ownerID := uuid.New()
	occursAt := futureDate()

	event, err := domain.NewEvent(ownerID, "  Go Workshop  ", "Learn Go basics", "Remote", 100, occursAt)

	require.NoError(t, err)
	is.NotNil(event)
	is.Equal("Go Workshop", event.Name, "Deve aplicar TrimSpace no nome")
	is.Equal(ownerID, event.OwnerId)
	is.Equal(int64(100), event.Capacity)
	is.True(event.OccursAt.Equal(occursAt))
	is.NotEmpty(event.ID, "Deve gerar um ID na base")
}

func TestNewEvent_Errors(t *testing.T) {
	is := assert.New(t)
	validOwner := uuid.New()
	pastDate := time.Now().Add(-48 * time.Hour)

	tests := []struct {
		name        string
		ownerId     uuid.UUID
		eventName   string
		location    string
		capacity    int64
		occursAt    time.Time
		expectedErr string
	}{
		{"Owner vazio", uuid.Nil, "Event", "Local", 10, futureDate(), "owner id is required"},
		{"Nome vazio", validOwner, "  ", "Local", 10, futureDate(), "event name is required"},
		{"Local vazio", validOwner, "Event", "  ", 10, futureDate(), "location is required"},
		{"Capacidade zero", validOwner, "Event", "Local", 0, futureDate(), "capacity must be greater than zero"},
		{"Capacidade negativa", validOwner, "Event", "Local", -5, futureDate(), "capacity must be greater than zero"},
		{"Data no passado", validOwner, "Event", "Local", 10, pastDate, "event date cannot be in the past"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := domain.NewEvent(tt.ownerId, tt.eventName, "desc", tt.location, tt.capacity, tt.occursAt)
			is.Nil(event)
			is.Error(err)
			is.Equal(tt.expectedErr, err.Error())
		})
	}
}

func TestEvent_BusinessRules(t *testing.T) {
	is := assert.New(t)
	ownerID := uuid.New()
	event, _ := domain.NewEvent(ownerID, "Tech Talk", "Desc", "Office", 50, futureDate())

	t.Run("CanEdit: Deve validar dono corretamente", func(t *testing.T) {
		is.True(event.CanEdit(ownerID))
		is.False(event.CanEdit(uuid.New()), "Outro usuário não pode editar")
	})

	t.Run("CanBeCanceled: Deve validar janela de 24h", func(t *testing.T) {
		// Evento em 2026 (longe) -> Pode cancelar
		is.True(event.CanBeCanceled())

		// Evento para daqui a 5 horas -> Não pode cancelar
		tightEvent, _ := domain.NewEvent(ownerID, "Quick Meet", "D", "L", 5, time.Now().Add(5*time.Hour))
		is.False(tightEvent.CanBeCanceled())
	})
}

func TestEvent_Updates(t *testing.T) {
	is := assert.New(t)
	event, _ := domain.NewEvent(uuid.New(), "Original", "Original Desc", "Local A", 10, futureDate())

	t.Run("UpdateDescription: Deve atualizar e mudar UpdatedAt", func(t *testing.T) {
		oldUpdate := event.UpdatedAt
		time.Sleep(time.Millisecond) // Garante diferença no clock

		err := event.UpdateDescription("Nova Descrição")
		is.NoError(err)
		is.Equal("Nova Descrição", event.Description)
		is.True(event.UpdatedAt.After(oldUpdate))
	})

	t.Run("UpdateLocation: Deve validar entrada vazia", func(t *testing.T) {
		err := event.UpdateLocation("  ")
		is.Error(err)
		is.Contains(err.Error(), "cannot be empty")
	})

	t.Run("UpdateCapacity: Deve validar valor positivo", func(t *testing.T) {
		err := event.UpdateCapacity(0)
		is.Error(err)
		is.Equal("capacity must be positive", err.Error())

		err = event.UpdateCapacity(200)
		is.NoError(err)
		is.Equal(int64(200), event.Capacity)
	})

	t.Run("Reschedule: Deve atualizar data", func(t *testing.T) {
		newDate := futureDate().Add(24 * time.Hour)
		event.Reschedule(newDate)
		is.True(event.OccursAt.Equal(newDate))
	})
}
