package web

import (
	"github.com/gin-gonic/gin"
	"github.com/visiontechw/ticketmaster/internal/app/handlers"
	"github.com/visiontechw/ticketmaster/internal/app/middleware"
	"github.com/visiontechw/ticketmaster/internal/domain/services"
)

// RouterConfig centraliza todos os handlers para facilitar a injeção
type RouterConfig struct {
	UserHandler  *handlers.UserHandler
	EventHandler *handlers.EventHandler
	TokenManager services.TokenManager
}

func SetupRoutes(r *gin.Engine, config RouterConfig) {
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		// Rotas Públicas (Auth)
		auth := api.Group("/auth")
		{
			auth.POST("/register", config.UserHandler.Create)
			auth.POST("/login", config.UserHandler.Login)
		}

		// Rotas Públicas (Events)
		publicEvents := api.Group("/events")
		{
			publicEvents.GET("/", config.EventHandler.ListEvents)
			publicEvents.GET("/:id", config.EventHandler.GetById)
		}

		// --- GRUPO PROTEGIDO ---
		protected := api.Group("/")
		protected.Use(middleware.TokenCatcherMiddleware(config.TokenManager))
		{
			// 1. Rota de Perfil do Usuário (Adicionada aqui)
			user := protected.Group("/user")
			{
				user.GET("/profile", config.UserHandler.Me)
			}

			// 2. Rotas de Admin/Eventos protegidos
			adminEvents := protected.Group("/events")
			{
				adminEvents.POST("/", config.EventHandler.Create)
				adminEvents.PUT("/:id", config.EventHandler.Update)
				adminEvents.DELETE("/:id", config.EventHandler.Delete)
			}
		}
	}
}
