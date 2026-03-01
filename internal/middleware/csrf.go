package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GenerateCSRFToken creates a cryptographically random CSRF token.
func GenerateCSRFToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic("csrf: failed to generate random token: " + err.Error())
	}
	return hex.EncodeToString(b)
}

// CSRFMiddleware validates the CSRF token on state-changing requests.
// The token in the X-CSRF-Token header must match the sc-csrf cookie.
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		csrfCookie, err := c.Cookie("sc-csrf")
		csrfHeader := c.GetHeader("X-CSRF-Token")

		if err != nil || csrfCookie == "" || csrfCookie != csrfHeader {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CSRF token mismatch"})
			return
		}

		c.Next()
	}
}
