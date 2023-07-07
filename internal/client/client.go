package client

import (
	"bytes"
	es "check_elasticsearch/internal/elasticsearch"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	Client http.Client
	Url    string
}

func NewClient(url string, rt http.RoundTripper) *Client {
	// Small wrapper
	c := &http.Client{
		Transport: rt,
	}

	return &Client{
		Url:    url,
		Client: *c,
	}
}

func (c *Client) Health() (r *es.HealthResponse, err error) {
	u, _ := url.JoinPath(c.Url, "/_cluster/health")
	resp, err := c.Client.Get(u)

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

	u, _ := url.JoinPath(c.Url, index, "/_search")

	req, err := http.NewRequest("GET", u, body)

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

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("request failed for search: %s", resp.Status)
		return
	}

	var response es.SearchResponse

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		err = fmt.Errorf("error parsing the response body: %w", err)
		return
	}

	total = response.Hits.Total.Value

	for _, hit := range response.Hits.Hits {
		if value, ok := hit.Source[messageKey]; ok {
			messages = append(messages, fmt.Sprint(value))
		} else {
			err = fmt.Errorf("message does not contain key '%s': %s", messageKey, hit.Id)
			return
		}
	}

	return
}

func (c *Client) NodeStats() (r *es.ClusterStats, err error) {
	u, _ := url.JoinPath(c.Url, "/_nodes/stats")
	resp, err := c.Client.Get(u)

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
