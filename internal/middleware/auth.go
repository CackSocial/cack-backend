package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/CackSocial/cack-backend/pkg/auth"
	"github.com/CackSocial/cack-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// extractToken tries to read the JWT from the cookie first, then falls back to
// the Authorization header for backward compatibility.
func extractToken(c *gin.Context) string {
	if token, err := c.Cookie("sc-token"); err == nil && token != "" {
		return token
	}
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			response.Error(c, http.StatusUnauthorized, "authorization required")
			c.Abort()
			return
		}

		token, err := auth.ValidateToken(tokenString, jwtSecret)
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
		tokenString := extractToken(c)
		if tokenString == "" {
			c.Set("userID", "")
			c.Next()
			return
		}

		token, err := auth.ValidateToken(tokenString, jwtSecret)
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
