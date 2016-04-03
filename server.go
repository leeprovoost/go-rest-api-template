package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// StartServer Wraps the mux Router and uses the Negroni Middleware
func StartServer(ctx appContext) {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = makeHandler(ctx, route.HandlerFunc)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	n := negroni.Classic()
	n.Use(ctx.Metrics)
	n.UseHandler(router)
	log.Println("===> 🌍 Starting app (v" + ctx.version + ") on port " + ctx.port + " in " + ctx.env + " mode.")
	if ctx.env == local {
		n.Run("localhost:" + ctx.port)
	} else {
		n.Run(":" + ctx.port)
	}
}
