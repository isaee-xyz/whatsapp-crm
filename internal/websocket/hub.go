package websocket

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/zerodha/logf"
)

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	// clients maps organization ID -> user ID -> client
	clients map[uuid.UUID]map[uuid.UUID]*Client

	// broadcast channel for messages
	broadcast chan BroadcastMessage

	// register channel for new clients
	register chan *Client

	// unregister channel for disconnecting clients
	unregister chan *Client

	// mutex for thread-safe access to clients map
	mu sync.RWMutex

	// logger
	log logf.Logger
}

// NewHub creates a new Hub instance
func NewHub(log logf.Logger) *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]map[uuid.UUID]*Client),
		broadcast:  make(chan BroadcastMessage, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		log:        log,
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient adds a client to the hub
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	orgClients, ok := h.clients[client.organizationID]
	if !ok {
		orgClients = make(map[uuid.UUID]*Client)
		h.clients[client.organizationID] = orgClients
	}

	// Close existing connection for same user if any
	if existing, ok := orgClients[client.userID]; ok {
		close(existing.send)
	}

	orgClients[client.userID] = client
	h.log.Info("WebSocket client registered",
		"user_id", client.userID,
		"org_id", client.organizationID,
		"total_clients", h.countClients())
}

// unregisterClient removes a client from the hub
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if orgClients, ok := h.clients[client.organizationID]; ok {
		if existing, ok := orgClients[client.userID]; ok && existing == client {
			delete(orgClients, client.userID)
			close(client.send)

			// Clean up empty org map
			if len(orgClients) == 0 {
				delete(h.clients, client.organizationID)
			}
		}
	}

	h.log.Info("WebSocket client unregistered",
		"user_id", client.userID,
		"org_id", client.organizationID,
		"total_clients", h.countClients())
}

// broadcastMessage sends a message to all relevant clients
func (h *Hub) broadcastMessage(msg BroadcastMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	orgClients, ok := h.clients[msg.OrgID]
	if !ok {
		return
	}

	data, err := json.Marshal(msg.Message)
	if err != nil {
		h.log.Error("Failed to marshal broadcast message", "error", err)
		return
	}

	for _, client := range orgClients {
		// If ContactID is specified, only send to clients viewing that contact
		if msg.ContactID != uuid.Nil && client.currentContact != nil && *client.currentContact != msg.ContactID {
			continue
		}

		select {
		case client.send <- data:
		default:
			// Client buffer full, skip
			h.log.Warn("Client send buffer full, skipping",
				"user_id", client.userID,
				"org_id", client.organizationID)
		}
	}
}

// Broadcast sends a message to the broadcast channel
func (h *Hub) Broadcast(msg BroadcastMessage) {
	select {
	case h.broadcast <- msg:
	default:
		h.log.Warn("Broadcast channel full, dropping message")
	}
}

// BroadcastToOrg sends a message to all clients in an organization
func (h *Hub) BroadcastToOrg(orgID uuid.UUID, msg WSMessage) {
	h.Broadcast(BroadcastMessage{
		OrgID:   orgID,
		Message: msg,
	})
}

// BroadcastToContact sends a message to clients viewing a specific contact
func (h *Hub) BroadcastToContact(orgID, contactID uuid.UUID, msg WSMessage) {
	h.Broadcast(BroadcastMessage{
		OrgID:     orgID,
		ContactID: contactID,
		Message:   msg,
	})
}

// countClients returns the total number of connected clients
func (h *Hub) countClients() int {
	count := 0
	for _, orgClients := range h.clients {
		count += len(orgClients)
	}
	return count
}

// GetClientCount returns the number of connected clients (thread-safe)
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.countClients()
}

// Register adds a client to the hub via the register channel
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client from the hub via the unregister channel
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}
