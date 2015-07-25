package main

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

type Passport struct {
	Id              string `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	DateOfBirth     string `json:"date_of_birth"`
	LocationOfBirth string `json:"location_of_birth"`
}

func PassportsHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "List of passports")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/passports", PassportsHandler).Methods("GET")

	n := negroni.Classic()
	n.UseHandler(router)

	fmt.Println("Starting server on :3009")
	n.Run(":3009")
}
