package models

import (
	"context"
	"log"
	"time"

	db "chat/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	RoomID      primitive.ObjectID `json:"room_id" bson:"room_id"`
	SenderID    primitive.ObjectID `json:"sender_id" bson:"sender_id"`
	Content     string             `json:"content" bson:"content"`
	Attachments string             `json:"attachments" bson:"attachments"`
	CreateAt    time.Time          `json:"created_at" bson:"created_at"`
}

type SaveMessage struct {
	message *Message `json:"message"`
}

func saveMessages(sm *chan SaveMessage) {
	for {
		sm, ok := <-*sm
		if !ok {
			log.Print("Error when receiving message to save")
			return
		}
		_, err := db.ConnectMongoDB().Client.Database("test").Collection("message").InsertOne(context.Background(), sm.message)
		if err != nil {
			log.Print(err)
		}
	}
}

// NewSaveMessageChan create a new SaveMessage channel
func NewSaveMessageChan() *chan SaveMessage {
	sm := make(chan SaveMessage, 256)
	go saveMessages(&sm)
	return &sm
}
