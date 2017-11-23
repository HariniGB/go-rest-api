package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/HariniGB/login-provider/controllers"
	"github.com/julienschmidt/httprouter"
)

func main() {
	// Instantiate a new router
	r := httprouter.New()

	username := os.Getenv("ROOT_USER")
	password := os.Getenv("ROOT_PASSWORD")
	host := os.Getenv("LDAP_HOST")
	portStr := os.Getenv("LDAP_PORT")
	port, _ := strconv.ParseInt(portStr, 10, 64)
	dn := os.Getenv("LDAP_DN")

	// Get a UserController instance
	uc := controllers.NewUserController(username, password, host, int(port), dn)

	// sign up page
	r.GET("/signup", uc.Signup)

	// Login page
	r.GET("/login", uc.Login)

	// Create a new user
	r.POST("/api/v1/user", uc.CreateUser)

	//  Update a user
	r.PUT("/api/v1/user/:id", uc.UpdateUser)

	// Remove an existing user
	r.DELETE("/api/v1/user/:id", uc.RemoveUser)

	// Fire up the server
	http.ListenAndServe("localhost:3000", r)
}
