package ws

import (
	"context"
	"fmt"
	"log"
	"main/config"

	"github.com/google/uuid"
)

type Room struct {
	id         uuid.UUID
	name       string
	capacity   int
	private    bool
	Broadcast  chan *Message
	Commands   chan *Message
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	ctx        context.Context
	cancelFunc func()
}

func NewRoom(private bool) *Room {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &Room{
		id:         uuid.New(),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 1<<3),
		Commands:   make(chan *Message),
		private:    private,
		capacity:   1 << 3,
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}
}

func (r *Room) Active() int {
	return len(r.clients)
}

func (r *Room) publishRoomMessage(message []byte) {
	err := config.Redis.Publish(config.CTX, r.id.String(), message).Err()

	if err != nil {
		log.Println(err)
	}
}

func (r *Room) unregisterClient(client *Client) {
	if _, ok := r.clients[client]; ok {
		delete(r.clients, client)
	}
}

func (r *Room) registerClient(client *Client) {
	r.clients[client] = true
}

func (r *Room) subscribeToRoomMessages() {
	pubsub := config.Redis.Subscribe(config.CTX, r.id.String())

	ch := pubsub.Channel()

	for msg := range ch {
		select {
		case <-r.ctx.Done():
			pubsub.Close()
			return
		default:
			r.broadcastToClients([]byte(msg.Payload))
		}
	}
}

func (r *Room) broadcastToClients(message []byte) {
	for client := range r.clients {
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
		case client := <-r.register:
			log.Println("Room; Client Registered")
			r.registerClient(client)
		case client := <-r.unregister:
			log.Println("Room; Client Unregistered")
			r.unregisterClient(client)
		case msg := <-r.Broadcast:
			log.Println("Room; Broadcasting")
			r.publishRoomMessage(msg.encode())
		case <-r.ctx.Done():
			log.Printf("Room %s cancelled\n", r.id.String())
			for k := range r.clients {
				k.Disconnect()
			}
			close(r.register)
			close(r.unregister)
			close(r.Broadcast)
			close(r.Commands)
			return
		}

	}
}

func (r *Room) Terminate() {
	r.cancelFunc()
}
