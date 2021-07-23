package models

import (
	db "chat/db"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rooms struct {
	ID           primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	Participants []primitive.ObjectID `json:"participants" bson:"participants"`
}

func RoomsList(w http.ResponseWriter, r *http.Request) {

	var data []Rooms

	c, err := db.ConnectMongoDB().Client.Database("test").Collection("rooms").Find(context.Background(), bson.D{})

	if err != nil {

		log.Print(err)
		Error(w, err, http.StatusBadGateway)
	}
	if err = c.All(context.Background(), &data); err != nil {
		Error(w, err, http.StatusBadGateway)
	}
	// return data, err

	JSON(w, http.StatusOK, data)

}

func NewRoomMongo(w http.ResponseWriter, r *http.Request) {
	// Write the user to mongo
	a, err := primitive.ObjectIDFromHex("60e3f9a7e1ab4c3dfc8fe4c1")
	b, err := primitive.ObjectIDFromHex("60e3b5dbe1ab4c388ce2d04c")

	room := Rooms{
		ID: primitive.NewObjectID(),
		Participants: []primitive.ObjectID{
			a,
			b,
		},
	}

	c, err := db.ConnectMongoDB().Client.Database("test").Collection("rooms").InsertOne(context.Background(), room)
	if err != nil {
		log.Print(err)
	}
	JSON(w, http.StatusOK, c)
}
func JSON(w http.ResponseWriter, status int, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
}

// Error write error, status to http response writer
func Error(w http.ResponseWriter, err error, status int) {
	http.Error(w, err.Error(), status)
}
