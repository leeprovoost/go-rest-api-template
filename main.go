package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/thoas/stats"
	"github.com/unrolled/render"
)

const local string = "LOCAL"

type appContext struct {
	metrics *stats.Stats
	render  *render.Render
	version string
	env     string
	port    string
	db      DataStorer
}

func main() {
	var (
		// environment variables
		env      = os.Getenv("ENV")      // LOCAL, DEV, STG, PRD
		port     = os.Getenv("PORT")     // server traffic on this port
		version  = os.Getenv("VERSION")  // path to VERSION file
		fixtures = os.Getenv("FIXTURES") // path to fixtures file
	)

	if env == "" || env == local {
		// running from localhost, so set some default values
		env = local
		port = "3001"
		version = "VERSION"
		fixtures = "fixtures.json"
	}

	// reading version from file
	dat, _ := ioutil.ReadFile(version)
	version = string(dat)
	version = strings.Trim(strings.Trim(version, "\n"), " ")

	// read JSON fixtures file
	var jsonObject map[string][]User
	log.Println("Location of fixtures.json file: " + fixtures)
	file, err := ioutil.ReadFile(fixtures)
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
	db := &MockDB{
		UserList:  list,
		MaxUserID: 1,
	}

	// initialse application context
	ctx := appContext{
		metrics: stats.New(),
		render:  render.New(),
		version: version,
		env:     env,
		port:    port,
		db:      db,
	}

	// start application
	StartServer(ctx)
}
