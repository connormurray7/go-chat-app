package server

import (
	"log"
	"net/http"

	"../message"
	"github.com/gorilla/websocket"
)

//Server holds clients andthe broadcast channel
type Server struct {
	clients   map[*websocket.Conn]bool
	broadcast chan message.Message
}

//NewServer creates a new server.
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

func (server *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println("Received something")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	server.clients[ws] = true

	for {
		var message message.Message
		err := ws.ReadJSON(&message)
		log.Println("This is a message:", message)
		if err != nil {
			log.Println("error: ", err)
			delete(server.clients, ws)
			break
		}
		server.broadcast <- message
	}
}

func (server *Server) handleMessages() {
	for {
		message := <-server.broadcast
		log.Println("Got a message:", message)
		for client := range server.clients {
			err := client.WriteJSON(message)
			if err != nil {
				log.Println("error: ", err)
				client.Close()
				delete(server.clients, client)
			}
		}
	}
}

//Run starts the server
func (server *Server) Run() {
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", server.handleConnections)

	go server.handleMessages()

	log.Println("Server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
