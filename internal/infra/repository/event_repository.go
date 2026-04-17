package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/visiontechw/ticketmaster/internal/domain"
	"gorm.io/gorm"
)

type EventRepository struct {
	gormDB *gorm.DB
	sqlxDB *sqlx.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	sqlDB, err := db.DB()
	if err != nil {
		panic("falha ao extrair sql.DB do gorm: " + err.Error())
	}

	// Criamos a instância do sqlx (ajuste o driver se não for postgres)
	sxDB := sqlx.NewDb(sqlDB, "postgres")

	return &EventRepository{
		gormDB: db,
		sqlxDB: sxDB,
	}
}

func (r *EventRepository) GetEventsByStatusPaginated(ctx context.Context, active bool, page, limit int) ([]*domain.Event, error) {

	offset := (page - 1) * limit

	query := `SELECT 
     		   id, name, description, occurs_at, location, capacity, owner_id, event_type_id 			   
			   FROM events 
			   WHERE deleted_at IS NULL`

	var args []interface{}
	if active {
		query += ` AND occurs_at >= $1 ORDER BY occurs_at ASC`
		args = append(args, time.Now())
	} else {
		query += ` AND occurs_at <= $1 ORDER BY occurs_at DESC`
		args = append(args, time.Now())
	}

	query += ` LIMIT $2 OFFSET $3`
	args = append(args, limit, offset)

	var events []*domain.Event
	err := r.sqlxDB.SelectContext(ctx, &events, query, args...)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar eventos: %w", err)
	}

	return events, nil

}

func (r *EventRepository) Create(ctx context.Context, event *domain.Event) (*domain.Event, error) {

	result := r.gormDB.WithContext(ctx).Create(event)
	if result.Error != nil {
		return nil, fmt.Errorf("erro gorm ao criar evento: %w", result.Error)
	}

	return event, nil
}

func (r *EventRepository) FindByID(ctx context.Context, eventId uuid.UUID) (*domain.Event, error) {
	var event domain.Event

	dbquery := ` SELECT id, name, description, occurs_at, location, capacity, owner_id, event_type_id 	   
			   FROM events 
              WHERE id = $1 AND deleted_at IS NULL LIMIT 1`

	err := r.sqlxDB.GetContext(ctx, &event, dbquery, eventId)

	if err != nil {
		// No sqlx, se não achar nada, ele retorna sql.ErrNoRows
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("evento não encontrado: %w", err)
		}
		return nil, fmt.Errorf("erro ao executar select no banco: %w", err)
	}

	return &event, nil
}

func (r *EventRepository) Update(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	result := r.gormDB.WithContext(ctx).Save(event)
	if result.Error != nil {
		return nil, fmt.Errorf("erro ao atualizar evento: %w", result.Error)
	}

	return event, nil
}

func (r *EventRepository) Delete(ctx context.Context, eventId uuid.UUID, userId uuid.UUID) error {
	// First check ownership
	var event domain.Event
	if err := r.gormDB.WithContext(ctx).Where("id = ? AND owner_id = ?", eventId, userId).First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("evento não encontrado ou você não tem permissão para deletar")
		}
		return fmt.Errorf("erro ao verificar evento: %w", err)
	}

	// Soft delete using GORM
	result := r.gormDB.WithContext(ctx).Delete(&event)
	if result.Error != nil {
		return fmt.Errorf("erro ao deletar evento: %w", result.Error)
	}

	return nil
}
