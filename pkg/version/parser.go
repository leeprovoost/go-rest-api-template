package version

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ParseVersionFile returns the version as a string, parsing and validating a file given the path.
func ParseVersionFile(versionPath string) (string, error) {
	dat, err := os.ReadFile(versionPath)
	if err != nil {
		return "", fmt.Errorf("reading version file: %w", err)
	}
	version := strings.TrimSpace(string(dat))
	// regex pulled from https://github.com/sindresorhus/semver-regex
	semverRegex := `^v?(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)(?:-[\da-z\-]+(?:\.[\da-z\-]+)*)?(?:\+[\da-z\-]+(?:\.[\da-z\-]+)*)?$`
	match, err := regexp.MatchString(semverRegex, version)
	if err != nil {
		return "", fmt.Errorf("executing regex match: %w", err)
	}
	if !match {
		return "", fmt.Errorf("string in VERSION is not a valid version number: %q", version)
	}
	return version, nil
}
