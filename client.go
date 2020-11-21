package main

import (
	"github.com/gorilla/websocket"
)

// Client represents a single client
type client struct {
	// socket is the websocket for the client
	socket *websocket.Conn
	// send is a channel on which messages
	// has to be sent
	send chan []byte
	// room is the room the client is chatting in
	room *room
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}
