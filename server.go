package main

// See this IMPORTANT LINK FOR REFERENCE:  https://stevenwhite.com/tag/golang/

import (
  // Standard library packages
  "net/http"
  // Connecting to the controller in same folder
  "github.com/HariniGB/go-rest-api/controllers"
  // Third party packages
  "github.com/julienschmidt/httprouter"
  "gopkg.in/mgo.v2"

)

func main() {
  // Instantiate a new router
  r := httprouter.New()

  // Get a UserController instance
  uc := controllers.NewUserController(getSession())

  // Get all users resources
  // r.GET("/users", uc.GetUsers)

  // Home page
   r.GET("/", uc.Home)

  // sign up page
  r.GET("/signup", uc.Signup)

  // Login page
  r.GET("/login", uc.Login)

  // Get a user resource
  r.GET("/api/v1/user/:id", uc.GetUser)

  // Create a new user
  r.POST("/api/v1/user", uc.CreateUser)

  //  Update a user
  r.PUT("/api/v1/user/:id", uc.UpdateUser)

  // Remove an existing user
  r.DELETE("/api/v1/user/:id", uc.RemoveUser)

  // Fire up the server
  http.ListenAndServe("localhost:3000", r)
}

// getSession creates a new mongo session and panics if connection error occurs
func getSession() *mgo.Session {
  // Connect to our local mongo
  s, err := mgo.Dial("mongodb://localhost:27017")

  // Check if connection error, is mongo running?
  if err != nil {
    panic(err)
  }

  // Deliver session
  return s
}


