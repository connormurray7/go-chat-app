package main

import (
	"net/http"

	"./message"
	"./server"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	var clients = make(map[*websocket.Conn]bool)
	var broadcast = make(chan message.Message)

	server.run(clients, broadcast)
}
