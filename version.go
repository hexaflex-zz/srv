package main

import "fmt"

// Some version constants.
const (
	AppVendor  = "hexaflex"
	AppName    = "srv"
	AppVersion = "v0.0.1"
)

// Version returns the version string.
func Version() string {
	return fmt.Sprintf("%s %s %s", AppVendor, AppName, AppVersion)
}
