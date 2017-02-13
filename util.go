package main

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/palantir/stacktrace"
)

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
