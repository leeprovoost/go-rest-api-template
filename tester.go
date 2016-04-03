package main

import (
	"time"

	"github.com/thoas/stats"
	"github.com/unrolled/render"
)

func createContextForTestSetup() appContext {
	testVersion := "0.0.0"
	ctx := appContext{
		metrics: stats.New(),
		render:  render.New(),
		version: testVersion,
		env:     "LOCAL",
		port:    "3001",
	}
	return ctx
}

func createTimestampForTest() time.Time {
	timestamp, _ := time.Parse(time.RFC3339, "2016-03-15T14:44:27Z")
	return timestamp
}
