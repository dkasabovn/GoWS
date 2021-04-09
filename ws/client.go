package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

type Client struct {
	conn  *websocket.Conn
	send  chan []byte
	ID    uuid.UUID
	Name  string
	room  *Room
	daHub *Hub
}

func NewClient(conn *websocket.Conn, server *Room, name string) *Client {
	// TODO fetch uuid from supabase given auth
	return &Client{
		ID:   uuid.New(),
		Name: name,
		conn: conn,
		room: server,
		send: make(chan []byte, 256),
	}
}

func Join(room *Room, w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	name, ok := params["name"]

	if !ok {
		log.Println("Client is a monkey")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, room, name[0])

	go client.writePump()
	go client.readPump()

	room.server.register <- client
}

func (c *Client) joinRoom(roomName string) *Room {
	room := c.daHub.GetRoom(roomName)
	if room == nil {
		return nil
	}
	if c.room != room {
		room.server.register <- c
		c.notifyRoomJoined(room)
	}
	return room
}

func (c *Client) notifyRoomJoined(room *Room) {
	message := Message{
		Action: RoomJoinedAction,
		Target: room,
		Sender: c,
	}
	c.send <- message.encode()
}

func (c *Client) handleNewMessage(jsonMessage []byte) {

	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}

	message.Sender = c

	switch message.Action {
	case SendMessageAction:
		c.room.broadcast <- &message

	case JoinRoomAction:
		c.handleJoinRoomMessage(message)

	case LeaveRoomAction:
		c.handleLeaveRoomMessage(message)
	}

}

func (c *Client) handleJoinRoomMessage(message Message) {
	roomName := message.Message

	c.joinRoom(roomName)
}

func (c *Client) handleLeaveRoomMessage(message Message) {
	room := c.daHub.GetRoom(message.Message)
	if room == nil {
		return
	}
	room.server.unregister <- c
}

func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		client.handleNewMessage(jsonMessage)
	}

}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) disconnect() {
	c.room.server.unregister <- c
	close(c.send)
	c.conn.Close()
}
