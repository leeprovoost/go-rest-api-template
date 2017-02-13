package main

import (
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/leeprovoost/go-rest-api-template/data"
	"github.com/leeprovoost/go-rest-api-template/server"
	"github.com/unrolled/render"
)

func init() {
	if "LOCAL" == strings.ToUpper(os.Getenv("ENV")) {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	// ===========================================================================
	// Load environment variables
	// ===========================================================================
	var (
		env      = strings.ToUpper(os.Getenv("ENV")) // LOCAL, DEV, STG, PRD
		port     = os.Getenv("PORT")                 // server traffic on this port
		version  = os.Getenv("VERSION")              // path to VERSION file
		fixtures = os.Getenv("FIXTURES")             // path to fixtures file
	)
	// ===========================================================================
	// Read version information
	// ===========================================================================
	version, err := ParseVersionFile(version)
	if err != nil {
		log.WithFields(log.Fields{
			"env":  env,
			"err":  err,
			"path": os.Getenv("VERSION"),
		}).Fatal("Can't find a VERSION file")
		return
	}
	log.WithFields(log.Fields{
		"env":     env,
		"path":    os.Getenv("VERSION"),
		"version": version,
	}).Info("Loaded VERSION file")
	// ===========================================================================
	//  Load fixtures data into mock database
	// ===========================================================================
	db, err := data.LoadFixturesIntoMockDatabase(fixtures)
	if err != nil {
		log.WithFields(log.Fields{
			"env":      env,
			"err":      err,
			"fixtures": fixtures,
		}).Fatal("Can't find a fixtures.json file")
		return
	}
	// ===========================================================================
	// Initialise application context
	// ===========================================================================
	appEnv := server.AppEnv{
		Render:  render.New(),
		Version: version,
		Env:     env,
		Port:    port,
		DB:      db,
	}
	// ===========================================================================
	// Start application
	// ===========================================================================
	server.StartServer(appEnv)
}
