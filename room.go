package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	tracer "github.com/neel229/tracing"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: messageBufferSize,
}

type room struct {
	// forward is a channel which takes the
	// incoming messages which needs to be
	// forwarded to other clients
	forward chan []byte
	// join is a channel for clients to
	// join the chat-room
	join chan *client
	// leave is a channel for clients to
	// leave the chat-room
	leave chan *client
	// clients holds the pool of all the
	// clients currently present in the room
	clients map[*client]bool
	// tracer will receive trace information
	// of the activity in the room
	tracer tracer.Tracer
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// add the client to the clients pool
			r.clients[client] = true
			r.tracer.Trace("New client joined")
		case client := <-r.leave:
			// remove the client from the clients pool
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left")
		case msg := <-r.forward:
			r.tracer.Trace("Message received: ", string(msg))
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace("--- sent to client ---")
			}
		}
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
