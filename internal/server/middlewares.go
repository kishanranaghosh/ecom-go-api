package server

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kishanranaghosh/ecom-go-api/internal/models"
	"github.com/kishanranaghosh/ecom-go-api/internal/utils"
)

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header is missing")
			c.Abort()
			return
		}
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.UnauthorizedResponse(c, "Invalid Authorization header format")
			c.Abort()
			return
		}
		tokenString := tokenParts[1]
		claims, err := utils.ValidateToken(tokenString, s.JWTConfig.SecretKey)
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid or expired token")
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func (s *Server) adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != string(models.UserRoleAdmin) {
			utils.ForbiddenResponse(c, "Admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}
