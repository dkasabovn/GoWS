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

type Claims2 struct {
	name    string
	email   string
	picture string
	sub     string
	iat     string
	exp     string
}

func main() {
	flag.Parse()

	config.CreateRedisClient()
	config.CreateFirestoreClient()
	fmt.Println("big pog")
	ws.Init()
	go ws.Main.Run()

	http.HandleFunc("/ws", ws.Main.ServeHub)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
