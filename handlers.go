package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	Render.Text(w, http.StatusNoContent, "")
}

func HealthcheckHandler(w http.ResponseWriter, req *http.Request) {
	Render.Text(w, http.StatusNoContent, "")
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
	decoder := json.NewDecoder(req.Body)
	var u User
	err := decoder.Decode(&u)
	if err != nil {
		Render.JSON(w, http.StatusBadRequest, err)
	} else {
		user := User{u.Id, u.FirstName, u.LastName, u.DateOfBirth, u.LocationOfBirth}
		user, err = db.Update(user)
		if err != nil {
			Render.JSON(w, http.StatusOK, user)
		} else {
			Render.JSON(w, http.StatusNotFound, err)
		}
	}
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
	fmt.Fprintf(w, "Handling Passports - Not implemented yet")
	Render.Text(w, http.StatusNotImplemented, "")
}
