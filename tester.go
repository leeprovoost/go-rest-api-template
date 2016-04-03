package main

import (
	"time"

	"github.com/thoas/stats"
	"github.com/unrolled/render"
)

func createContextForTestSetup() appContext {
	testVersion := "0.0.0"
	db := createMockDatabase()
	ctx := appContext{
		metrics: stats.New(),
		render:  render.New(),
		version: testVersion,
		env:     "LOCAL",
		port:    "3001",
		db:      db,
	}
	return ctx
}

func createMockDatabase() *MockDB {
	list := make(map[int]User)
	dt, _ := time.Parse(time.RFC3339, "1985-12-31T00:00:00Z")
	list[0] = User{0, "John", "Doe", dt, "London"}
	dt, _ = time.Parse(time.RFC3339, "1992-01-01T00:00:00Z")
	list[1] = User{1, "Jane", "Doe", dt, "Milton Keynes"}
	return &MockDB{list, 1}
}

func createTimestampForTest() time.Time {
	timestamp, _ := time.Parse(time.RFC3339, "2016-03-15T14:44:27Z")
	return timestamp
}
