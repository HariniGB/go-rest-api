package models

import "gopkg.in/mgo.v2/bson"

type (
  // User represents the structure of our resource
  User struct {
    Id        bson.ObjectId `json:"id" bson:"_id"`
    Name      string        `json:"name" bson:"name"`
    Email     string        `json:"email" bson:"email"`
    Password  string        `json:"password" bson:"password"`
  }
)

type  (
Url struct {
  Id         bson.ObjectId `json:"id" bson:"_id"`
  oURL       string        `json:"oURL" bson:"oURL"`
  sURL       string        `json:"sURL" bson:"sURL"`
  userId     string        `json:"userId" bson:"userId"`
}
)
