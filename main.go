// Example REST API for managing passports
package main

import (
	"fmt"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/unrolled/render"
)

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
