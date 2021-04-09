package ws

type Hub struct {
	server    *Room
	roomStore map[string]*Room
}

func (h *Hub) GetRoom(id string) *Room {
	return h.roomStore[id]
}
