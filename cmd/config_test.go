package cmd

import (
	"os"
	"testing"
)

func TestLoadFromEnv(t *testing.T) {
	c := Config{}

	err := os.Setenv("CHECK_ELASTICSEARCH_USERNAME", "foobar")
	defer os.Unsetenv("CHECK_ELASTICSEARCH_USERNAME") // to not impact other tests

	if err != nil {
		t.Error("Did not expect error, got: %w", err)
	}

	loadFromEnv(&c)

	if "foobar" != c.Username {
		t.Error("\nActual: ", c.Username, "\nExpected: ", "foobar")
	}
	if "" != c.Password {
		t.Error("\nActual: ", c.Password, "\nExpected: ", "empty-string")
	}
}
