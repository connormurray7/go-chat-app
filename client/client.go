package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

//Message contains the information for each message
type Message struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func main() {
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(address, nil)

	if err != nil {
		fmt.Println("Error connecting to server", err)
		return
	}
	fmt.Println("You are connected to")
	go writeMessages(conn)
	waitForMessages(conn)
}

func getAddress() {
	const address string = "ws://localhost:8000"
	fmt.Print("Address:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	fmt.Println("This was the address given", text)
}

func writeMessages(conn *websocket.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		m := Message{Name: "me", Message: text}
		conn.WriteJSON(m)
	}
}

func waitForMessages(conn *websocket.Conn) {
	for {
		var m Message
		err := conn.ReadJSON(&m)
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println(m.Name, "->", m.Message)
	}
}
