//Thanks to the author of https://scotch.io/bar-talk/build-a-realtime-chat-server-with-go-and-websockets
//for a nice template to build off of.
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//Server holds clients and the broadcast channel
type Server struct {
	clients   map[*websocket.Conn]bool
	messageCh chan Message
}

//Message contains the information for each message
type Message struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func newServer() *Server {
	var s Server
	s.clients = make(map[*websocket.Conn]bool)
	s.messageCh = make(chan Message)
	return &s
}

func (server *Server) handleMessages() {
	for {
		message := <-server.messageCh
		server.broadcastMessage(message)
	}
}

func (server *Server) broadcastMessage(message Message) {
	for client := range server.clients {
		err := client.WriteJSON(message)
		if err != nil {
			log.Println("error: ", err)
			client.Close()
			delete(server.clients, client)
		}
	}
}

// Configure the upgrader
var upgrader = websocket.Upgrader{}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Got here....")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	server.clients[ws] = true

	for {
		var message Message
		err := ws.ReadJSON(&message)
		if err != nil {
			log.Println("error: ", err)
			delete(server.clients, ws)
			break
		}
		server.messageCh <- message
	}
}

//Run starts the server
func (server *Server) Run() {
	http.Handle("/", server)

	go server.handleMessages()

	log.Println("Server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	var server = newServer()
	server.Run()
}
