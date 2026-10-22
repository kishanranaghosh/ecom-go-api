package server

import (
	"github.com/gin-gonic/gin"
	"github.com/kishanranaghosh/ecom-go-api/internal/config"
	"github.com/kishanranaghosh/ecom-go-api/internal/services"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Server struct {
	serverConfig   *config.ServerConfig
	JWTConfig      *config.JWTConfig
	db             *gorm.DB
	logger         zerolog.Logger
	authService    services.AuthServiceInterface
	productService services.ProductServiceInterface
	uploadService  services.UploadServiceInterface
	userService    services.UserServiceInterface
	cartService    services.CartServiceInterface
	orderService   services.OrderServiceInterface
}

func New(
	cfg *config.Config,
	logger *zerolog.Logger,
	authService services.AuthServiceInterface,
	productService services.ProductServiceInterface,
	userService services.UserServiceInterface,
	uploadService services.UploadServiceInterface,
	cartService services.CartServiceInterface,
	orderService services.OrderServiceInterface,
) *Server {
	return &Server{
		serverConfig:   &cfg.Server,
		JWTConfig:      &cfg.JWT,
		logger:         logger.With().Str("component", "server").Logger(),
		authService:    authService,
		productService: productService,
		userService:    userService,
		uploadService:  uploadService,
		cartService:    cartService,
		orderService:   orderService,
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

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{ //nolint:gocritic // I need this for readability
			auth.POST("/register", s.register)
			auth.POST("/login", s.login)
			auth.POST("/refresh", s.refreshToken)
			auth.POST("/logout", s.logout)

		}
		protected := api.Group("/")
		protected.Use(s.authMiddleware())
		{
			// User routes
			users := protected.Group("/users")
			{
				userRoutes := users
				userRoutes.GET("/profile", s.getProfile)
				userRoutes.PUT("/profile", s.updateProfile)
			}

			// category routes
			categories := protected.Group("/categories")
			{
				categoryRoute := categories
				categoryRoute.POST("/", s.adminMiddleware(), s.createCategory)
				categoryRoute.PUT("/:id", s.adminMiddleware(), s.updateCategory)
				categoryRoute.DELETE("/:id", s.adminMiddleware(), s.deleteCategory)
			}

			// product routes
			products := protected.Group("/products")
			{
				productRoutes := products
				productRoutes.POST("/", s.adminMiddleware(), s.createProduct)
				productRoutes.PUT("/:id", s.adminMiddleware(), s.updateProduct)
				productRoutes.DELETE("/:id", s.adminMiddleware(), s.deleteProduct)
				productRoutes.POST("/:id/images", s.adminMiddleware(), s.uploadProductImage)

			}

		}
	}
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
