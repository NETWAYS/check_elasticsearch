package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	es "github.com/NETWAYS/check_elasticsearch/internal/elasticsearch"
)

type Client struct {
	Client http.Client
	URLs   []*url.URL
}

func NewClient(urls []*url.URL, rt http.RoundTripper) *Client {
	c := &http.Client{
		Transport: rt,
	}

	return &Client{
		URLs:   urls,
		Client: *c,
	}
}

func (c *Client) Perform(req *http.Request) (*http.Response, error) {
	originalPath := req.URL.String()

	for _, hostURL := range c.URLs {
		// For each URL take the request, prepend the URL
		u, _ := url.JoinPath(hostURL.String(), originalPath)

		req.URL, _ = url.Parse(u)

		resp, errDo := c.Client.Do(req)

		if errDo != nil {
			// If there's an error we try the next host
			continue
		}

		return resp, errDo
	}

	return &http.Response{}, errors.New("no node reachable")
}

func (c *Client) Health() (*es.HealthResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	u := "/_cluster/health"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)

	r := &es.HealthResponse{}

	if err != nil {
		return r, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.Perform(req)

	if err != nil {
		return r, fmt.Errorf("could not fetch cluster health: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return r, fmt.Errorf("request failed for cluster health: %s", resp.Status)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(r)

	if err != nil {
		return r, fmt.Errorf("error parsing the response body: %w", err)
	}

	return r, nil
}

// SearchMessages runs a query_string query and returns the
// count of documents and the requesed values via messageKey
func (c *Client) SearchMessages(index string, query string, messageKey string) (uint, []string, error) {
	queryBody := es.SearchRequest{
		Query: es.Query{
			QueryString: &es.QueryString{
				Query: query,
			},
		},
	}

	var total uint

	var messages []string

	data, err := json.Marshal(queryBody)
	body := bytes.NewReader(data)

	if err != nil {
		return total, messages, fmt.Errorf("error encoding query: %w", err)
	}

	u := index + "/_search"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, body)

	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return total, messages, fmt.Errorf("error creating request: %w", err)
	}

	p := req.URL.Query()
	p.Add("track_total_hits", "true")
	p.Add("size", "1")

	req.URL.RawQuery = p.Encode()

	resp, err := c.Perform(req)

	if err != nil {
		return total, messages, fmt.Errorf("could not execute search request: %s", err.Error())
	}

	var response es.SearchResponse

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return total, messages, fmt.Errorf("error parsing the response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		queryErrors := response.GetErrors()

		return total, messages, fmt.Errorf("failed to run query: %s", queryErrors)
	}

	total = response.Hits.Total.Value

	for _, hit := range response.Hits.Hits {
		// When the user does not request a field we skip here
		if messageKey == "" {
			continue
		}

		// Append the requested values if a key is given
		if value, ok := hit.Source[messageKey]; ok {
			messages = append(messages, fmt.Sprint(value))
		} else {
			return total, messages, fmt.Errorf("document does not contain key '%s': %s", messageKey, hit.ID)
		}
	}

	return total, messages, nil
}

func (c *Client) NodeStats() (*es.ClusterStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	u := "/_nodes/stats"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)

	r := &es.ClusterStats{}

	if err != nil {
		return r, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.Perform(req)

	if err != nil {
		return r, fmt.Errorf("could not fetch cluster nodes statistics: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return r, fmt.Errorf("request failed for cluster nodes statistics: %s", resp.Status)
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(r)

	if err != nil {
		return r, fmt.Errorf("error parsing the response body: %w", err)
	}

	return r, nil
}

func (c *Client) Snapshot(repository string, snapshot string) (*es.SnapshotResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r := &es.SnapshotResponse{}

	u, _ := url.JoinPath("/_snapshot/", repository, snapshot)

	// Retrieve snapshots in descending order to get latest
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u+"?order=desc", nil)

	if err != nil {
		return r, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.Perform(req)

	if err != nil {
		return r, fmt.Errorf("could not fetch snapshots: %s", err.Error())
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
