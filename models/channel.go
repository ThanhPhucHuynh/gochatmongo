package models

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *MessageSocket
	rooms      map[*RoomSocket]bool
}

func NewWebsocketServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *MessageSocket),
		rooms:      make(map[*RoomSocket]bool),
	}
}

func (server *WsServer) Run() {
	for {
		select {

		case client := <-server.register:
			server.registerClient(client)

		case client := <-server.unregister:
			server.unregisterClient(client)

		case message := <-server.broadcast:
			server.broadcastToClients(message)
		}

	}
}

func (server *WsServer) registerClient(client *Client) {
	server.clients[client] = true
}

func (server *WsServer) unregisterClient(client *Client) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
	}
}

func (server *WsServer) broadcastToClients(messageSocket *MessageSocket) {
	for client := range server.clients {
		client.send <- messageSocket
	}
}

func (server *WsServer) findRoomByID(ID primitive.ObjectID) *RoomSocket {
	var foundRoom *RoomSocket
	fmt.Println(server.rooms)

	for room := range server.rooms {
		if room.GetId() == ID {
			fmt.Println(room.ID, room.GetId())
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (server *WsServer) createRoom(id primitive.ObjectID, private bool) *RoomSocket {
	room := NewRoom(id, private)
	go room.RunRoomSocket()
	server.rooms[room] = true

	return room
}

func (server *WsServer) findClientByID(ID primitive.ObjectID) *Client {
	var foundClient *Client
	for client := range server.clients {
		if client.ID == ID {
			foundClient = client
			break
		}
	}

	return foundClient
}
