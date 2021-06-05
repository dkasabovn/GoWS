package ws

import (
	"context"
	"log"
	"main/config"

	"github.com/google/uuid"
)

type Room struct {
	id         uuid.UUID
	name       string
	capacity   int
	private    bool
	server     *Server
	ctx        context.Context
	cancelFunc func()
}

func NewRoom(private bool) *Room {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &Room{
		id:         uuid.New(),
		server:     NewServer(),
		private:    private,
		capacity:   1 << 3,
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}
}

func (r *Room) ID() string {
	return r.id.String()
}

func (r *Room) Active() int {
	return len(r.server.clients)
}

func (r *Room) BroadcastMessage(message *Message) {
	r.server.broadcast <- message
}

func (r *Room) publishRoomMessage(message []byte) {
	err := config.Redis.Publish(config.CTX, r.id.String(), message).Err()

	if err != nil {
		log.Println(err)
	}
}

func (r *Room) unregisterClient(client *Client) {
	if _, ok := r.server.clients[client]; ok {
		delete(r.server.clients, client)
		r.notifyClientLeft(client)
	}
}

func (r *Room) registerClient(client *Client) {
	r.server.clients[client] = true
	r.notifyClientJoined(client)
	r.notifyClientOfParticipants(client)
}

func (r *Room) notifyClientOfParticipants(client *Client) {
	participants := make([]string, 0)
	for k := range r.server.clients {
		participants = append(participants, k.Name)
	}
	message := &Message{
		Action: BootstrapData,
		Data: map[string]interface{}{
			"users": participants,
		},
	}
	client.send <- message.encode()
}

func (r *Room) subscribeToRoomMessages() {
	pubsub := config.Redis.Subscribe(config.CTX, r.id.String())

	ch := pubsub.Channel()

	for {

		select {
		case msg := <-ch:
			r.broadcastToClients([]byte(msg.Payload))
		case <-r.ctx.Done():
			pubsub.Close()
			return
		}
	}
}

func (r *Room) broadcastToClients(message []byte) {
	for client := range r.server.clients {
		client.send <- message
	}
}

func (r *Room) notifyClientJoined(client *Client) {
	message := &Message{
		Action: UserJoinedAction,
		Data: map[string]interface{}{
			"user": client.Name,
		},
	}

	r.publishRoomMessage(message.encode())
}

func (r *Room) notifyClientLeft(client *Client) {
	message := &Message{
		Action: UserLeftAction,
		Data: map[string]interface{}{
			"user": client.Name,
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
		case msg := <-r.server.broadcast:
			log.Println("Room; Broadcasting")
			r.publishRoomMessage(msg.encode())
		case <-r.ctx.Done():
			log.Printf("Room %s cancelled\n", r.id.String())
			for k := range r.server.clients {
				k.Disconnect()
			}
			close(r.server.register)
			close(r.server.unregister)
			close(r.server.broadcast)
			close(r.server.internal)
			return
		}

	}
}

func (r *Room) Terminate() {
	r.cancelFunc()
}
