FROM golang:latest

EXPOSE 8000

WORKDIR /server
ADD . /server

RUN go get github.com/gorilla/websocket
RUN go build -o /server/s server/server.go
CMD /server/s
