package main

import (
	"fmt"
	"regexp"
	"time"
)

var (
	builtAt string
	version string
)

func init() {
	if builtAt == "" {
		builtAt = time.Now().Format(time.RFC3339)
	}

	if version == "" {
		version = "dev"
	}
}

func getBuildInfo() string {
	friendlyVersion := version
	isHash, _ := regexp.MatchString("^[a-f0-9]+$", friendlyVersion)

	if friendlyVersion == "dev" {
		friendlyVersion = "development build"
	} else if isHash {
		friendlyVersion = "rev. " + friendlyVersion[:7]
	} else {
		friendlyVersion = "ver. " + friendlyVersion
	}

	return fmt.Sprintf("%s, built at %s", friendlyVersion, builtAt)
}
