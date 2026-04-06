package middleware

import (
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/utils/apperrors"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a Gin middleware that validates the "api_key" header.
func AuthMiddleware(expectedKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("api_key")
		if apiKey == "" {
			c.JSON(apperrors.MapToHTTP(apperrors.NewUnauthorized("API key is required")))
			c.Abort()
			return
		}

		if apiKey != expectedKey {
			c.JSON(apperrors.MapToHTTP(apperrors.NewUnauthorized("Invalid API key")))
			c.Abort()
			return
		}

		c.Next()
	}
}
