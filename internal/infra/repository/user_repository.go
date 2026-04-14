package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/visiontechw/ticketmaster/internal/domain"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *GormUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error

	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return &user, nil
}

func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	// SELECT * FROM users WHERE email = ? AND deleted_at IS NULL LIMIT 1;
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return &user, nil
}

func (r *GormUserRepository) Update(ctx context.Context, user *domain.User) error {
	// O GORM atualizará apenas os campos que mudaram e o UpdatedAt automaticamente
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *GormUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Como temos o campo DeletedAt na Base, o GORM fará um Soft Delete automaticamente
	return r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", id).Error
}
