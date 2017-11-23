package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/HariniGB/login-provider/controllers"
	"github.com/HariniGB/login-provider/storage"
	"github.com/HariniGB/login-provider/storage/memcache"
	"github.com/HariniGB/login-provider/storage/memory"
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

	timeoutStr := os.Getenv("SESSION_EXPIRY")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		log.Printf("Invalid session expiry %s. Defaulting to 1h", timeoutStr)
		timeout = time.Hour
	}

	memHost := os.Getenv("MEMCACHE_HOST")
	memPort := os.Getenv("MEMCACHE_PORT")

	var st storage.Storage
	if memHost != "" && memPort != "" {
		st = memcache.NewMemoryStorage(memHost, memPort, timeout)
	} else {
		st = memory.NewMemoryStorage(1000, timeout)
	}

	// Get a UserController instance
	uc := controllers.NewUserController(username, password, host, int(port), dn, st, timeout)

	// sign up page
	r.GET("/signup", uc.Signup)

	// Login page
	r.GET("/login", uc.Login)

	// Authenticate user
	r.POST("/auth", uc.Auth)
	r.GET("/auth", uc.AuthGet)

	// Create a new user
	r.POST("/api/v1/user", uc.CreateUser)

	//  Update a user
	r.PUT("/api/v1/user/:id", uc.UpdateUser)

	// Remove an existing user
	r.DELETE("/api/v1/user/:id", uc.RemoveUser)

	// Fire up the server
	http.ListenAndServe("localhost:3000", r)
}
