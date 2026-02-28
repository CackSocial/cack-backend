package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/CackSocial/cack-backend/pkg/auth"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

type WSHandler struct {
	hub       *Hub
	jwtSecret string
}

func NewWSHandler(hub *Hub, jwtSecret string) *WSHandler {
	return &WSHandler{hub: hub, jwtSecret: jwtSecret}
}

func (h *WSHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/api/v1/ws", h.HandleWebSocket)
}

func (h *WSHandler) HandleWebSocket(c *gin.Context) {
	// Authenticate via query parameter token
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	token, err := auth.ValidateToken(tokenString, h.jwtSecret)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	userID, err := auth.ExtractUserID(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	h.hub.register <- client

	go h.hub.WritePump(client)
	go h.hub.ReadPump(client)
}
