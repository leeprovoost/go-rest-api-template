package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// HandlerFunc is a custom implementation of the http.HandlerFunc
type HandlerFunc func(http.ResponseWriter, *http.Request, AppContext)

// makeHandler allows us to pass an environment struct to our handlers, without resorting to global
// variables. It accepts an environment (Env) struct and our own handler function. It returns
// a function of the type http.HandlerFunc so can be passed on to the HandlerFunc in main.go.
func makeHandler(ctx AppContext, fn func(http.ResponseWriter, *http.Request, AppContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, ctx)
	}
}

// HealthcheckHandler returns useful info about the app
func HealthcheckHandler(w http.ResponseWriter, req *http.Request, ctx AppContext) {
	check := Healthcheck{
		AppName: "go-rest-api-template",
		Version: ctx.Version,
	}
	ctx.Render.JSON(w, http.StatusOK, check)
}

// ListUsersHandler returns a list of users
func ListUsersHandler(w http.ResponseWriter, req *http.Request, ctx AppContext) {
	list, err := ctx.DB.ListUsers()
	if err != nil {
		response := Status{
			Status:  "404",
			Message: "can't find any users",
		}
		log.Println(err)
		ctx.Render.JSON(w, http.StatusNotFound, response)
		return
	}
	responseObject := make(map[string][]User)
	responseObject["users"] = list
	ctx.Render.JSON(w, http.StatusOK, responseObject)
}

// GetUserHandler returns a user object
func GetUserHandler(w http.ResponseWriter, req *http.Request, ctx AppContext) {
	vars := mux.Vars(req)
	uid, _ := strconv.Atoi(vars["uid"])
	user, err := ctx.DB.GetUser(uid)
	if err != nil {
		response := Status{
			Status:  "404",
			Message: "can't find user",
		}
		log.Println(err)
		ctx.Render.JSON(w, http.StatusNotFound, response)
		return
	}
	ctx.Render.JSON(w, http.StatusOK, user)
}

// CreateUserHandler adds a new user
func CreateUserHandler(w http.ResponseWriter, req *http.Request, ctx AppContext) {
	decoder := json.NewDecoder(req.Body)
	var u User
	err := decoder.Decode(&u)
	if err != nil {
		response := Status{
			Status:  "400",
			Message: "malformed user object",
		}
		log.Println(err)
		ctx.Render.JSON(w, http.StatusBadRequest, response)
		return
	}
	user := User{-1, u.FirstName, u.LastName, u.DateOfBirth, u.LocationOfBirth}
	user, _ = ctx.DB.AddUser(user)
	ctx.Render.JSON(w, http.StatusCreated, user)
}

// UpdateUserHandler updates a user object
func UpdateUserHandler(w http.ResponseWriter, req *http.Request, ctx AppContext) {
	decoder := json.NewDecoder(req.Body)
	var u User
	err := decoder.Decode(&u)
	if err != nil {
		response := Status{
			Status:  "400",
			Message: "malformed user object",
		}
		log.Println(err)
		ctx.Render.JSON(w, http.StatusBadRequest, response)
		return
	}
	user := User{
		ID:              u.ID,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		DateOfBirth:     u.DateOfBirth,
		LocationOfBirth: u.LocationOfBirth,
	}
	user, err = ctx.DB.UpdateUser(user)
	if err != nil {
		response := Status{
			Status:  "500",
			Message: "something went wrong",
		}
		log.Println(err)
		ctx.Render.JSON(w, http.StatusInternalServerError, response)
		return
	}
	ctx.Render.JSON(w, http.StatusOK, user)
}

// DeleteUserHandler deletes a user
func DeleteUserHandler(w http.ResponseWriter, req *http.Request, ctx AppContext) {
	vars := mux.Vars(req)
	uid, _ := strconv.Atoi(vars["uid"])
	err := ctx.DB.DeleteUser(uid)
	if err != nil {
		response := Status{
			Status:  "500",
			Message: "something went wrong",
		}
		log.Println(err)
		ctx.Render.JSON(w, http.StatusInternalServerError, response)
		return
	}
	ctx.Render.Text(w, http.StatusNoContent, "")
}

// PassportsHandler not implemented yet
func PassportsHandler(w http.ResponseWriter, req *http.Request, ctx AppContext) {
	log.Println("Handling Passports - Not implemented yet")
	ctx.Render.Text(w, http.StatusNotImplemented, "")
}
