package client

import (
	"crypto/tls"
	"fmt"
	es7 "github.com/elastic/go-elasticsearch/v7"
	"net/http"
)

type Client struct {
	Url      string
	Username string
	Password string
	Insecure bool
	Client   *es7.Client
	Version  string
}

func NewClient(url, username, password string) *Client {
	return &Client{
		Url:      url,
		Username: username,
		Password: password,
		Insecure: false,
	}
}

func (c *Client) Connect() error {
	cfg := es7.Config{
		Addresses: []string{c.Url},
		Username:  c.Username,
		Password:  c.Password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: c.Insecure,
				MinVersion:         tls.VersionTLS11,
			},
		},
	}

	esClient, err := es7.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("could not connect to cluster: %w", err)
	}

	c.Client = esClient

	info, err := c.Info()
	if err != nil {
		return err
	}

	c.Version = info.Version.Number

	return nil
}
