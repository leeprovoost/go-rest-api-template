package passport

import "net/http"

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthcheck", s.handleHealthcheck)
	mux.HandleFunc("GET /ready", s.handleReady)

	// Users
	mux.HandleFunc("GET /users", s.handleListUsers)
	mux.HandleFunc("GET /users/{id}", s.handleGetUser)
	mux.HandleFunc("POST /users", s.handleCreateUser)
	mux.HandleFunc("PUT /users/{id}", s.handleUpdateUser)
	mux.HandleFunc("DELETE /users/{id}", s.handleDeleteUser)

	// Passports
	mux.HandleFunc("GET /users/{uid}/passports", s.handleListUserPassports)
	mux.HandleFunc("GET /passports/{id}", s.handleGetPassport)
	mux.HandleFunc("POST /users/{uid}/passports", s.handleCreatePassport)
	mux.HandleFunc("PUT /passports/{id}", s.handleUpdatePassport)
	mux.HandleFunc("DELETE /passports/{id}", s.handleDeletePassport)

	return mux
}
