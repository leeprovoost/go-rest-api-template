package main

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

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
