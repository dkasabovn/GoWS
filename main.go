package main

import (
	"flag"
	"fmt"
	"log"
	"main/config"
	"main/ws"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http server address")

func main() {
	flag.Parse()

	config.CreateRedisClient()

	fmt.Println("big pog")
	hub := ws.NewHub()

	http.HandleFunc("/ws", hub.ServeRoom)
	http.HandleFunc("/create", hub.StartGame)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
