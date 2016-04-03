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
	list, err := ctx.db.ListUsers()
	if err != nil {
		response := Status{
			Status:  "404",
			Message: "can't find any users",
		}
		log.Println(err)
		ctx.render.JSON(w, http.StatusNotFound, response)
		return
	}
	responseObject := make(map[string][]User)
	responseObject["users"] = list
	ctx.render.JSON(w, http.StatusOK, responseObject)
}

// GetUserHandler returns a user object
func GetUserHandler(w http.ResponseWriter, req *http.Request, ctx appContext) {
	vars := mux.Vars(req)
	uid, _ := strconv.Atoi(vars["uid"])
	user, err := ctx.db.GetUser(uid)
	if err != nil {
		response := Status{
			Status:  "404",
			Message: "can't find user",
		}
		log.Println(err)
		ctx.render.JSON(w, http.StatusNotFound, response)
		return
	}
	ctx.render.JSON(w, http.StatusOK, user)
}

// CreateUserHandler adds a new user
func CreateUserHandler(w http.ResponseWriter, req *http.Request, ctx appContext) {
	decoder := json.NewDecoder(req.Body)
	var u User
	err := decoder.Decode(&u)
	if err != nil {
		response := Status{
			Status:  "400",
			Message: "malformed user object",
		}
		log.Println(err)
		ctx.render.JSON(w, http.StatusBadRequest, response)
		return
	}
	user := User{-1, u.FirstName, u.LastName, u.DateOfBirth, u.LocationOfBirth}
	user, _ = ctx.db.AddUser(user)
	ctx.render.JSON(w, http.StatusCreated, user)
}

// UpdateUserHandler updates a user object
func UpdateUserHandler(w http.ResponseWriter, req *http.Request, ctx appContext) {
	decoder := json.NewDecoder(req.Body)
	var u User
	err := decoder.Decode(&u)
	if err != nil {
		response := Status{
			Status:  "400",
			Message: "malformed user object",
		}
		log.Println(err)
		ctx.render.JSON(w, http.StatusBadRequest, response)
		return
	}
	user := User{
		ID:              u.ID,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		DateOfBirth:     u.DateOfBirth,
		LocationOfBirth: u.LocationOfBirth,
	}
	user, err = ctx.db.UpdateUser(user)
	if err != nil {
		response := Status{
			Status:  "500",
			Message: "something went wrong",
		}
		log.Println(err)
		ctx.render.JSON(w, http.StatusInternalServerError, response)
		return
	}
	ctx.render.JSON(w, http.StatusOK, user)
}

// DeleteUserHandler deletes a user
func DeleteUserHandler(w http.ResponseWriter, req *http.Request, ctx appContext) {
	vars := mux.Vars(req)
	uid, _ := strconv.Atoi(vars["uid"])
	err := ctx.db.DeleteUser(uid)
	if err != nil {
		response := Status{
			Status:  "500",
			Message: "something went wrong",
		}
		log.Println(err)
		ctx.render.JSON(w, http.StatusInternalServerError, response)
		return
	}
	ctx.render.Text(w, http.StatusNoContent, "")
}

// PassportsHandler not implemented yet
func PassportsHandler(w http.ResponseWriter, req *http.Request, ctx appContext) {
	log.Println("Handling Passports - Not implemented yet")
	ctx.render.Text(w, http.StatusNotImplemented, "")
}
