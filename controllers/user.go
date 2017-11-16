package controllers

import (
  "encoding/json"
  "fmt"
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/HariniGB/go-rest-api/models"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "html/template"
  "golang.org/x/crypto/bcrypt"
)

type (
  // UserController represents the controller for operating on the User resource
  UserController struct {
    session *mgo.Session
  }
)

// NewUserController provides a reference to a UserController with provided mongo session
func NewUserController(s *mgo.Session) *UserController {
  return &UserController{s}
}

//Method to get hash password for user password and show how to hash passwords using bcrypt
func HashPassword(password string) (string, error) {
  bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
  return string(bytes), err
}
//  Method to check the hashed password and user given password while login
func CheckPasswordHash(password, hash string) bool {
  err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
  return err == nil
}

// Home retrieves the home page
func (uc UserController) Home(w http.ResponseWriter, r *http.Request, p httprouter.Params){
  tmpl, err := template.ParseFiles("home.html")
  if err != nil {
    panic(err)
  }
  tmpl.Execute(w,"Home page")
}

// Sign up retrieves the signup form for new users
func (uc UserController) Signup(w http.ResponseWriter, r *http.Request, p httprouter.Params){
  tmpl, err := template.ParseFiles("signup.html")
  if err != nil {
    panic(err)
  }
  tmpl.Execute(w,"Signup page")
}

// Login retrieves the login form for the users
func (uc UserController) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params){
  tmpl, err := template.ParseFiles("login.html")
  if err != nil {
    panic(err)
  }
   tmpl.Execute(w, "Login page")
}

// CreateUser creates a new user resource
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  // Stub an user to be populated from the body
  u := models.User{}
  r.ParseForm()

  // Populate the user data
  json.NewDecoder(r.Body).Decode(&u)

  // Add an Id
  u.Id = bson.NewObjectId()
  u.Name = r.FormValue("name")
  u.Email = r.FormValue("email")
  password, _:= HashPassword(r.FormValue("password"))
  u.Password = password

  // Testing Bcrypt
  fmt.Println("Password:", r.FormValue("password"))
  fmt.Println("Hash:    ", u.Password)
  match := CheckPasswordHash(r.FormValue("password"), u.Password)
  fmt.Println("Match:   ", match)

  // Write the user to mongo
  uc.session.DB("go_rest").C("users").Insert(u)
  fmt.Printf("\ncreated",u)
  // Marshal provided interface into JSON structure
  uj, _ := json.Marshal(u)

  // Write content-type, status code, payload
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(201)
  fmt.Fprintf(w, "%s", uj)

}

// GetUsers retrieves all the user's resources
func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  // Grab id
  id := p.ByName("id")

  // Verify id is ObjectId, otherwise bail
  if !bson.IsObjectIdHex(id) {
    w.WriteHeader(404)
    return
  }

  // Grab id
  oid := bson.ObjectIdHex(id)

  // Stub user
  u := models.User{}

  // Fetch user
  if err := uc.session.DB("go_rest").C("users").FindId(oid).One(&u); err != nil {
    w.WriteHeader(404)
    return
  }

  // Marshal provided interface into JSON structure
  uj, _ := json.Marshal(u)

  // Write content-type, status code, payload
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(200)
  fmt.Fprintf(w, "%s", uj)
}

// UpdateUser updates the user resource
func (uc UserController) UpdateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  // Grab id
  id := p.ByName("id")

  // Verify id is ObjectId, otherwise bail
  if !bson.IsObjectIdHex(id) {
    w.WriteHeader(404)
    return
  }

  // Grab id
  oid := bson.ObjectIdHex(id)

  // Stub user
  u := models.User{}

  // Fetch user
  if err := uc.session.DB("go_rest").C("users").FindId(oid).One(&u); err != nil {
    w.WriteHeader(404)
    return
  }

  // Populate the user data
  json.NewDecoder(r.Body).Decode(&u)

  //Update the user to mongo
  uc.session.DB("go_rest").C("users").Update(oid, u)

  // Marshal provided interface into JSON structure
  uj, _ := json.Marshal(u)

  // Write content-type, status code, payload
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(201)
  fmt.Fprintf(w, "%s", uj)
}

// RemoveUser removes an existing user resource
func (uc UserController) RemoveUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  // Grab id
  id := p.ByName("id")

  // Verify id is ObjectId, otherwise bail
  if !bson.IsObjectIdHex(id) {
    w.WriteHeader(404)
    return
  }

  // Grab id
  oid := bson.ObjectIdHex(id)

  // Remove user
  if err := uc.session.DB("go_rest").C("users").RemoveId(oid); err != nil {
    w.WriteHeader(404)
    return
  }

  // Write status
  w.WriteHeader(200)
}




// GetUser retrieves an individual user resource
// func (uc UserController) GetUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

//   // Stub user
//   u := models.User{}

//   // Fetch user
//   if err := uc.session.DB("go_rest").C("users"); err != nil {
//     w.WriteHeader(404)
//     return
//   }

//   // Marshal provided interface into JSON structure
//   uj, _ := json.Marshal(u)

//   // Write content-type, status code, payload
//   w.Header().Set("Content-Type", "application/json")
//   w.WriteHeader(200)
//   fmt.Fprintf(w, "%s", uj)
// }

