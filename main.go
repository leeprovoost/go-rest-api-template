package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/thoas/stats"
	"github.com/unrolled/render"
)

type env struct {
	Metrics *stats.Stats
	Render  *render.Render
}

var fPort string
var fFixtures string

func init() {
	// parse command line flags
	flag.StringVar(&fFixtures, "fixtures", "./fixtures.json", "location of fixtures.json file")
	flag.StringVar(&fPort, "port", "3009", "serve traffic on this port")
	flag.Parse()

	// read JSON fixtures file
	var jsonObject map[string][]User
	fmt.Println("Location of fixtures.json file: " + fFixtures)
	file, err := ioutil.ReadFile(fFixtures)
	if err != nil {
		log.Fatalf("File error: %v\n", err)
	}
	err = json.Unmarshal(file, &jsonObject)
	if err != nil {
		log.Fatal(err)
	}
	// load data in database
	list := make(map[int]User)
	list[0] = jsonObject["users"][0]
	list[1] = jsonObject["users"][1]
	db = &Database{
		UserList:  list,
		MaxUserID: 1,
	}
}

func main() {
	env := env{
		Metrics: stats.New(),
		Render:  render.New(),
	}
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		var handler http.Handler

		handler = makeHandler(env, route.Handler)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	n := negroni.Classic()
	n.Use(env.Metrics)
	n.UseHandler(router)
	fmt.Println("===> 🌍 Starting server on port: " + fPort)
	n.Run(":" + fPort)
}
