package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kishanranaghosh/ecom-go-api/internal/config"
	"github.com/kishanranaghosh/ecom-go-api/internal/database"
	internal "github.com/kishanranaghosh/ecom-go-api/internal/logger"
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

	log.Info().Msgf("Starting server on port %s in %s mode", cfg.Server.Port, cfg.Server.GinMode)
	
}
