//Thanks to the author of https://scotch.io/bar-talk/build-a-realtime-chat-server-with-go-and-websockets
//for a nice template to build off of.
package main

import (
	"bufio"
	"log"
	"net/http"
	"os"

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
	port := getPort()
	server.Run(port)
}

func newServer() *Server {
	var s Server
	s.clients = make(map[*websocket.Conn]bool)
	s.messageCh = make(chan Message)
	return &s
}

func getPort() string {
	const defaultPort string = "8000"
	log.Printf("Default %s\nPort:", defaultPort)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	port := scanner.Text()
	if port == "" {
		return defaultPort
	}
	return port
}

//Run starts the server
func (server *Server) Run(port string) {
	http.Handle("/", server)
	go server.handleMessages()

	log.Println("Starting server on", port)
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
