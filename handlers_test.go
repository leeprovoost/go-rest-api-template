package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
	fmt.Println(w)
	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}
	fmt.Println(w.Code)
	fmt.Println(w.HeaderMap)
	fmt.Println(w.Body)
}
