package main

import (
	"flag"
	"fmt"
	"log"
	"main/config"
	"main/game"
	"main/ws"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http server address")

func main() {
	flag.Parse()

	config.CreateRedisClient()

	fmt.Println("big pog")
	hub := ws.NewHub()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		room, ok := hub.ServeRoom(w, r)
		if !ok {
			log.Println("Room connection failed")
		}
		gme := game.NewGameManager(room)
		go gme.Run()
	})
	http.HandleFunc("/create", hub.StartGame)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
