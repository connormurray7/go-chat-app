package main

import (
	"net/http"

	"./server"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	s := server.NewServer()
	server.run(s)
}
