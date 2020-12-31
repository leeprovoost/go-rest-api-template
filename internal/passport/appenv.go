package passport

import (
	"github.com/leeprovoost/go-rest-api-template/internal/passport/models"
	"github.com/unrolled/render"
)

// AppEnv holds application configuration data
type AppEnv struct {
	Render    *render.Render
	Version   string
	Env       string
	Port      string
	UserStore models.UserStorage
}

// CreateContextForTestSetup initialises an application context struct
// for testing purposes
func CreateContextForTestSetup() AppEnv {
	testVersion := "0.0.0"
	appEnv := AppEnv{
		Render:    render.New(),
		Version:   testVersion,
		Env:       "LOCAL",
		Port:      "3001",
		UserStore: NewUserService(CreateMockDataSet()),
	}
	return appEnv
}
