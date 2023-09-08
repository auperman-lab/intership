package chat

import (
	"github.com/gorilla/websocket"
	"strconv"
	"strings"
	"sync"
)

type Client struct {
	Conn *websocket.Conn
	ID   uint
	Hub  *Hub
	Send chan []byte
}

type Message struct {
	SenderID   uint   `json:"sender_id"`
	ReceiverID uint   `json:"receiver_id"`
	Text       []byte `json:"text"`
}

func (c *Client) Read() {
	var unregisterOnce sync.Once // Use a sync.Once to ensure unregistering only happens once.

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		parts := strings.SplitN(string(message), " ", 3)
		if len(parts) >= 3 && parts[0] == "/pm" {
			recipientID, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				continue
			}

			// Send the private message to the Hub.
			c.Hub.Send <- Message{
				SenderID:   c.ID,
				ReceiverID: uint(recipientID),
				Text:       []byte(parts[2]),
			}
		} else {
			// This is a public message. Broadcast it to all clients.
			c.Hub.Send <- Message{
				SenderID:   c.ID,
				ReceiverID: 0, // 0 indicates a public message.
				Text:       message,
			}
		}
	}

	// Use sync.Once to ensure unregistering only happens once.
	unregisterOnce.Do(func() {
		c.Hub.Unregister <- c
		close(c.Hub.Send)
	})

	// Close the WebSocket connection.
	c.Conn.Close()
}

func (c *Client) Write() {
	defer func() {
		c.Conn.Close()
	}()

	for message := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}
