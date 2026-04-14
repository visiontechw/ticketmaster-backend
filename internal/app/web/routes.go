package web

import (
	"github.com/gin-gonic/gin"
	"github.com/visiontechw/ticketmaster/internal/app/handlers"
)

// RouterConfig centraliza todos os handlers para facilitar a injeção
type RouterConfig struct {
	UserHandler *handlers.UserHandler
	// EventHandler *EventHandler (Adicione conforme criar novos)
}

func SetupRoutes(r *gin.Engine, config RouterConfig) {
	// Middleware Global (Ex: Logger, Recovery, CORS)
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		// Rotas de Usuário
		auth := api.Group("/auth")
		{
			// Registro de novo usuário
			auth.POST("/register", config.UserHandler.Create)

			// Login: recebe email/senha e retorna o Token JWT
			auth.POST("/login", config.UserHandler.Login)
		}

		// Espaço para as próximas rotas:
		// events := api.Group("/events")
		// {
		//     events.GET("/", config.EventHandler.List)
		// }
	}
}
