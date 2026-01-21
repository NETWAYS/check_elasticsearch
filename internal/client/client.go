package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	es "github.com/NETWAYS/check_elasticsearch/internal/elasticsearch"
)

type Client struct {
	Client http.Client
	URL    string
}

func NewClient(url string, rt http.RoundTripper) *Client {
	// Small wrapper
	c := &http.Client{
		Transport: rt,
	}

	return &Client{
		URL:    url,
		Client: *c,
	}
}

func (c *Client) Health() (r *es.HealthResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	u, _ := url.JoinPath(c.URL, "/_cluster/health")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)

	if err != nil {
		err = fmt.Errorf("error creating request: %w", err)
		return
	}

	resp, err := c.Client.Do(req)

	if err != nil {
		err = fmt.Errorf("could not fetch cluster health: %w", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("request failed for cluster health: %s", resp.Status)
		return
	}

	r = &es.HealthResponse{}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(r)

	if err != nil {
		err = fmt.Errorf("could not decode health json: %w", err)
		return
	}

	return
}

func (c *Client) SearchMessages(index string, query string, messageKey string) (total uint, messages []string, err error) {
	queryBody := es.SearchRequest{
		Query: es.Query{
			QueryString: &es.QueryString{
				Query: query,
			},
		},
	}

	data, err := json.Marshal(queryBody)
	body := bytes.NewReader(data)

	if err != nil {
		err = fmt.Errorf("error encoding query: %w", err)
		return
	}

	u, _ := url.JoinPath(c.URL, index, "/_search")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, body)

	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		err = fmt.Errorf("error creating request: %w", err)
		return
	}

	p := req.URL.Query()
	p.Add("track_total_hits", "true")
	p.Add("size", "1")

	req.URL.RawQuery = p.Encode()

	resp, err := c.Client.Do(req)

	if err != nil {
		err = fmt.Errorf("could not execute search request: %w", err)
		return
	}

	var response es.SearchResponse

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		err = fmt.Errorf("error parsing the response body: %w", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		queryErrors := response.GetErrors()
		err = fmt.Errorf("failed to run query: %s", queryErrors)

		return
	}

	total = response.Hits.Total.Value

	for _, hit := range response.Hits.Hits {
		if value, ok := hit.Source[messageKey]; ok {
			messages = append(messages, fmt.Sprint(value))
		} else {
			err = fmt.Errorf("message does not contain key '%s': %s", messageKey, hit.ID)
			return
		}
	}

	return
}

func (c *Client) NodeStats() (r *es.ClusterStats, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	u, _ := url.JoinPath(c.URL, "/_nodes/stats")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)

	if err != nil {
		err = fmt.Errorf("error creating request: %w", err)
		return
	}

	resp, err := c.Client.Do(req)

	if err != nil {
		err = fmt.Errorf("could not fetch cluster nodes statistics: %w", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("request failed for cluster nodes statistics: %s", resp.Status)
		return
	}

	r = &es.ClusterStats{}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(r)

	if err != nil {
		err = fmt.Errorf("could not decode nodes statistics json: %w", err)
		return
	}

	return
}

func (c *Client) Snapshot(repository string, snapshot string) (*es.SnapshotResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r := &es.SnapshotResponse{}

	u, _ := url.JoinPath(c.URL, "/_snapshot/", repository, snapshot)

	// Retrieve snapshots in descending order to get latest
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u+"?order=desc", nil)

	if err != nil {
		return r, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.Client.Do(req)

	if err != nil {
		return r, fmt.Errorf("could not fetch snapshots: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return r, fmt.Errorf("request failed for snapshots: %s", resp.Status)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(r)

	if err != nil {
		return r, fmt.Errorf("could not decode snapshot response: %w", err)
	}

	return r, nil
}
