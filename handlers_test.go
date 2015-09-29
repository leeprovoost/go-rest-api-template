package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thoas/stats"
	"github.com/unrolled/render"
)

func HandlersTestSetup() Env {
	env := Env{
		Metrics: stats.New(),
		Render:  render.New(),
	}
	return env
}

func teardown() {

}

func TestListUsersHandler(t *testing.T) {
	env := HandlersTestSetup()
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	makeHandler(env, ListUsersHandler).ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "they should be equal")
	assert.Equal(t, "application/json; charset=UTF-8", w.HeaderMap["Content-Type"][0], "they should be equal")
	// parse json body
	var f interface{}
	json.Unmarshal(w.Body.Bytes(), &f)
	m := f.(map[string]interface{})
	//fmt.Println(m)
	fmt.Println(m["users"])
}
