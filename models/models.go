package models

import "go.mongodb.org/mongo-driver/bson/primitive"

//Create Struct
type Test struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	User string             `json:"user" bson:"user,omitempty"`
}

// type User struct {
// 	FirstName string `json:"firstname,omitempty" bson:"firstname,omitempty"`
// 	LastName  string `json:"lastname,omitempty" bson:"lastname,omitempty"`
//}
