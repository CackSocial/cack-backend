package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/CackSocial/cack-backend/internal/dto"
	messageUC "github.com/CackSocial/cack-backend/internal/usecase/message"
)

type Client struct {
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
}

type Hub struct {
	clients    map[string]*Client // userID -> Client
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	messageUC  *messageUC.MessageUseCase
}

func NewHub(muc *messageUC.MessageUseCase) *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		messageUC:  muc,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			// Close existing connection if user reconnects
			if existing, ok := h.clients[client.UserID]; ok {
				close(existing.Send)
				if err := existing.Conn.Close(); err != nil {
					log.Printf("failed to close existing connection: %v", err)
				}
			}
			h.clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("Client connected: %s", client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if existing, ok := h.clients[client.UserID]; ok && existing == client {
				close(client.Send)
				delete(h.clients, client.UserID)
			}
			h.mu.Unlock()
			log.Printf("Client disconnected: %s", client.UserID)
		}
	}
}

// SendToUser sends a message to a specific user if they're online.
func (h *Hub) SendToUser(userID string, data []byte) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()
	if ok {
		select {
		case client.Send <- data:
		default:
			// Client's send buffer is full, remove them
			h.mu.Lock()
			close(client.Send)
			delete(h.clients, userID)
			h.mu.Unlock()
		}
	}
}

// ReadPump reads messages from the WebSocket connection.
func (h *Hub) ReadPump(client *Client) {
	defer func() {
		h.unregister <- client
		if err := client.Conn.Close(); err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}()

	for {
		_, msgBytes, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var wsMsg dto.WSMessage
		if err := json.Unmarshal(msgBytes, &wsMsg); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		switch wsMsg.Type {
		case "message":
			h.handleMessage(client, &wsMsg)
		default:
			log.Printf("Unknown message type: %s", wsMsg.Type)
		}
	}
}

// WritePump writes messages to the WebSocket connection.
func (h *Hub) WritePump(client *Client) {
	defer func() {
		if err := client.Conn.Close(); err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}()

	for message := range client.Send {
		if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("WebSocket write error: %v", err)
			break
		}
	}
}

func (h *Hub) handleMessage(client *Client, wsMsg *dto.WSMessage) {
	if wsMsg.ReceiverID == "" || (wsMsg.Content == "" && wsMsg.ImageURL == "") {
		return
	}

	// Persist message to DB
	msg, err := h.messageUC.SendFromWS(client.UserID, wsMsg.ReceiverID, wsMsg.Content, wsMsg.ImageURL)
	if err != nil {
		log.Printf("Failed to save message: %v", err)
		return
	}

	// Build response
	responseMsg := map[string]interface{}{
		"type":        "message",
		"id":          msg.ID,
		"sender_id":   msg.SenderID,
		"receiver_id": msg.ReceiverID,
		"content":     msg.Content,
		"image_url":   msg.ImageURL,
		"created_at":  msg.CreatedAt,
	}
	data, _ := json.Marshal(responseMsg)

	// Send to receiver if online
	h.SendToUser(wsMsg.ReceiverID, data)
	// Also echo back to sender
	h.SendToUser(client.UserID, data)
}
