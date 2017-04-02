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

func main() {
	var server = newServer()
	const port string = "8000" //Can swap out for config file or user input.
	server.Run(port)
}

func newServer() *Server {
	var s Server
	s.clients = make(map[*websocket.Conn]bool)
	s.messageCh = make(chan Message)
	return &s
}

//Run starts the server
func (server *Server) Run(port string) {
	http.Handle("/", server)
	go server.handleMessages()

	log.Println("Starting server on port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
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
			log.Println("error:", err)
			client.Close()
			delete(server.clients, client)
		}
	}
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	server.clients[ws] = true
	log.Println("New user connected from:", r.RemoteAddr)

	for {
		var message Message
		err := ws.ReadJSON(&message)
		if err != nil {
			log.Println("error:", err)
			delete(server.clients, ws)
			break
		}
		server.messageCh <- message
	}
}
