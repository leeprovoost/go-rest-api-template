package main

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// StartServer Wraps the mux Router and uses the Negroni Middleware
func StartServer(env env, port string) {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = makeHandler(env, route.HandlerFunc)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	n := negroni.Classic()
	n.Use(env.Metrics)
	n.UseHandler(router)

	n.Run(":" + port)
}
