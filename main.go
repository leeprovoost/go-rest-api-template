package main

import (
	"log"
	"os"

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
	// load fixtures data into mock database
	db, err := LoadFixturesIntoMockDatabase(fixtures)
	if err != nil {
		log.Fatal(err)
	}
	// initialse application context
	ctx := AppContext{
		Render:  render.New(),
		Version: version,
		Env:     env,
		Port:    port,
		DB:      db,
	}
	// start application
	StartServer(ctx)
}
