package main

import "fmt"

var buildVersion string
var buildDate string
var buildCommit string

var defaultValue = "N/A"

func init() {
	if buildVersion == "" {
		buildVersion = defaultValue
	}
	if buildDate == "" {
		buildDate = defaultValue
	}
	if buildCommit == "" {
		buildCommit = defaultValue
	}

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
}
