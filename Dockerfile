FROM golang

ADD ./go/src/github.com/connormurray7/go-chat-app

RUN go install github.com/connormurray7/go-chat-app
RUN go build github.com/connormurray7/go-chat-app/server/server.go
RUN go build github.com/connormurray7/go-chat-app/client/client.go

CMD github.com/connormurray7/go-chat-app/server/server

EXPOSE 8000
