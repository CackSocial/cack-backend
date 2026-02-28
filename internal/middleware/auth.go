package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/CackSocial/cack-backend/pkg/auth"
	"github.com/CackSocial/cack-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "invalid authorization header format")
			c.Abort()
			return
		}

		token, err := auth.ValidateToken(parts[1], jwtSecret)
		if err != nil || !token.Valid {
			slog.Warn("auth: invalid or expired token",
				"error", err,
				"ip", c.ClientIP(),
				"path", c.Request.URL.Path,
			)
			response.Error(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		userID, err := auth.ExtractUserID(token)
		if err != nil {
			slog.Warn("auth: invalid token claims",
				"error", err,
				"ip", c.ClientIP(),
				"path", c.Request.URL.Path,
			)
			response.Error(c, http.StatusUnauthorized, "invalid token claims")
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

// OptionalAuth is like AuthMiddleware but doesn't abort if no token provided.
// Sets userID in context if valid token exists, otherwise sets empty string.
func OptionalAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Set("userID", "")
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Set("userID", "")
			c.Next()
			return
		}

		token, err := auth.ValidateToken(parts[1], jwtSecret)
		if err != nil || !token.Valid {
			c.Set("userID", "")
			c.Next()
			return
		}

		userID, err := auth.ExtractUserID(token)
		if err != nil {
			c.Set("userID", "")
			c.Next()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
