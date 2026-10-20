package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kishanranaghosh/ecom-go-api/internal/config"
	"github.com/kishanranaghosh/ecom-go-api/internal/database"
	internal "github.com/kishanranaghosh/ecom-go-api/internal/logger"
	"github.com/kishanranaghosh/ecom-go-api/internal/server"
)

func main() {
	log := internal.New().With().Logger()
	cfg, err := config.Load()

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	log.Info().Msgf("Configuration loaded: %+v", cfg)

	db, err := database.New(&cfg.Database)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
	}
	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get database instance")
	}
	defer mainDB.Close()

	log.Info().Msg("Database connection established")
	gin.SetMode(cfg.Server.GinMode)
	srv := server.New(&cfg.Server, db, log)

	router := srv.SetupRoutes()

	httServer := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		log.Info().Msgf("Starting server on port %s in %s mode", cfg.Server.Port, cfg.Server.GinMode)
		if err := httServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httServer.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")

}
