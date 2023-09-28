package cmd

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/NETWAYS/check_elasticsearch/internal/client"
	"github.com/NETWAYS/go-check"
	checkhttpconfig "github.com/NETWAYS/go-check-network/http/config"
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
	tlsConfig, err := checkhttpconfig.NewTLSConfig(&checkhttpconfig.TLSConfig{
		InsecureSkipVerify: c.Insecure,
		CAFile:             c.CAFile,
		KeyFile:            c.KeyFile,
		CertFile:           c.CertFile,
	})

	if err != nil {
		check.ExitError(err)
	}

	var rt http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     tlsConfig,
	}

	// Using a Bearer Token for authentication
	if c.Bearer != "" {
		rt = checkhttpconfig.NewAuthorizationCredentialsRoundTripper("Bearer", c.Bearer, rt)
	}

	// Using a BasicAuth for authentication
	if c.Username != "" {
		if c.Password == "" {
			check.ExitError(fmt.Errorf("specify the user name and password for server authentication"))
		}

		rt = checkhttpconfig.NewBasicAuthRoundTripper(c.Username, c.Password, rt)
	}

	return client.NewClient(u.String(), rt)
}
