package ws

const PubSubGeneralChannel = "general"

type Server struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	room       *Room
}

func NewWS() *Server {
	ws := &Server{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	return ws
}
