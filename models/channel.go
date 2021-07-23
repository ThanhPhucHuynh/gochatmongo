package models

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Channel struct {
	forward chan *Message
	join    chan *Client
	leave   chan *Client

	clients map[*Client]bool
}

func (r *Channel) Run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
			}
		case msg := <-r.forward:
			for client := range r.clients {
				client.send <- msg
			}

		}
	}
}

func NewChanRoom() *Channel {
	r := &Channel{
		forward: make(chan *Message),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		clients: make(map[*Client]bool),
	}

	return r
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize,
	WriteBufferSize: socketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ChannelChat(c *Channel, sm *chan SaveMessage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		Parameter := r.URL.Query().Get("room")

		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("ServeHTTP:", err)
			return
		}
		a, _ := primitive.ObjectIDFromHex(Parameter)

		client := &Client{
			socket:  socket,
			send:    make(chan *Message, messageBufferSize),
			channel: c,
			user:    &User{},
			room:    a,
			save:    sm,
		}

		fmt.Println(c.clients)

		c.join <- client
		defer func() {
			c.leave <- client
		}()

		go client.write()
		client.read()
	}
}
