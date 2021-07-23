package models

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	room   primitive.ObjectID
	socket *websocket.Conn
	send   chan *Message

	channel *Channel
	user    *User
	save    *chan SaveMessage
}

func (c *Client) read() {
	defer c.socket.Close()
	for {
		var mgs *Message
		err := c.socket.ReadJSON(&mgs)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("mgs: ", mgs)
		c.channel.forward <- mgs

		sm := &SaveMessage{
			message: mgs,
		}
		*c.save <- *sm
	}
}
func (c *Client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}
