package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// HandlerFunc is a custom implementation of the http.HandlerFunc
type HandlerFunc func(http.ResponseWriter, *http.Request, appContext)

// makeHandler allows us to pass an environment struct to our handlers, without resorting to global
// variables. It accepts an environment (Env) struct and our own handler function. It returns
// a function of the type http.HandlerFunc so can be passed on to the HandlerFunc in main.go.
func makeHandler(ctx appContext, fn func(http.ResponseWriter, *http.Request, appContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, ctx)
	}
}

// HealthcheckHandler returns useful info about the app
func HealthcheckHandler(w http.ResponseWriter, req *http.Request, ctx appContext) {
	ctx.render.Text(w, http.StatusNoContent, "")
}

// MetricsHandler returns application metrics
func MetricsHandler(w http.ResponseWriter, req *http.Request, ctx appContext) {
	stats := ctx.metrics.Data()
	ctx.render.JSON(w, http.StatusOK, stats)
}

// ListUsersHandler returns a list of users
func ListUsersHandler(w http.ResponseWriter, req *http.Request, ctx appContext) {
	list, _ := db.List()
	ctx.render.JSON(w, http.StatusOK, list)
}

// GetUserHandler returns a user object
func GetUserHandler(w http.ResponseWriter, req *http.Request, ctx appContext) {
	vars := mux.Vars(req)
	uid, _ := strconv.Atoi(vars["uid"])
	user, err := db.Get(uid)
	if err == nil {
		ctx.render.JSON(w, http.StatusOK, user)
	} else {
		ctx.render.JSON(w, http.StatusNotFound, err)
	}
}

// CreateUserHandler adds a new user
func CreateUserHandler(w http.ResponseWriter, req *http.Request, ctx appContext) {
	decoder := json.NewDecoder(req.Body)
	var u User
	err := decoder.Decode(&u)
	if err != nil {
		ctx.render.JSON(w, http.StatusBadRequest, err)
	} else {
		user := User{-1, u.FirstName, u.LastName, u.DateOfBirth, u.LocationOfBirth}
		user, _ = db.Add(user)
		ctx.render.JSON(w, http.StatusCreated, user)
	}
}

// UpdateUserHandler updates a user object
func UpdateUserHandler(w http.ResponseWriter, req *http.Request, ctx appContext) {
	decoder := json.NewDecoder(req.Body)
	var u User
	err := decoder.Decode(&u)
	if err != nil {
		ctx.render.JSON(w, http.StatusBadRequest, err)
	} else {
		user := User{u.ID, u.FirstName, u.LastName, u.DateOfBirth, u.LocationOfBirth}
		user, err = db.Update(user)
		if err != nil {
			ctx.render.JSON(w, http.StatusOK, user)
		} else {
			ctx.render.JSON(w, http.StatusNotFound, err)
		}
	}
}

// DeleteUserHandler deletes a user
func DeleteUserHandler(w http.ResponseWriter, req *http.Request, ctx appContext) {
	vars := mux.Vars(req)
	uid, _ := strconv.Atoi(vars["uid"])
	ok, err := db.Delete(uid)
	if ok {
		// TO DO return empty body?
		ctx.render.Text(w, http.StatusNoContent, "")
	} else {
		ctx.render.JSON(w, http.StatusNotFound, err)
	}
}

// PassportsHandler not implemented yet
func PassportsHandler(w http.ResponseWriter, req *http.Request, ctx appContext) {
	log.Println("Handling Passports - Not implemented yet")
	ctx.render.Text(w, http.StatusNotImplemented, "")
}
