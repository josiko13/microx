package middleware

import (
	"net/http"
	"strconv"

	"microx/internal/repository"

	"github.com/gin-gonic/gin"
)

const (
	UserIDHeader = "X-User-ID"
	UserIDKey    = "user_id"
)

// AuthMiddleware extrae el ID de usuario del header X-User-ID (validación básica)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetHeader(UserIDHeader)

		if userIDStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "X-User-ID header is required",
			})
			c.Abort()
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID format",
			})
			c.Abort()
			return
		}

		if userID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User ID must be positive",
			})
			c.Abort()
			return
		}

		// Guardar el userID en el contexto para uso posterior
		c.Set(UserIDKey, userID)
		c.Next()
	}
}

// AuthWithUserValidationMiddleware extrae el ID de usuario y valida que existe en la BD
func AuthWithUserValidationMiddleware(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetHeader(UserIDHeader)

		if userIDStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "X-User-ID header is required",
			})
			c.Abort()
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID format",
			})
			c.Abort()
			return
		}

		if userID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User ID must be positive",
			})
			c.Abort()
			return
		}

		// Validar que el usuario existe en la base de datos
		_, err = userRepo.GetByID(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			c.Abort()
			return
		}

		// Guardar el userID en el contexto para uso posterior
		c.Set(UserIDKey, userID)
		c.Next()
	}
}

// GetUserID obtiene el ID de usuario del contexto
func GetUserID(c *gin.Context) int64 {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0
	}
	return userID.(int64)
}
