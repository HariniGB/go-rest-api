package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/HariniGB/login-provider/ldap"
	"github.com/HariniGB/login-provider/models"
	"github.com/julienschmidt/httprouter"
)

type UserController struct {
	lp *ldap.Ldap
}

// NewUserController provides a reference to a UserController with provided mongo session
func NewUserController(username, password, host string, port int, dn string) *UserController {
	lp, err := ldap.NewLdap(username, password, host, port, dn)
	if err != nil {
		return nil
	}
	return &UserController{lp}
}

// Sign up retrieves the signup form for new users
func (uc UserController) Signup(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	tmpl, err := template.ParseFiles("templates/signup.html")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, "Signup page")
}

// Login retrieves the login form for the users
func (uc UserController) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, "Login page")
}

// CreateUser creates a new user resource
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	u := models.User{}
	if r.Header.Get("Content-Type") == "application/json" {
		json.NewDecoder(r.Body).Decode(&u)
	} else {
		u.Username = r.FormValue("name")
		u.Email = r.FormValue("email")
		u.FirstName = r.FormValue("first_name")
		u.LastName = r.FormValue("last_name")
		u.Password = r.FormValue("password")
	}

	if uc.lp.ExistsUser(u.Username) == true {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "User %s already exists", u.Username)
		return
	}

	err := uc.lp.AddUser(&u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to create user %s. Please try again", u.Username)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User %s created", u.Username)
}

// UpdateUser updates the user resource
func (uc UserController) UpdateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	u := models.User{}
	if r.Header.Get("Content-Type") == "application/json" {
		json.NewDecoder(r.Body).Decode(&u)
	} else {
		u.Username = r.FormValue("name")
		u.Email = r.FormValue("email")
		u.FirstName = r.FormValue("first_name")
		u.LastName = r.FormValue("last_name")
		u.Password = r.FormValue("password")
	}

	if uc.lp.ExistsUser(u.Username) == false {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User %s doesn't exist", u.Username)
		return
	}

	err := uc.lp.UpdateUser(&u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to update user %s. Please try again", u.Username)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %s updated", u.Username)
}

// RemoveUser removes an existing user resource
func (uc UserController) RemoveUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if uc.lp.ExistsUser(id) == false {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User %s doesn't exist", id)
		return
	}

	err := uc.lp.DeleteUser(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to delete user %s. Please try again", id)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %s updated", id)
}
