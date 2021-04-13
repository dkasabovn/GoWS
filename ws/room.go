package ws

import (
	"fmt"
	"log"
	"main/config"

	"github.com/google/uuid"
)

type Room struct {
	id        uuid.UUID
	name      string
	capacity  int
	private   bool
	server    *Server
	broadcast chan *Message
	Commands  chan *Message
}

func NewRoom(private bool) *Room {
	return &Room{
		id:        uuid.New(),
		server:    NewWS(),
		broadcast: make(chan *Message, 1<<3),
		private:   private,
		capacity:  1 << 3,
	}
}

func (r *Room) Active() int {
	return len(r.server.clients)
}

func (r *Room) publishRoomMessage(message []byte) {
	err := config.Redis.Publish(config.CTX, r.name, message).Err()

	if err != nil {
		log.Println(err)
	}
}

func (r *Room) unregisterClient(client *Client) {
	if _, ok := r.server.clients[client]; ok {
		delete(r.server.clients, client)
	}
}

func (r *Room) registerClient(client *Client) {
	r.server.clients[client] = true
}

func (r *Room) subscribeToRoomMessages() {
	pubsub := config.Redis.Subscribe(config.CTX, r.name)

	ch := pubsub.Channel()

	for msg := range ch {
		r.broadcastToClients([]byte(msg.Payload))
	}
}

func (r *Room) broadcastToClients(message []byte) {
	for client := range r.server.clients {
		client.send <- message
	}
}

func (r *Room) notifyClientJoined(client *Client) {
	message := &Message{
		Action: SendMessageAction,
		Data: map[string]interface{}{
			"message": fmt.Sprintf("%s joined", client.Name),
		},
	}

	r.publishRoomMessage(message.encode())
}

func (r *Room) Run() {
	go r.subscribeToRoomMessages()

	for {
		select {
		case client := <-r.server.register:
			log.Println("Room; Client Registered")
			r.registerClient(client)
		case client := <-r.server.unregister:
			log.Println("Room; Client Unregistered")
			r.unregisterClient(client)
		case msg := <-r.broadcast:
			log.Println("Room; Broadcasting")
			r.publishRoomMessage(msg.encode())
		}
	}
}
