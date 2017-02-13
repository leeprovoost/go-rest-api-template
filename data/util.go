package data

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/leeprovoost/go-rest-api-template/models"
	"github.com/palantir/stacktrace"
)

// CreateMockDatabase initialises a database for test purposes
func CreateMockDatabase() *MockDB {
	list := make(map[int]models.User)
	dt, _ := time.Parse(time.RFC3339, "1985-12-31T00:00:00Z")
	list[0] = models.User{
		ID:              0,
		FirstName:       "John",
		LastName:        "Doe",
		DateOfBirth:     dt,
		LocationOfBirth: "London",
	}
	dt, _ = time.Parse(time.RFC3339, "1992-01-01T00:00:00Z")
	list[1] = models.User{
		ID: 1, FirstName: "Jane",
		LastName:        "Doe",
		DateOfBirth:     dt,
		LocationOfBirth: "Milton Keynes",
	}
	return &MockDB{list, 1}
}

// LoadFixturesIntoMockDatabase loads data from fixtures file into MockDB
func LoadFixturesIntoMockDatabase(fixturesFile string) (*MockDB, error) {
	var jsonObject map[string][]models.User
	file, err := ioutil.ReadFile(fixturesFile)
	if err != nil {
		return nil, stacktrace.Propagate(err, "error reading fixtures file")
	}
	err = json.Unmarshal(file, &jsonObject)
	if err != nil {
		return nil, stacktrace.Propagate(err, "error parsing fixtures file")
	}
	list := make(map[int]models.User)
	list[0] = jsonObject["users"][0]
	list[1] = jsonObject["users"][1]
	return &MockDB{
		UserList:  list,
		MaxUserID: 1,
	}, nil
}
