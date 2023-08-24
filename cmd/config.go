package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"check_elasticsearch/internal/client"
	"check_elasticsearch/internal/config"

	"github.com/NETWAYS/go-check"
)

type Config struct {
	Hostname  string
	BasicAuth string
	Bearer    string
	CAFile    string
	CertFile  string
	KeyFile   string
	Username  string
	Password  string
	Port      int
	TLS       bool
	Insecure  bool
}

var cliConfig Config

func (c *Config) NewClient() *client.Client {
	u := url.URL{
		Scheme: "http",
		Host:   c.Hostname + ":" + strconv.Itoa(c.Port),
	}

	if c.TLS {
		u.Scheme = "https"
	}

	// Create TLS configuration for default RoundTripper
	tlsConfig, err := config.NewTLSConfig(&config.TLSConfig{
		InsecureSkipVerify: c.Insecure,
		CAFile:             c.CAFile,
		KeyFile:            c.KeyFile,
		CertFile:           c.CertFile,
	})

	if err != nil {
		check.ExitError(err)
	}

	var rt http.RoundTripper = &http.Transport{
		TLSClientConfig:       tlsConfig,
		IdleConnTimeout:       10 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	}

	// Using a Bearer Token for authentication
	if c.Bearer != "" {
		var t = config.Secret(c.Bearer)
		rt = config.NewAuthorizationCredentialsRoundTripper("Bearer", t, rt)
	}

	// Using a BasicAuth for authentication
	if c.Username != "" {
		if c.Password == "" {
			check.ExitError(fmt.Errorf("specify the user name and password for server authentication"))
		}

		var p = config.Secret(c.Password)

		rt = config.NewBasicAuthRoundTripper(c.Username, p, "", rt)
	}

	return client.NewClient(u.String(), rt)
}
