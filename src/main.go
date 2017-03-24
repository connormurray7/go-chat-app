package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//Server holds clients and the broadcast channel
type Server struct {
	clients   map[*websocket.Conn]bool
	broadcast chan Message
}

//Message contains the information for each message
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func newServer() *Server {
	var s Server
	s.clients = make(map[*websocket.Conn]bool)
	s.broadcast = make(chan Message)
	return &s
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

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	server.clients[ws] = true

	for {
		var message Message
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

//Run starts the server
func (server *Server) Run() {
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	http.Handle("/ws", server)

	go server.handleMessages()

	log.Println("Server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	var server = newServer()
	server.Run()
}
