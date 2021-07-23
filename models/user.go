package models

import (
	"github.com/globalsign/mgo/bson"
)

// User hold information about an user
type User struct {
	ID     bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Email  string        `json:"email" bson:"email"`
	Name   string        `json:"name" bson:"name"`
	Active bool          `json:"active" bson:"active"`
}
