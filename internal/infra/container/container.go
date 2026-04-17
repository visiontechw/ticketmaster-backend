package container

import (
	"github.com/jmoiron/sqlx"
	"github.com/visiontechw/ticketmaster/internal/app/handlers"
	eventsusecases "github.com/visiontechw/ticketmaster/internal/app/usecases/events"
	userusecases "github.com/visiontechw/ticketmaster/internal/app/usecases/user"
	"github.com/visiontechw/ticketmaster/internal/domain/services"
	"github.com/visiontechw/ticketmaster/internal/infra/auth"
	"github.com/visiontechw/ticketmaster/internal/infra/config"
	"github.com/visiontechw/ticketmaster/internal/infra/repository"
	"gorm.io/gorm"
)

type DependencyContainer struct {
	UserHandler  *handlers.UserHandler
	EventHandler *handlers.EventHandler
	TokenManager services.TokenManager
}

func NewDependencyContainer(gormDB *gorm.DB, sqlxDB *sqlx.DB, cfg *config.Config) *DependencyContainer {
	// 1. Repositories
	userRepo := repository.NewGormUserRepository(gormDB)
	eventRepo := repository.NewEventRepository(gormDB)

	// services
	passwordHasher := auth.NewBcryptService()
	tokenManager := auth.NewJWTService(cfg.JwtKey)

	// 2. Use Cases
	createUserUC := userusecases.NewCreateUserUseCases(userRepo, passwordHasher)
	loginUC := userusecases.NewLoginUseCase(userRepo, passwordHasher, tokenManager)
	userById := userusecases.NewGetUserByIdUseCase(userRepo)

	getEventsUC := eventsusecases.NewGetEventsUseCases(eventRepo)
	createEventUC := eventsusecases.NewCreateventUseCases(eventRepo)
	getEventById := eventsusecases.NewGetEventByIdUseCases(eventRepo)
	updateEventUC := eventsusecases.NewUpdateEventUseCase(eventRepo)
	deleteEventUC := eventsusecases.NewDeleteEventUseCase(eventRepo)

	// 3. Handlers
	userHandler := handlers.NewUserHandler(createUserUC, loginUC, userById)
	eventHandler := handlers.NewEventHandler(getEventsUC, createEventUC, getEventById, updateEventUC, deleteEventUC)

	return &DependencyContainer{
		UserHandler:  userHandler,
		EventHandler: eventHandler,
		TokenManager: tokenManager,
	}
}
