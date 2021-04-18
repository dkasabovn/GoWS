package ws

import (
	"encoding/json"
	"log"
	"net/http"
)

var MainHub *Hub

type Hub struct {
	registeredRooms map[string]*Room
}

func NewHub() {
	MainHub = &Hub{
		registeredRooms: make(map[string]*Room),
	}
}

func (h *Hub) GetRoom(room string) (*Room, bool) {
	// check if Room is registered in memory
	if r, ok := h.registeredRooms[room]; ok {
		return r, true
	}
	// ? If using multiple cluster system the above check will not always work depending on which server the game is being hosted on
	return nil, false
}

func (h *Hub) RemoveRoom(room string) {
	delete(h.registeredRooms, room)
}

func (h *Hub) CreateRoom() *Room {
	r := NewRoom(false)
	h.registeredRooms[r.ID()] = r
	go r.Run()
	return r
}

func (h *Hub) StartGame(w http.ResponseWriter, r *http.Request) *Room {
	// TODO once done with game manager logic like totally do this
	log.Println("Got create room request")
	room := h.CreateRoom()
	data := map[string]interface{}{
		"socket": room.ID(),
	}
	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println("Finished request")
	return room
}

func (h *Hub) ServeRoom(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	name, ok := params["name"]

	if !ok {
		log.Println("Client does not have a name apparently. Really rude of them tbh")
	}

	roomName, rok := params["room"]

	if !rok {
		log.Println("No room defined; Central Hub is not configured")
	}

	room, exists := h.GetRoom(roomName[0])
	if !exists {
		log.Println("Room defined but doesn't exist; Front end error?")
	}

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := NewClient(conn, room, name[0])

	client.Run()

	log.Println("Client Connected; Pumps Started")
}

func (h *Hub) TerminateRoom(id string) {
	room, ok := h.GetRoom(id)
	if !ok {
		return
	}
	room.Terminate()
}
