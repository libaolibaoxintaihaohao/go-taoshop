package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"taoshop/internal/service"
)

func Auth(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing authorization header"})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		if token == header {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization scheme"})
			return
		}

		userID, err := authService.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
