package main

import (
	"io/ioutil"
	"regexp"
	"strings"
	"time"

	"github.com/palantir/stacktrace"
	"github.com/thoas/stats"
	"github.com/unrolled/render"
)

// AppContext holds application configuration data
type AppContext struct {
	Metrics *stats.Stats
	Render  *render.Render
	Version string
	Env     string
	Port    string
	DB      DataStorer
}

// Healthcheck will store information about its name and version
type Healthcheck struct {
	AppName string `json:"appName"`
	Version string `json:"version"`
}

// Status is a custom response object we pass around the system and send back to the customer
// 404: Not found
// 500: Internal Server Error
type Status struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// CreateContextForTestSetup initialises an application context struct
// for testing purposes
func CreateContextForTestSetup() AppContext {
	testVersion := "0.0.0"
	db := CreateMockDatabase()
	ctx := AppContext{
		Metrics: stats.New(),
		Render:  render.New(),
		Version: testVersion,
		Env:     local,
		Port:    "3001",
		DB:      db,
	}
	return ctx
}

// CreateMockDatabase initialises a database for test purposes
func CreateMockDatabase() *MockDB {
	list := make(map[int]User)
	dt, _ := time.Parse(time.RFC3339, "1985-12-31T00:00:00Z")
	list[0] = User{0, "John", "Doe", dt, "London"}
	dt, _ = time.Parse(time.RFC3339, "1992-01-01T00:00:00Z")
	list[1] = User{1, "Jane", "Doe", dt, "Milton Keynes"}
	return &MockDB{list, 1}
}

// ParseVersionFile returns the version as a string, parsing and validating a file given the path
func ParseVersionFile(versionPath string) (string, error) {
	dat, err := ioutil.ReadFile(versionPath)
	if err != nil {
		return "", stacktrace.Propagate(err, "error reading version file")
	}
	version := string(dat)
	version = strings.Trim(strings.Trim(version, "\n"), " ")
	// regex pulled from official https://github.com/sindresorhus/semver-regex
	semverRegex := `^v?(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)(?:-[\da-z\-]+(?:\.[\da-z\-]+)*)?(?:\+[\da-z\-]+(?:\.[\da-z\-]+)*)?$`
	match, err := regexp.MatchString(semverRegex, version)
	if err != nil {
		return "", stacktrace.Propagate(err, "error executing regex match")
	}
	if !match {
		return "", stacktrace.NewError("string in VERSION is not a valid version number")
	}
	return version, nil
}
