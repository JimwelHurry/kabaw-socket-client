package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	mathrand "math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Message represents a chat message
type Message struct {
	Type      string `json:"type"`
	Username  string `json:"username"`
	UserID    string `json:"user_id,omitempty"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
	Channel   string `json:"channel,omitempty"`
}

// Client represents a connected WebSocket client
type Client struct {
	conn     *websocket.Conn
	send     chan Message
	hub      *Hub
	username string
	userID   string
	channel  string
}

// Hub maintains the set of active clients and broadcasts messages to clients
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin (for development)
		return true
	},
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Simulated messages for the general channel
var simulatedMessages = []string{
	"Welcome to the chat! 👋",
	"How's everyone doing today?",
	"This is a simulated message to keep the chat active",
	"Feel free to join the conversation!",
	"Testing the WebSocket connection...",
	"Anyone else working on Go projects?",
	"The weather is nice today ☀️",
	"Don't forget to stay hydrated! 💧",
	"What's your favorite programming language?",
	"This chat supports multiple channels",
	"WebSocket connections are pretty cool!",
	"Hope you're having a great day! 😊",
	"Remember to take breaks while coding",
	"Coffee or tea? ☕",
	"The server is running smoothly",
	"Cross-origin requests work perfectly here",
	"Real-time messaging is awesome!",
	"Thanks for testing the chat system",
}

// SimulatedUser represents a test user with UUID
type SimulatedUser struct {
	Username string
	UserID   string
}

var simulatedUsers = []SimulatedUser{
	{
		Username: "ChatBot",
		UserID:   "550e8400-e29b-41d4-a716-446655440001", // Fixed UUID for ChatBot
	},
	{
		Username: "Developer",
		UserID:   "550e8400-e29b-41d4-a716-446655440002", // Fixed UUID for Developer
	},
	{
		Username: "SystemHelper",
		UserID:   "550e8400-e29b-41d4-a716-446655440003", // Fixed UUID for SystemHelper
	},
}

// StartMessageSimulation starts generating simulated messages for the general channel
func (h *Hub) StartMessageSimulation() {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Only send messages if there are clients connected to general channel
				h.mutex.RLock()
				hasGeneralClients := false
				for client := range h.clients {
					if client.channel == "general" {
						hasGeneralClients = true
						break
					}
				}
				h.mutex.RUnlock()

				if hasGeneralClients {
					// Generate a random simulated message
					selectedUser := simulatedUsers[mathrand.Intn(len(simulatedUsers))]
					message := Message{
						Type:      "message",
						Username:  selectedUser.Username,
						UserID:    selectedUser.UserID,
						Content:   simulatedMessages[mathrand.Intn(len(simulatedMessages))],
						Timestamp: getCurrentTimestamp(),
						Channel:   "general",
					}

					// Log the simulated message
					log.Printf("[SIMULATED] Channel: %s | User: %s (ID: %s) | Content: %s", message.Channel, message.Username, message.UserID, message.Content)

					h.broadcast <- message
				}

				// Reset ticker to 10 seconds
				ticker.Reset(10 * time.Second)
			}
		}
	}()
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()

			log.Printf("[CONNECT] User: %s | Channel: %s | Total clients: %d", client.username, client.channel, len(h.clients))

			// Send welcome message with user ID
			welcomeMsg := Message{
				Type:      "user_connected",
				Content:   "Welcome to the chat!",
				Username:  "System",
				UserID:    client.userID,
				Timestamp: getCurrentTimestamp(),
				Channel:   client.channel,
			}
			select {
			case client.send <- welcomeMsg:
			default:
				close(client.send)
				delete(h.clients, client)
			}

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("[DISCONNECT] User: %s | Channel: %s | Total clients: %d", client.username, client.channel, len(h.clients))
			}
			h.mutex.Unlock()

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				// Only send message to clients in the same channel (or all channels if no channel specified)
				if message.Channel == "" || client.channel == message.Channel {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// GetConnectedClients returns the number of connected clients
func (h *Hub) GetConnectedClients() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Set message metadata
		msg.Username = c.username
		msg.UserID = c.userID
		msg.Timestamp = getCurrentTimestamp()
		msg.Channel = c.channel

		// Log the incoming message
		log.Printf("[MESSAGE] Channel: %s | User: %s | Content: %s", msg.Channel, msg.Username, msg.Content)

		// Broadcast the message
		c.hub.broadcast <- msg
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				log.Printf("error: %v", err)
				return
			}
		}
	}
}

// WebSocket handler
func handleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// Get username and channel from query parameters
	username := r.URL.Query().Get("username")
	channel := r.URL.Query().Get("channel")

	if username == "" {
		username = "Anonymous"
	}
	if channel == "" {
		channel = "general"
	}

	// Generate a unique user ID for this connection
	userID := generateUserID()

	client := &Client{
		conn:     conn,
		send:     make(chan Message, 256),
		hub:      hub,
		username: username,
		userID:   userID,
		channel:  channel,
	}

	client.hub.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// Health check endpoint
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":  "ok",
		"service": "kabaw-discord-server",
	}
	json.NewEncoder(w).Encode(response)
}

// Stats endpoint
func handleStats(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"connected_clients": hub.GetConnectedClients(),
		}
		json.NewEncoder(w).Encode(response)
	}
}

// Get current timestamp in ISO format
func getCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

// generateUserID creates a random user ID
func generateUserID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func main() {
	// Seed random number generator
	mathrand.Seed(time.Now().UnixNano())

	hub := NewHub()
	go hub.Run()

	// Start message simulation for general channel
	hub.StartMessageSimulation()

	// Routes
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(hub, w, r)
	})
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/stats", handleStats(hub))

	// Serve CSS file specifically
	http.HandleFunc("/static/styles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "styles.css")
	})

	// Serve test client
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "test_client.html")
	})

	// Default route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <title>Kabaw Discord Server</title>
</head>
<body>
    <h1>Kabaw Discord WebSocket Server</h1>
    <p>WebSocket endpoint: <code>ws://localhost:8080/ws</code></p>
    <p>Health check: <a href="/health">/health</a></p>
    <p>Stats: <a href="/stats">/stats</a></p>
    <p>Test Client: <a href="/test">Launch Test Client</a></p>
    <h2>Usage:</h2>
    <p>Connect to WebSocket with query parameters:</p>
    <ul>
        <li><code>username</code> - Your display name</li>
        <li><code>channel</code> - Channel name (default: "general")</li>
    </ul>
    <p>Example: <code>ws://localhost:8080/ws?username=John&channel=general</code></p>
</body>
</html>
		`))
	})

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
