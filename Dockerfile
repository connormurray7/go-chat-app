FROM golang:latest

EXPOSE 8000

WORKDIR /server
ADD . /server

RUN go build -o /server/server server/server.go
CMD /server/server
