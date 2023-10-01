package main

import "fmt"

var buildVersion = defaultValue
var buildDate = defaultValue
var buildCommit = defaultValue

const defaultValue = "N/A"

func printBuildInfo() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
}
