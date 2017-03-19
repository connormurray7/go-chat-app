package main

import (
	// "log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{}

type Message struct {
	email    string `json:"email"`
	username string `json: "username"`
	body     string `json: "body"`
}

func main() {
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)
}
