package ws

type Server struct {
	broadcast  chan *Message
	register   chan *Client
	internal   chan *Message
	unregister chan *Client
	clients    map[*Client]bool
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
		internal:   make(chan *Message),
	}
}
