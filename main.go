package main

import (
	"fmt"

	"github.com/NETWAYS/check_elasticsearch/cmd"
)

var (
	// These get filled at build time with the proper vaules.
	version = "development"
	commit  = "HEAD"
	date    = "latest"
)

func main() {
	cmd.Execute(buildVersion())
}

func buildVersion() string {
	result := version

	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}

	if date != "" {
		result = fmt.Sprintf("%s\ndate: %s", result, date)
	}

	return result
}
