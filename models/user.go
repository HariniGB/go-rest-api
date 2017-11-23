package models

import "gopkg.in/mgo.v2/bson"

type User struct {
	Username  string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Url struct {
	Id     bson.ObjectId `json:"id" bson:"_id"`
	oURL   string        `json:"oURL" bson:"oURL"`
	sURL   string        `json:"sURL" bson:"sURL"`
	userId string        `json:"userId" bson:"userId"`
}
