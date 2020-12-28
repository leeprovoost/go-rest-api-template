package db

import (
	"time"

	"github.com/leeprovoost/go-rest-api-template/models"
)

// CreateMockDataSet initialises a database for test purposes
func CreateMockDataSet() map[int]models.User {
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
	return list
}
