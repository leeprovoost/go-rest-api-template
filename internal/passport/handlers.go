package passport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	models "github.com/leeprovoost/go-rest-api-template/internal/passport/models"
	"github.com/leeprovoost/go-rest-api-template/pkg/health"
	"github.com/leeprovoost/go-rest-api-template/pkg/status"
	log "github.com/sirupsen/logrus"
)

// HandlerFunc is a custom implementation of the http.HandlerFunc
type HandlerFunc func(http.ResponseWriter, *http.Request, AppEnv)

// MakeHandler allows us to pass an environment struct to our handlers, without resorting to global
// variables. It accepts an environment (Env) struct and our own handler function. It returns
// a function of the type http.HandlerFunc so can be passed on to the HandlerFunc in main.go.
func MakeHandler(appEnv AppEnv, fn func(http.ResponseWriter, *http.Request, AppEnv)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Terry Pratchett tribute
		w.Header().Set("X-Clacks-Overhead", "GNU Terry Pratchett")
		// return function with AppEnv
		fn(w, r, appEnv)
	}
}

// HealthcheckHandler returns useful info about the app
func HealthcheckHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	check := health.Check{
		AppName: "go-rest-api-template",
		Version: appEnv.Version,
	}
	appEnv.Render.JSON(w, http.StatusOK, check)
}

// ListUsersHandler returns a list of users
func ListUsersHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	list, err := appEnv.UserStore.ListUsers()
	if err != nil {
		response := status.Response{
			Status:  strconv.Itoa(http.StatusNotFound),
			Message: "can't find any users",
		}
		log.WithFields(log.Fields{
			"env":    appEnv.Env,
			"status": http.StatusNotFound,
		}).Error("Can't find any users")
		appEnv.Render.JSON(w, http.StatusNotFound, response)
		return
	}
	responseObject := make(map[string]interface{})
	responseObject["users"] = list
	responseObject["count"] = len(list)
	appEnv.Render.JSON(w, http.StatusOK, responseObject)
}

// GetUserHandler returns a user object
func GetUserHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	vars := mux.Vars(req)
	uid, _ := strconv.Atoi(vars["uid"])
	user, err := appEnv.UserStore.GetUser(uid)
	if err != nil {
		response := status.Response{
			Status:  strconv.Itoa(http.StatusNotFound),
			Message: "can't find user",
		}
		log.WithFields(log.Fields{
			"env":    appEnv.Env,
			"status": http.StatusNotFound,
		}).Error("Can't find user")
		appEnv.Render.JSON(w, http.StatusNotFound, response)
		return
	}
	appEnv.Render.JSON(w, http.StatusOK, user)
}

// CreateUserHandler adds a new user
func CreateUserHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	decoder := json.NewDecoder(req.Body)
	var u models.User
	err := decoder.Decode(&u)
	if err != nil {
		response := status.Response{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: "malformed user object",
		}
		log.WithFields(log.Fields{
			"env":    appEnv.Env,
			"status": http.StatusBadRequest,
		}).Error("malformed user object")
		appEnv.Render.JSON(w, http.StatusBadRequest, response)
		return
	}
	user := models.User{
		ID:              -1,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		DateOfBirth:     u.DateOfBirth,
		LocationOfBirth: u.LocationOfBirth,
	}
	user, _ = appEnv.UserStore.AddUser(user)
	appEnv.Render.JSON(w, http.StatusCreated, user)
}

// UpdateUserHandler updates a user object
func UpdateUserHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	decoder := json.NewDecoder(req.Body)
	var u models.User
	err := decoder.Decode(&u)
	if err != nil {
		response := status.Response{
			Status:  strconv.Itoa(http.StatusBadRequest),
			Message: "malformed user object",
		}
		log.WithFields(log.Fields{
			"env":    appEnv.Env,
			"status": http.StatusBadRequest,
		}).Error("malformed user object")
		appEnv.Render.JSON(w, http.StatusBadRequest, response)
		return
	}
	user := models.User{
		ID:              u.ID,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		DateOfBirth:     u.DateOfBirth,
		LocationOfBirth: u.LocationOfBirth,
	}
	user, err = appEnv.UserStore.UpdateUser(user)
	if err != nil {
		response := status.Response{
			Status:  strconv.Itoa(http.StatusInternalServerError),
			Message: "something went wrong",
		}
		log.WithFields(log.Fields{
			"env":    appEnv.Env,
			"status": http.StatusInternalServerError,
		}).Error("something went wrong")
		appEnv.Render.JSON(w, http.StatusInternalServerError, response)
		return
	}
	appEnv.Render.JSON(w, http.StatusOK, user)
}

// DeleteUserHandler deletes a user
func DeleteUserHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	vars := mux.Vars(req)
	uid, _ := strconv.Atoi(vars["uid"])
	err := appEnv.UserStore.DeleteUser(uid)
	if err != nil {
		response := status.Response{
			Status:  strconv.Itoa(http.StatusInternalServerError),
			Message: "something went wrong",
		}
		log.WithFields(log.Fields{
			"env":    appEnv.Env,
			"status": http.StatusInternalServerError,
		}).Error("something went wrong")
		appEnv.Render.JSON(w, http.StatusInternalServerError, response)
		return
	}
	appEnv.Render.Text(w, http.StatusNoContent, "")
}

// PassportsHandler not implemented yet
func PassportsHandler(w http.ResponseWriter, req *http.Request, appEnv AppEnv) {
	log.WithFields(log.Fields{
		"env":    appEnv.Env,
		"status": http.StatusInternalServerError,
	}).Error("Handling Passports - Not implemented yet")
	appEnv.Render.Text(w, http.StatusNotImplemented, "")
}
