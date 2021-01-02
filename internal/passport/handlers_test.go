package passport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
)

func TestHealthcheckHandler(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	r, _ := http.NewRequest(http.MethodGet, "/healthcheck", nil)
	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	// recreate route from routes.go
	router.
		Methods(http.MethodGet).
		Path("/healthcheck").
		Name("HealthcheckHandler").
		Handler(MakeHandler(appEnv, HealthcheckHandler))
	n := negroni.New()
	n.UseHandler(router)
	n.ServeHTTP(w, r)
	// test response headers and codes
	assert.Equal(t, http.StatusOK, w.Code, "they should be equal")
	assert.Equal(t, "GNU Terry Pratchett", w.HeaderMap["X-Clacks-Overhead"][0], "they should be equal")
	// parse json body
	var f interface{}
	json.Unmarshal(w.Body.Bytes(), &f)
	m := f.(map[string]interface{})
	assert.Equal(t, "go-rest-api-template", m["appName"], "they should be equal")
	assert.Equal(t, appEnv.Version, m["version"], "they should be equal")
}

func TestListUsersHandler(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	r, _ := http.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	// recreate route from routes.go
	router.
		Methods(http.MethodGet).
		Path("/users").
		Name("ListUsersHandler").
		Handler(MakeHandler(appEnv, ListUsersHandler))
	n := negroni.New()
	n.UseHandler(router)
	n.ServeHTTP(w, r)
	// test response headers and codes
	assert.Equal(t, http.StatusOK, w.Code, "they should be equal")
	assert.Equal(t, "GNU Terry Pratchett", w.HeaderMap["X-Clacks-Overhead"][0], "they should be equal")
	// parse json body
	var f interface{}
	json.Unmarshal(w.Body.Bytes(), &f)
}
