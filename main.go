package main

import (
	"check_elasticsearch/cmd"
	"fmt"
)

var (
	version = "0.2.0"
	commit  = "f9eefa1"
	date    = "25.11.2022"
)

func main() {
	cmd.Execute(buildVersion())
}

//goland:noinspection GoBoolExpressions
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
