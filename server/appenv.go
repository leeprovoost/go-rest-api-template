package server

import (
	"github.com/leeprovoost/go-rest-api-template/data"
	"github.com/unrolled/render"
)

// AppEnv holds application configuration data
type AppEnv struct {
	Render  *render.Render
	Version string
	Env     string
	Port    string
	DB      data.DataStorer
}
