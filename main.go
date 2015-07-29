package main

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
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
	UserList  []User
	MaxUserId int
}

var db Database

// var userList []User
// var maxUserId int // to mimic database Id
var Render *render.Render

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Nothing to see here. #kthxbai")
}

func HealthcheckHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "HandleHealthchecks")
}

func ListUsersHandler(w http.ResponseWriter, req *http.Request) {
	//Render.JSON(w, http.StatusOK, ServiceGetUserList)
	responseObject := make(map[string][]User)
	responseObject["users"] = db.UserList
	Render.JSON(w, http.StatusOK, responseObject)
}

func GetUserHandler(w http.ResponseWriter, req *http.Request) {
	//ServiceGetUserById(1)
	fmt.Fprintf(w, "TO DO")
}

func CreateUserHandler(w http.ResponseWriter, req *http.Request) {
	// userList = append(userList, User{3, "Davide", "Tassinari", "01-01-1992", "Bologna"})
	// responseObject := make(map[string][]User)
	// responseObject["users"] = userList
	// Render.JSON(w, http.StatusOK, responseObject)
}

func UpdateUserHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "TO DO")
}

func DeleteUserHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "TO DO")
}

func PassportsHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Handling Passports")
}

func main() {
	// Initialise mock database
	userList := make([]User, 0)
	userList = append(userList, User{0, "John", "Doe", "31-12-1985", "London"})
	userList = append(userList, User{1, "Jane", "Doe", "01-01-1992", "Milton Keynes"})
	db = Database{userList, 1}

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
	n.UseHandler(router)

	fmt.Println("Starting server on :3009")
	n.Run(":3009")
}

// DB helper functions, move later to another file

// func ServiceGetUserList() map[string][]User {
// 	responseObject := make(map[string][]User)
// 	responseObject["users"] = db.UserList
// 	return responseObject
// }

// func ServiceGetUserById(id int) User {
// 	var u User
// 	for _, value := range userList {
// 		if value.Id == id {

// 		}
// 	}
// 	return u
// }

// func DbCreateUser(u User) map[string][]User {
// 	userList = append(userList, User{3, "Davide", "Tassinari", "01-01-1992", "Bologna"})
// 	responseObject := make(map[string][]User)
// 	responseObject["users"] = userList
// 	return responseObject
// }
