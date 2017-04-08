FROM golang:latest
ARG path=$GOPATH
EXPOSE 8000

RUN echo $path

ADD . $path/src/github.com/connormurray7/go-chat-app/server
WORKDIR $path/src/github.com/connormurray7/go-chat-app/server

RUN go build $path/src/github.com/connormurray7/go-chat-app/server/server.go
CMD $path/src/github.com/connormurray7/go-chat-app/server/server

