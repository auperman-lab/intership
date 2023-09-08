package chat

type Hub struct {
	Clients    map[uint]*Client
	Send       chan Message
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Send:       make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[uint]*Client),
	}
}

func Run(h *Hub) {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.ID] = client
		case client := <-h.Unregister:
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
			}
		case message := <-h.Send:
			if client, ok := h.Clients[message.ReceiverID]; ok {
				select {
				case client.Send <- message.Text:
				default:
					close(client.Send)
				}
			}
		}
	}
}
