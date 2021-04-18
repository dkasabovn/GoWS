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
	config.CreateFirestoreClient()
	fmt.Println("big pog")
	ws.NewHub()

	http.HandleFunc("/ws", ws.MainHub.ServeRoom)
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		room := ws.MainHub.StartGame(w, r)
		gme := game.NewGameManager(room)
		go gme.Run()
	})

	log.Fatal(http.ListenAndServe(*addr, nil))

}
