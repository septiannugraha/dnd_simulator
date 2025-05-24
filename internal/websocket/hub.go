package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"dnd-simulator/internal/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, implement proper origin checking
		return true
	},
}

// Client represents a WebSocket client connection
type Client struct {
	ID       string
	Conn     *websocket.Conn
	Hub      *Hub
	Send     chan []byte
	UserID   primitive.ObjectID
	Username string
	SessionID primitive.ObjectID
	CharacterID primitive.ObjectID
}

// Hub maintains active clients and broadcasts messages
type Hub struct {
	// Registered clients grouped by session
	Sessions map[primitive.ObjectID]map[*Client]bool
	
	// Register requests from clients
	Register chan *Client
	
	// Unregister requests from clients
	Unregister chan *Client
	
	// Broadcast to all clients in a session
	Broadcast chan *SessionMessage
	
	// Mutex for thread safety
	mu sync.RWMutex
}

// SessionMessage represents a message to broadcast to a session
type SessionMessage struct {
	SessionID primitive.ObjectID
	Message   models.WSMessage
	Exclude   *Client // Optional: exclude this client from broadcast
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		Sessions:   make(map[primitive.ObjectID]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *SessionMessage),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)
			
		case client := <-h.Unregister:
			h.unregisterClient(client)
			
		case sessionMsg := <-h.Broadcast:
			h.broadcastToSession(sessionMsg)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if h.Sessions[client.SessionID] == nil {
		h.Sessions[client.SessionID] = make(map[*Client]bool)
	}
	
	h.Sessions[client.SessionID][client] = true
	log.Printf("Client %s joined session %s", client.Username, client.SessionID.Hex())
	
	// Notify other clients in the session
	joinMessage := models.WSMessage{
		Type:      models.MessageTypePlayerJoined,
		Timestamp: time.Now(),
		UserID:    client.UserID,
		Username:  client.Username,
		SessionID: client.SessionID,
		Data: map[string]interface{}{
			"user_id":      client.UserID,
			"username":     client.Username,
			"character_id": client.CharacterID,
		},
	}
	
	h.broadcastToSessionExclude(client.SessionID, joinMessage, client)
	
	// Send success message to the joining client
	successMessage := models.WSMessage{
		Type:      models.MessageTypeSuccess,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message":    "Successfully joined session",
			"session_id": client.SessionID,
		},
	}
	h.sendToClient(client, successMessage)
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if clients, ok := h.Sessions[client.SessionID]; ok {
		if _, ok := clients[client]; ok {
			delete(clients, client)
			close(client.Send)
			
			// Clean up empty sessions
			if len(clients) == 0 {
				delete(h.Sessions, client.SessionID)
			}
			
			log.Printf("Client %s left session %s", client.Username, client.SessionID.Hex())
			
			// Notify other clients in the session
			leaveMessage := models.WSMessage{
				Type:      models.MessageTypePlayerLeft,
				Timestamp: time.Now(),
				UserID:    client.UserID,
				Username:  client.Username,
				SessionID: client.SessionID,
				Data: map[string]interface{}{
					"user_id":  client.UserID,
					"username": client.Username,
				},
			}
			
			h.broadcastToSessionExclude(client.SessionID, leaveMessage, nil)
		}
	}
}

func (h *Hub) broadcastToSession(sessionMsg *SessionMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	if clients, ok := h.Sessions[sessionMsg.SessionID]; ok {
		messageBytes, err := json.Marshal(sessionMsg.Message)
		if err != nil {
			log.Printf("Error marshaling message: %v", err)
			return
		}
		
		for client := range clients {
			if sessionMsg.Exclude != nil && client == sessionMsg.Exclude {
				continue
			}
			
			select {
			case client.Send <- messageBytes:
			default:
				// Client's send channel is full, remove them
				delete(clients, client)
				close(client.Send)
			}
		}
	}
}

func (h *Hub) broadcastToSessionExclude(sessionID primitive.ObjectID, message models.WSMessage, exclude *Client) {
	sessionMsg := &SessionMessage{
		SessionID: sessionID,
		Message:   message,
		Exclude:   exclude,
	}
	
	select {
	case h.Broadcast <- sessionMsg:
	default:
		log.Println("Broadcast channel full, dropping message")
	}
}

func (h *Hub) sendToClient(client *Client, message models.WSMessage) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}
	
	select {
	case client.Send <- messageBytes:
	default:
		log.Printf("Client %s send channel full", client.Username)
	}
}

// BroadcastToSession sends a message to all clients in a session
func (h *Hub) BroadcastToSession(sessionID primitive.ObjectID, message models.WSMessage) {
	sessionMsg := &SessionMessage{
		SessionID: sessionID,
		Message:   message,
	}
	
	select {
	case h.Broadcast <- sessionMsg:
	default:
		log.Println("Broadcast channel full, dropping message")
	}
}

// GetSessionClients returns the number of clients in a session
func (h *Hub) GetSessionClients(sessionID primitive.ObjectID) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	if clients, ok := h.Sessions[sessionID]; ok {
		return len(clients)
	}
	return 0
}

// HandleWebSocket handles the WebSocket upgrade and client management
func (h *Hub) HandleWebSocket(c *gin.Context) {
	// Get user info from auth middleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found"})
		return
	}
	
	// Get session ID from query parameter
	sessionIDStr := c.Query("session_id")
	if sessionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID required"})
		return
	}
	
	sessionID, err := primitive.ObjectIDFromHex(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}
	
	// Optional character ID
	characterIDStr := c.Query("character_id")
	var characterID primitive.ObjectID
	if characterIDStr != "" {
		characterID, err = primitive.ObjectIDFromHex(characterIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}
	}
	
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	
	// Create client
	client := &Client{
		ID:          primitive.NewObjectID().Hex(),
		Conn:        conn,
		Hub:         h,
		Send:        make(chan []byte, 256),
		UserID:      userID.(primitive.ObjectID),
		Username:    username.(string),
		SessionID:   sessionID,
		CharacterID: characterID,
	}
	
	// Register client
	h.Register <- client
	
	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// Client methods
func (c *Client) readPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	
	// Set read deadline and pong handler
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	
	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		
		// Parse message
		var message models.WSMessage
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}
		
		// Set message metadata
		message.Timestamp = time.Now()
		message.UserID = c.UserID
		message.Username = c.Username
		message.SessionID = c.SessionID
		
		// Handle the message based on type
		c.handleMessage(message)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
			
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(message models.WSMessage) {
	// Route message based on type
	switch message.Type {
	case models.MessageTypeChat, models.MessageTypeChatIC, models.MessageTypeChatOOC:
		c.handleChatMessage(message)
	case models.MessageTypeDiceRoll:
		c.handleDiceRoll(message)
	case models.MessageTypeCharacterUpdate:
		c.handleCharacterUpdate(message)
	default:
		log.Printf("Unknown message type: %s", message.Type)
	}
}

func (c *Client) handleChatMessage(message models.WSMessage) {
	// Broadcast chat message to all clients in the session
	c.Hub.BroadcastToSession(c.SessionID, message)
}

func (c *Client) handleDiceRoll(message models.WSMessage) {
	// Process dice roll and broadcast result
	c.Hub.BroadcastToSession(c.SessionID, message)
}

func (c *Client) handleCharacterUpdate(message models.WSMessage) {
	// Broadcast character update to all clients in the session
	c.Hub.BroadcastToSession(c.SessionID, message)
}