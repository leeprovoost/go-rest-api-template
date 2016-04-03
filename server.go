package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/secure"
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
	// security
	secureMiddleware := secure.New(secure.Options{
		AllowedHosts: []string{}, // AllowedHosts is a list of fully qualified domain names that are allowed (CORS)
	})
	// start now
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.Use(ctx.metrics)
	n.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))
	n.UseHandler(router)
	log.Println("===> 🌍 Starting app (v" + ctx.version + ") on port " + ctx.port + " in " + ctx.env + " mode.")
	if ctx.env == local {
		n.Run("localhost:" + ctx.port)
	} else {
		n.Run(":" + ctx.port)
	}
}
