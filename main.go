package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/thoas/stats"
	"github.com/unrolled/render"
)

type appContext struct {
	Metrics *stats.Stats
	Render  *render.Render
}

var fPort string
var fFixtures string

func init() {
	// parse command line flags
	flag.StringVar(&fFixtures, "fixtures", "./fixtures.json", "location of fixtures.json file")
	flag.StringVar(&fPort, "port", "3001", "serve traffic on this port")
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
	ctx := appContext{
		Metrics: stats.New(),
		Render:  render.New(),
	}

	fmt.Println("===> 🌍 Starting server on port: " + fPort)
	StartServer(ctx, fPort)
}
