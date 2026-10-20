package server

import (
	"github.com/gin-gonic/gin"
	"github.com/kishanranaghosh/ecom-go-api/internal/config"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Server struct {
	config *config.ServerConfig
	db     *gorm.DB
	logger zerolog.Logger
}

func New(config *config.ServerConfig, db *gorm.DB, logger zerolog.Logger) *Server {
	return &Server{
		config: config,
		db:     db,
		logger: logger,
	}
}

// this is like member function of Server struct, so we can access the config, db and logger from this function
func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.New()

	// middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())

	// routes
	router.GET("/health", s.healthCheckHandler)
	return router
}

func (s *Server) healthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
