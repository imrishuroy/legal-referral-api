package chat

import (
	"context"
	"fmt"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/rs/zerolog/log"
)

// Hub is a struct that holds all the clients and the messages that are sent to them
type Hub struct {
	// Registered clients.
	clients map[string]map[*Client]bool
	//Unregistered clients.
	unregister chan *Client
	// Register requests from the clients.
	register chan *Client
	// Inbound messages from the clients.
	broadcast chan Message
	// storage for messages
	store db.Store
}

// Message struct to hold message data
type Message struct {
	MessageID       int32  `json:"message_id"`
	ParentMessageID int32  `json:"parent_message_id"` // fix:
	SenderID        string `json:"sender_id"`
	RecipientID     string `json:"recipient_id"`
	Message         string `json:"message"`
	HasAttachment   bool   `json:"has_attachment"`
	AttachmentID    int32  `json:"attachment_url"`
	IsRead          bool   `json:"is_read"`
	RoomID          string `json:"room_id"`
}

func NewHub(store db.Store) *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		unregister: make(chan *Client),
		register:   make(chan *Client),
		broadcast:  make(chan Message),
		store:      store,
	}
}

// Run Core function to run the hub
func (h *Hub) Run() {
	for {
		select {
		// Register a client.
		case client := <-h.register:
			h.RegisterNewClient(client)
			// Unregister a client.
		case client := <-h.unregister:
			h.RemoveClient(client)
			// Broadcast a message to all clients.
		case message := <-h.broadcast:
			//Check if the message is a type of "message"
			h.HandleMessage(message)

		}
	}
}

// RegisterNewClient function check if room exists and if not create it and add client to it
func (h *Hub) RegisterNewClient(client *Client) {
	connections := h.clients[client.RoomID]
	if connections == nil {
		connections = make(map[*Client]bool)
		h.clients[client.RoomID] = connections
	}
	h.clients[client.RoomID][client] = true
	// here we can add the presence of the client to the room

	fmt.Println("Size of clients: ", len(h.clients[client.RoomID]))
}

// RemoveClient function to remove client from room
func (h *Hub) RemoveClient(client *Client) {
	if _, ok := h.clients[client.RoomID]; ok {
		delete(h.clients[client.RoomID], client)
		close(client.send)
		fmt.Println("Removed client")
	}
}

// HandleMessage function to handle message based on type of message
func (h *Hub) HandleMessage(message Message) {

	// here we can add the message to the database

	arg := db.CreateMessageParams{
		ParentMessageID: &message.ParentMessageID,
		SenderID:        message.SenderID,
		RecipientID:     message.RecipientID,
		Message:         message.Message,
		HasAttachment:   message.HasAttachment,
		AttachmentID:    &message.AttachmentID,
		IsRead:          message.IsRead,
		RoomID:          message.RoomID,
	}
	_, err := h.store.CreateMessage(context.Background(), arg)
	if err != nil {
		log.Error().Err(err).Msg("Error creating message")
		return
	}

	clients := h.clients[message.RoomID]
	for client := range clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients[message.RoomID], client)
		}
	}
}
