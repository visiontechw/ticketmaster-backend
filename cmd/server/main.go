package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/visiontechw/ticketmaster/internal/app/web"
	"github.com/visiontechw/ticketmaster/internal/infra/config"
	"github.com/visiontechw/ticketmaster/internal/infra/container"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Não foi possível carregar as configurações:", err)
	}

	gormDB, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		PrepareStmt:    true,
		QueryFields:    true,
		TranslateError: true,
	})
	if err != nil {
		log.Fatal("Falha ao conectar via GORM:", err)
	}

	sqlDB, _ := gormDB.DB()
	sqlxDB := sqlx.NewDb(sqlDB, "postgres")

	di := container.NewDependencyContainer(gormDB, sqlxDB, cfg)

	gin := gin.Default()
	web.SetupRoutes(gin, web.RouterConfig{
		UserHandler: di.UserHandler,
	})

	log.Printf("Servidor rodando na porta %s...", cfg.ApiServerPort)
	if err := gin.Run(cfg.ApiServerPort); err != nil {
		log.Fatal("Erro ao subir o servidor:", err)
	}

}
