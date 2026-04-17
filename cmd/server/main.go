package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/gzip"
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
	gin.Use(gzip.Gzip(gzip.DefaultCompression))
	web.SetupRoutes(gin, web.RouterConfig{
		UserHandler:  di.UserHandler,
		EventHandler: di.EventHandler,
		TokenManager: di.TokenManager,
	})

	srv := &http.Server{
		Addr:    cfg.ApiServerPort,
		Handler: gin,
	}

	go func() {
		log.Printf("Servidor rodando na porta %s...", cfg.ApiServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao subir o servidor: %s\n", err)
		}
	}()

	// 3. Canal para escutar sinais do Sistema Operacional (Interrupt, Kill, SIGTERM)
	quit := make(chan os.Signal, 1)
	// kill (sem param) envia syscall.SIGTERM
	// ctrl+c envia os.Interrupt
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Bloqueia até receber um sinal
	<-quit
	log.Println("Iniciando desligamento (Graceful Shutdown)...")

	// 4. Definir um tempo limite (timeout) para o encerramento
	// Isso evita que o servidor fique "pendurado" se houver conexões travadas
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 5. Finalizar o servidor HTTP
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Servidor forçado a encerrar:", err)
	}

	log.Println("Fechando conexões de banco de dados...")
	if err := sqlDB.Close(); err != nil {
		log.Printf("Erro ao fechar DB: %v", err)
	}

	log.Println("Servidor finalizado com sucesso.")

}
