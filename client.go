package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a single client
type client struct {
	// socket is the websocket for the client
	socket *websocket.Conn
	// send is a channel on which messages
	// has to be sent
	send chan *message
	// room is the room the client is chatting in
	room *room
	// userdata holds the data of the user
	userData map[string]interface{}
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
}
