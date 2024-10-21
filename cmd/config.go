package cmd

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/NETWAYS/check_elasticsearch/internal/client"
	"github.com/NETWAYS/go-check"
	checkhttpconfig "github.com/NETWAYS/go-check-network/http/config"
)

type Config struct {
	Bearer   string // Currently unused in CLI
	Hostname string `env:"CHECK_ELASTICSEARCH_HOSTNAME"`
	CAFile   string `env:"CHECK_ELASTICSEARCH_CA_FILE"`
	CertFile string `env:"CHECK_ELASTICSEARCH_CERT_FILE"`
	KeyFile  string `env:"CHECK_ELASTICSEARCH_KEY_FILE"`
	Username string `env:"CHECK_ELASTICSEARCH_USERNAME"`
	Password string `env:"CHECK_ELASTICSEARCH_PASSWORD"`
	Port     int
	TLS      bool
	Insecure bool
}

// LoadFromEnv can be used to load struct values from 'env' tags.
// Mainly used to avoid passing secrets via the CLI
//
//	type Config struct {
//		Token    string `env:"BEARER_TOKEN"`
//	}
func loadFromEnv(config interface{}) {
	configValue := reflect.ValueOf(config).Elem()
	configType := configValue.Type()

	for i := range configValue.NumField() {
		field := configType.Field(i)
		tag := field.Tag.Get("env")

		// If there's no "env" tag, skip this field.
		if tag == "" {
			continue
		}

		envValue := os.Getenv(tag)

		if envValue == "" {
			continue
		}

		// Potential for addding different types
		// nolint: exhaustive, gocritic
		switch field.Type.Kind() {
		case reflect.String:
			configValue.Field(i).SetString(envValue)
		}
	}
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
			// nolint: perfsprint
			check.ExitError(fmt.Errorf("specify the user name and password for server authentication"))
		}

		rt = checkhttpconfig.NewBasicAuthRoundTripper(c.Username, c.Password, rt)
	}

	return client.NewClient(u.String(), rt)
}
