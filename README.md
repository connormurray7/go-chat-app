# Go chat
Go chat is a simple chat program that can be run from the terminal. It uses websockets to communicate between the client and server. This is my first experience using Go and the starting point was [this tutorial](https://scotch.io/bar-talk/build-a-realtime-chat-server-with-go-and-websockets).

# Install
	go get github.com/connormurray7/go-chat-app

# Run
    go run server/server.go
    go run client/client.go # Another session

# Run using docker
Within the directory run,

    docker build .
    docker run <image>
    go run client/client.go # Another session

# Example
### Server
The server defaults to serving `localhost:8000`. The server logs when a new user is connected and from what remote address.
![Alt text](example/server.png?raw=true "Optional Title")
### Clients
Clients default to`localhost:8000` for connecting to the server. Every client that is connected to the server will receive messages. Currently all clients will be broadcasted messages, but in the future I may the options for different feeds.. 
![Alt text](example/connor.png?raw=true "Optional Title")
![Alt text](example/other.png?raw=true "Optional Title") 
