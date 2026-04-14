package container

import (
	"github.com/jmoiron/sqlx"
	"github.com/visiontechw/ticketmaster/internal/app/handlers"
	usecases "github.com/visiontechw/ticketmaster/internal/app/usecases/user"
	"github.com/visiontechw/ticketmaster/internal/infra/auth"
	"github.com/visiontechw/ticketmaster/internal/infra/config"
	"github.com/visiontechw/ticketmaster/internal/infra/repository"
	"gorm.io/gorm"
)

type DependencyContainer struct {
	UserHandler *handlers.UserHandler
	// EventHandler *handlers.EventHandler
}

// NewDependencyContainer faz todo o "wire-up" que estava na main
func NewDependencyContainer(gormDB *gorm.DB, sqlxDB *sqlx.DB, cfg *config.Config) *DependencyContainer {
	// 1. Repositories
	userRepo := repository.NewGormUserRepository(gormDB)

	// services
	passwordHasher := auth.NewBcryptService()
	tokenManager := auth.NewJWTService(cfg.JwtKey)

	// 2. Use Cases
	createUserUC := usecases.NewCreateUserUseCases(userRepo, passwordHasher)
	loginUC := usecases.NewLoginUseCase(userRepo, passwordHasher, tokenManager)

	// 3. Handlers
	userHandler := handlers.NewUserHandler(createUserUC, loginUC)

	return &DependencyContainer{
		UserHandler: userHandler,
	}
}
