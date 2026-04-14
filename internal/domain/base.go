package domain

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID        uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt time.Time `db:"created_at" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `db:"updated_at" gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt time.Time `db:"deleted_at" gorm:"column:deleted_at;index"`
}

func NewBase() Base {
	now := time.Now()
	return Base{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}
