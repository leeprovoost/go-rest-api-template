package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// makeHandler allows us to pass an environment struct to our handlers, without resorting to global
// variables. It accepts an environment (Env) struct and our own handler function. It returns
// a function of the type http.HandlerFunc so can be passed on to the HandlerFunc in main.go.
func makeHandler(env env, fn func(http.ResponseWriter, *http.Request, env)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, env)
	}
}

func HomeHandler(w http.ResponseWriter, req *http.Request, env env) {
	log.Println("Home - Not implemented yet")
	env.Render.Text(w, http.StatusNotImplemented, "")
}

func HealthcheckHandler(w http.ResponseWriter, req *http.Request, env env) {
	env.Render.Text(w, http.StatusNoContent, "")
}

func MetricsHandler(w http.ResponseWriter, req *http.Request, env env) {
	stats := env.Metrics.Data()
	env.Render.JSON(w, http.StatusOK, stats)
}

func ListUsersHandler(w http.ResponseWriter, req *http.Request, env env) {
	list, _ := db.List()
	env.Render.JSON(w, http.StatusOK, list)
}

func GetUserHandler(w http.ResponseWriter, req *http.Request, env env) {
	vars := mux.Vars(req)
	uid, _ := strconv.Atoi(vars["uid"])
	user, err := db.Get(uid)
	if err == nil {
		env.Render.JSON(w, http.StatusOK, user)
	} else {
		env.Render.JSON(w, http.StatusNotFound, err)
	}
}

func CreateUserHandler(w http.ResponseWriter, req *http.Request, env env) {
	decoder := json.NewDecoder(req.Body)
	var u User
	err := decoder.Decode(&u)
	if err != nil {
		env.Render.JSON(w, http.StatusBadRequest, err)
	} else {
		user := User{-1, u.FirstName, u.LastName, u.DateOfBirth, u.LocationOfBirth}
		user, _ = db.Add(user)
		env.Render.JSON(w, http.StatusCreated, user)
	}
}

func UpdateUserHandler(w http.ResponseWriter, req *http.Request, env env) {
	decoder := json.NewDecoder(req.Body)
	var u User
	err := decoder.Decode(&u)
	if err != nil {
		env.Render.JSON(w, http.StatusBadRequest, err)
	} else {
		user := User{u.Id, u.FirstName, u.LastName, u.DateOfBirth, u.LocationOfBirth}
		user, err = db.Update(user)
		if err != nil {
			env.Render.JSON(w, http.StatusOK, user)
		} else {
			env.Render.JSON(w, http.StatusNotFound, err)
		}
	}
}

func DeleteUserHandler(w http.ResponseWriter, req *http.Request, env env) {
	vars := mux.Vars(req)
	uid, _ := strconv.Atoi(vars["uid"])
	ok, err := db.Delete(uid)
	if ok {
		// TO DO return empty body?
		env.Render.Text(w, http.StatusNoContent, "")
	} else {
		env.Render.JSON(w, http.StatusNotFound, err)
	}
}

func PassportsHandler(w http.ResponseWriter, req *http.Request, env env) {
	log.Println("Handling Passports - Not implemented yet")
	env.Render.Text(w, http.StatusNotImplemented, "")
}
