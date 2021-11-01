package models

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	ID primitive.ObjectID

	wsServer *WsServer
	conn     *websocket.Conn
	send     chan *MessageSocket
	user     *User

	save *chan SaveMessage

	Name  string `json:"name"`
	rooms map[*RoomSocket]bool
}

func newClient(conn *websocket.Conn, wsServer *WsServer, id primitive.ObjectID, sm *chan SaveMessage) *Client {
	return &Client{
		ID:       primitive.NewObjectID(),
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan *MessageSocket),
		rooms:    make(map[*RoomSocket]bool),
		save:     sm,
	}

}

func (c *Client) read() {
	defer c.conn.Close()
	for {
		var mgs *MessageSocket
		err := c.conn.ReadJSON(&mgs)
		if err != nil {
			log.Println(err)
			return
		}
		c.handleNewMessage(mgs)
	}
}

func (c *Client) write() {
	defer c.conn.Close()
	for msg := range c.send {
		err := c.conn.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}

func (client *Client) handleNewMessage(jsonMessage *MessageSocket) {

	// jsonMessage.Sender = client
	switch jsonMessage.Action {
	case SendMessageAction:

		roomID := jsonMessage.RoomID

		room := client.wsServer.findRoomByID(roomID)
		fmt.Println(roomID, room)

		if room := client.wsServer.findRoomByID(roomID); room != nil {
			jsonMessage.ID = primitive.NewObjectID()
			room.broadcast <- jsonMessage
			sm := &SaveMessage{
				message: &Message{
					ID:          jsonMessage.ID,
					RoomID:      jsonMessage.RoomID,
					SenderID:    jsonMessage.SenderID,
					Content:     jsonMessage.Content,
					Attachments: jsonMessage.Attachments,
					CreateAt:    jsonMessage.CreateAt,
				},
			}
			fmt.Println(sm)
			*client.save <- *sm
		}

	case JoinRoomAction:
		fmt.Println(jsonMessage.RoomID)
		client.handleJoinRoomMessage(*jsonMessage)
	case LeaveRoomAction:
		client.handleLeaveRoomMessage(*jsonMessage)

	}

}

func (client *Client) handleLeaveRoomMessage(message MessageSocket) {
	room := client.wsServer.findRoomByID(message.RoomID)
	if room == nil {
		return
	}

	if _, ok := client.rooms[room]; ok {
		delete(client.rooms, room)
	}

	room.unregister <- client
}
func (client *Client) handleJoinRoomMessage(message MessageSocket) {
	roomID := message.RoomID
	fmt.Println("handleJoinRoomMessage", roomID)
	client.joinRoom(roomID, nil)
}

func (client *Client) joinRoom(roomID primitive.ObjectID, sender *Client) {

	room := client.wsServer.findRoomByID(roomID)
	if room == nil {

		room = client.wsServer.createRoom(roomID, sender != nil)
	}

	// Don't allow to join private rooms through public room message
	if sender == nil && room.Private {
		return
	}

	if !client.isInRoom(room) {

		client.rooms[room] = true
		room.register <- client

	}

}
func (client *Client) isInRoom(room *RoomSocket) bool {
	if _, ok := client.rooms[room]; ok {
		return true
	}

	return false
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

// ServeWs handles websocket requests from clients requests.
func ServeWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {

	sm := NewSaveMessageChan()

	id, ok := r.URL.Query()["id"]

	if !ok || len(id[0]) < 1 {
		log.Println("Url Param 'id rooms' is missing")
		return
	}

	a, _ := primitive.ObjectIDFromHex(id[0])

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, wsServer, a, sm)

	go client.write()
	go client.read()

	wsServer.register <- client
}
