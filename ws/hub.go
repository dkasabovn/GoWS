package ws

import (
	"log"
	"net/http"
	"time"
)

// TODO store all rooms internally -> route users to rooms
// ? disconnect users who are afk
// TODO find matches between users
// ? write tests
// ? Factor in ELO
// * user joins -> search for game -> if game is found register client to room via register channel; deregister client from lobby

var Main *Hub

type Hub struct {
	server    *Server
	algoTimer *time.Ticker
	rooms     map[string]*Room
}

func Init() {
	Main = &Hub{
		algoTimer: time.NewTicker(time.Second * 5),
		rooms:     make(map[string]*Room),
		server:    NewServer(),
	}
}

func (h *Hub) onNewClient(client *Client) {
	h.server.clients[client] = true
}

func (h *Hub) unregisterUser(client *Client) {
	if _, ok := h.server.clients[client]; ok {
		delete(h.server.clients, client)
	}
}

func (h *Hub) matchmakingAlgorithm() {
	var clientOne, clientTwo *Client
	for k := range h.server.clients {
		if clientOne == nil {
			clientOne = k
		} else {
			clientTwo = k
			h.createGame(clientOne, clientTwo)
			clientOne = nil
			clientTwo = nil
		}
		log.Println("looping")
	}
}

func (h *Hub) createGame(one, two *Client) {
	room := NewRoom(false)
	go room.Run()
	room.server.register <- one
	room.server.register <- two
	h.unregisterUser(one)
	h.unregisterUser(two)
	log.Println("Users connected to Game Room")
	gme := NewGameManager(room)
	go gme.Run()
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.server.register:
			log.Println("User In Q")
			h.onNewClient(client)
		case client := <-h.server.unregister:
			log.Println("Dequeued user")
			h.unregisterUser(client)
		case <-h.algoTimer.C:
			if len(h.server.clients) > 1 {
				h.matchmakingAlgorithm()
				log.Println("Matchmaking and finding users")
			}
			log.Printf("Not enough players %d\n", len(h.server.clients))
			break
		}
	}
}

func (h *Hub) ServeHub(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	name, ok := params["name"]

	if !ok {
		log.Fatalln("Put in a name doofenshmirtz")
	}

	conn, err := Upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatalln("Idiot")
	}

	client := NewClient(conn, h.server, name[0])
	client.Run()
	log.Println("Client connected; Pumps started")
}
