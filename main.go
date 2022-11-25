package main

import (
	"check_elasticsearch/cmd"
	"fmt"
)

var (
	version = "0.1.2"
	commit  = "c7e4d7b"
	date    = "25.04.2022"
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
