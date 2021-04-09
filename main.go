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

	mainRoom := ws.NewRoom("main", false)
	go mainRoom.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.Join(mainRoom, w, r)
	})

	log.Fatal(http.ListenAndServe(*addr, nil))
}
