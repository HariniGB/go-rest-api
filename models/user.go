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

// type  (
// sessionData struct {
//   User
//   LoggedIn  bool `json:"LoggedIn" bson:"LoggedIn"`
//   LoginFail bool `json:"LoginFail" bson:"LoginFail"`
// }
// )
