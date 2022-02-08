package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) Info() (i *InfoResponse, err error) {
	info, err := c.Client.Info()
	if err != nil {
		err = fmt.Errorf("could not fetch cluster info: %w", err)
		return
	}

	if info.IsError() {
		err = fmt.Errorf("request failed for cluster info: %s", info.Status())
		return
	}

	i = &InfoResponse{}
	err = json.NewDecoder(info.Body).Decode(i)

	if err != nil {
		err = fmt.Errorf("could not decode info json: %w", err)
		return
	}

	return
}

func (c *Client) Health() (r *HealthResponse, err error) {
	health, err := c.Client.Cluster.Health()

	if err != nil {
		err = fmt.Errorf("could not fetch cluster health: %w", err)
		return
	}

	if health.IsError() {
		err = fmt.Errorf("request failed for cluster health: %s", health.Status())
		return
	}

	r = &HealthResponse{}
	err = json.NewDecoder(health.Body).Decode(r)

	if err != nil {
		err = fmt.Errorf("could not decode health json: %w", err)
		return
	}

	return
}

func (c *Client) SearchMessages(index string, query string, messageKey string) (total uint, messages []string, err error) {
	var body bytes.Buffer

	queryBody := SearchRequest{
		Query: Query{
			QueryString: &QueryString{
				Query: query,
			},
		},
	}

	err = json.NewEncoder(&body).Encode(queryBody)

	if err != nil {
		err = fmt.Errorf("error encoding query: %w", err)
		return
	}

	s := c.Client.Search
	res, err := c.Client.Search(
		s.WithContext(context.Background()),
		s.WithIndex(index),
		s.WithBody(&body),
		s.WithTrackTotalHits(true),
		s.WithSize(1), //TODO config?
	)

	if err != nil {
		err = fmt.Errorf("could not execute search request: %w", err)
		return
	}

	if res.IsError() {
		err = fmt.Errorf("request failed for search: %s", res.Status())
		return
	}

	var response SearchResponse
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
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

	// TODO Add exclude to query
	/*if c.exclude {
		if hit.(map[string]interface{})["_source"].(map[string]interface{})[c.excludeKey] == nil {
			for key, _ := range hit.(map[string]interface{})["_source"].(map[string]interface{}) {
				availableMessageKeys += key + "\n"
			}
			err = fmt.Errorf("exclude key: \""+
				c.excludeKey+"\" was not found. Available keys:\n%v", availableMessageKeys)

			return
		}
	}*/

	return
}
