package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thoas/stats"
	"github.com/unrolled/render"
)

func TestListUsersHandler(t *testing.T) {
	env := Env{
		Metrics: stats.New(),
		Render:  render.New(),
	}
	req, _ := http.NewRequest("GET", "http://localhost:3009/users", nil)
	w := httptest.NewRecorder()
	makeHandler(env, ListUsersHandler).ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}
}
