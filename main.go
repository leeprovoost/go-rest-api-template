package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/thoas/stats"
	"github.com/unrolled/render"
)

const local string = "LOCAL"

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
	version, err := ParseVersionFile(version)
	if err != nil {
		log.Fatal(err)
	}
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
	ctx := AppContext{
		Metrics: stats.New(),
		Render:  render.New(),
		Version: version,
		Env:     env,
		Port:    port,
		DB:      db,
	}
	// start application
	StartServer(ctx)
}
