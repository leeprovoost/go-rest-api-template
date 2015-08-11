// Example REST API for managing passports
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/unrolled/render"
)

type User struct {
	Id              int    `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	DateOfBirth     string `json:"date_of_birth"`
	LocationOfBirth string `json:"location_of_birth"`
}

type Passport struct {
	Id           string `json:"id"`
	DateOfIssue  string `json:"date_of_issue"`
	DateOfExpiry string `json:"date_of_expiry"`
	Authority    string `json:"authority"`
	UserId       int    `json:"user_id"`
}

type Database struct {
	UserList  map[int]User
	MaxUserId int
}

// List returns a list of JSON documents
func (db *Database) List() map[string][]User {
	var list []User = make([]User, 0)
	for _, v := range db.UserList {
		list = append(list, v)
	}
	responseObject := make(map[string][]User)
	responseObject["users"] = list
	return responseObject
}

// Retrieve a single JSON document
func (db *Database) Get(i int) (User, error) {
	user, ok := db.UserList[i]
	if ok {
		return user, nil
	} else {
		return user, errors.New("User does not exist")
	}
}

// Add a User JSON document, returns the JSON document with the generated id
func (db *Database) Add(u User) User {
	db.MaxUserId = db.MaxUserId + 1
	newUser := User{
		Id:              db.MaxUserId,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		DateOfBirth:     u.DateOfBirth,
		LocationOfBirth: u.LocationOfBirth,
	}
	db.UserList[db.MaxUserId] = newUser
	return newUser
}

// Delete a user
func (db *Database) Delete(i int) (bool, error) {
	_, ok := db.UserList[i]
	if ok {
		delete(db.UserList, i)
		return true, nil
	} else {
		return false, errors.New("Could not delete this user")
	}
}

// Update an existing user
func (db *Database) Update(u User) (User, error) {
	id := u.Id
	user, ok := db.UserList[id]
	if ok {
		db.UserList[id] = user
		return db.UserList[id], nil
	} else {
		return user, errors.New("User does not exist")
	}
}

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Nothing to see here. #kthxbai")
}

func HealthcheckHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "HandleHealthchecks")
}

func ListUsersHandler(w http.ResponseWriter, req *http.Request) {
	Render.JSON(w, http.StatusOK, db.List())
}

func GetUserHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uid, _ := strconv.Atoi(vars["uid"])
	user, err := db.Get(uid)
	if err == nil {
		Render.JSON(w, http.StatusOK, user)
	} else {
		Render.JSON(w, http.StatusNotFound, err)
	}
}

func CreateUserHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var u User
	err := decoder.Decode(&u)
	if err != nil {
		Render.JSON(w, http.StatusBadRequest, err)
	} else {
		user := User{-1, u.FirstName, u.LastName, u.DateOfBirth, u.LocationOfBirth}
		user = db.Add(user)
		Render.JSON(w, http.StatusCreated, user)
	}
}

func UpdateUserHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "TO DO")
}

func DeleteUserHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uid, _ := strconv.Atoi(vars["uid"])
	ok, err := db.Delete(uid)
	if ok {
		// TO DO return empty body?
		Render.JSON(w, http.StatusOK, nil)
	} else {
		Render.JSON(w, http.StatusNotFound, err)
	}
}

func PassportsHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Handling Passports")
}

var db *Database
var Render *render.Render

func init() {
	list := make(map[int]User)
	list[0] = User{0, "John", "Doe", "31-12-1985", "London"}
	list[1] = User{1, "Jane", "Doe", "01-01-1992", "Milton Keynes"}
	db = &Database{list, 1}
}

func main() {
	Render = render.New()
	router := mux.NewRouter()

	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/healthcheck", HealthcheckHandler).Methods("GET")

	router.HandleFunc("/users", ListUsersHandler).Methods("GET")
	router.HandleFunc("/users/{uid}", GetUserHandler).Methods("GET")
	router.HandleFunc("/users", CreateUserHandler).Methods("POST")
	router.HandleFunc("/users/{uid}", UpdateUserHandler).Methods("PUT")
	router.HandleFunc("/users/{uid}", DeleteUserHandler).Methods("DELETE")

	router.HandleFunc("/users/{uid}/passports", PassportsHandler).Methods("GET")
	router.HandleFunc("/passports/{pid}", PassportsHandler).Methods("GET")
	router.HandleFunc("/users/{uid}/passports", PassportsHandler).Methods("POST")
	router.HandleFunc("/passports/{pid}", PassportsHandler).Methods("PUT")
	router.HandleFunc("/passports/{pid}", PassportsHandler).Methods("DELETE")

	n := negroni.Classic()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	n.Use(c)
	n.UseHandler(router)
	fmt.Println("Starting server on :3009")
	n.Run(":3009")
}
