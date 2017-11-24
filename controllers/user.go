package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/HariniGB/login-provider/common"
	"github.com/HariniGB/login-provider/ldap"
	"github.com/HariniGB/login-provider/storage"
	"github.com/julienschmidt/httprouter"
	"github.com/sony/sonyflake"
)

const (
	s3url = "s3url"
)

type UserController struct {
	lp *ldap.Ldap
	st storage.Storage
	fl *sonyflake.Sonyflake
	to time.Duration
}

// NewUserController provides a reference to a UserController with provided mongo session
func NewUserController(username, password, host string,
	port int, dn string, st storage.Storage, timeout time.Duration) *UserController {
	lp, err := ldap.NewLdap(username, password, host, port, dn)
	if err != nil {
		return nil
	}
	fl := sonyflake.NewSonyflake(sonyflake.Settings{StartTime: time.Now()})
	return &UserController{lp: lp, st: st, fl: fl, to: timeout}
}

// Sign up retrieves the signup form for new users
func (uc UserController) Signup(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	tmpl, err := template.ParseFiles("templates/signup.html")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, "Signup page")
}

func (uc UserController) AuthGet(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	_, info := uc.getUserInfoFromCookie(r)
	if info != nil {
		uc.setHttpHeaders(w, info)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Login successful")
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, "Forbidden")
	return
}

// Login retrieves the login form for the users
func (uc UserController) Auth(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	target := r.Header.Get("Referer")
	id, info := uc.getUserInfoFromCookie(r)
	if info != nil {
		// Delete old entry
		uc.st.Delete(id)
	}

	u := common.User{}
	if r.Header.Get("Content-Type") == "application/json" {
		json.NewDecoder(r.Body).Decode(&u)
	} else {
		u.Username = r.FormValue("id")
		u.Password = r.FormValue("password")
	}

	err := u.ValidateLogin()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
	}

	if uc.lp.ExistsUser(u.Username) == false {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User %s doesn't exist", u.Username)
		return
	}

	if uc.lp.Validate(u.Username, u.Password) == false {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Invalid password for username %s", u.Username)
		return
	}

	info = &common.UserInfo{
		Id: u.Username,
	}
	if groups, ok := uc.lp.GetUsersGroups(u.Username); ok {
		info.Groups = groups
	}

	uid, err := uc.fl.NextID()
	if err != nil {
		log.Println("Unable to generate ID due to error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to authenticate user %s. Please try again", u.Username)
		return
	}

	id = fmt.Sprintf("%d", uid)
	err = uc.st.Insert(id, info)
	if err != nil {
		log.Println("Unable to persist user info due to error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to authenticate user %s. Please try again", u.Username)
		return

	}
	cookie := &http.Cookie{
		Name:     s3url,
		HttpOnly: true,
		Secure:   false,
		Value:    id,
		Expires:  time.Now().Add(uc.to),
	}

	uc.setHttpHeaders(w, info)
	http.SetCookie(w, cookie)
	fmt.Println(target)
	if target != "" {
		http.Redirect(w, r, target, http.StatusFound)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Login successful")
	}

}

// Login retrieves the login form for the users
func (uc UserController) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, r.Header.Get("Referer"))
}

// CreateUser creates a new user resource
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	u := common.User{}
	if r.Header.Get("Content-Type") == "application/json" {
		json.NewDecoder(r.Body).Decode(&u)
	} else {
		u.Username = r.FormValue("id")
		u.Email = r.FormValue("email")
		u.FirstName = r.FormValue("first_name")
		u.LastName = r.FormValue("last_name")
		u.Password = r.FormValue("password")
	}

	err := u.ValidateSignup()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
	}

	if uc.lp.ExistsUser(u.Username) == true {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "User %s already exists", u.Username)
		return
	}

	err = uc.lp.AddUser(&u)
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
	u := common.User{}
	if r.Header.Get("Content-Type") == "application/json" {
		json.NewDecoder(r.Body).Decode(&u)
	} else {
		u.Username = r.FormValue("id")
		u.Email = r.FormValue("email")
		u.FirstName = r.FormValue("first_name")
		u.LastName = r.FormValue("last_name")
		u.Password = r.FormValue("password")
	}

	err := u.ValidateSignup()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
	}

	if uc.lp.ExistsUser(u.Username) == false {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User %s doesn't exist", u.Username)
		return
	}

	err = uc.lp.UpdateUser(&u)
	if err == ldap.UserAuthFailure {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Invalid password for username %s", u.Username)
		return
	}

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

func (uc UserController) getUserInfoFromCookie(r *http.Request) (string, *common.UserInfo) {
	c, err := r.Cookie(s3url)

	// Check if cookie is present or not
	if err == nil {
		id := c.Value
		// If there is a valid entry in meta store then mark as success
		info, _ := uc.st.Get(id)
		return id, info
	}

	return "", nil
}

func (uc UserController) setHttpHeaders(w http.ResponseWriter, u *common.UserInfo) {
	w.Header().Set("X-Forwarded-User", u.Id)
	w.Header().Set("X-Forwarded-Groups", strings.Join(u.Groups, "|"))
}
