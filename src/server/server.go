package server

import (
	"log"
	"net/http"

	"../message"
	"github.com/gorilla/websocket"
)

type Server struct {
	clients   map[*websocket.Conn]bool
	broadcast chan message.Message
}

func NewServer() *Server {
	var s Server
	s.clients = make(map[*websocket.Conn]bool)
	s.broadcast = make(chan message.Message)
	return &s
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	clients[ws] = true

	for {
		var message message.Message
		err := ws.ReadJSON(&message)
		if err != nil {
			log.Println("error: ", err)
			delete(clients, ws)
			break
		}
		broadcast <- message
	}
}

func handleMessages() {
	for {
		message := <-broadcast
		for client := range clients {
			err := client.WriteJSON(message)
			if err != nil {
				log.Println("error: ", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func run(server Server) {
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	log.Println("Server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
