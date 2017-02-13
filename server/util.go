package server

import (
	"github.com/leeprovoost/go-rest-api-template/data"
	"github.com/unrolled/render"
)

// CreateContextForTestSetup initialises an application context struct
// for testing purposes
func CreateContextForTestSetup() AppEnv {
	testVersion := "0.0.0"
	db := data.CreateMockDatabase()
	appEnv := AppEnv{
		Render:  render.New(),
		Version: testVersion,
		Env:     "LOCAL",
		Port:    "3001",
		DB:      db,
	}
	return appEnv
}
