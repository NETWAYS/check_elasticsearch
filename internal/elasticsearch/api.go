package elasticsearch

import (
	"slices"
	"strings"
)

type HealthResponse struct {
	ClusterName                 string  `json:"cluster_name"`
	Status                      string  `json:"status"`
	TimedOut                    bool    `json:"timed_out"`
	NumberOfNodes               int     `json:"number_of_nodes"`
	NumberOfDataNodes           int     `json:"number_of_data_nodes"`
	ActivePrimaryShards         int     `json:"active_primary_shards"`
	ActiveShards                int     `json:"active_shards"`
	RelocatingShards            int     `json:"relocating_shards"`
	InitializingShards          int     `json:"initializing_shards"`
	UnassignedShards            int     `json:"unassigned_shards"`
	DelayedUnassignedShards     int     `json:"delayed_unassigned_shards"`
	NumberOfPendingTasks        int     `json:"number_of_pending_tasks"`
	NumberOfInFlightFetch       int     `json:"number_of_in_flight_fetch"`
	TaskMaxWaitingInQueueMillis int     `json:"task_max_waiting_in_queue_millis"`
	ActiveShardsPercentAsNumber float64 `json:"active_shards_percent_as_number"`
}

// https://www.elastic.co/guide/en/elasticsearch/reference/current/search-search.html#search-api-response-body
type SearchResponse struct {
	Hits  SearchHits `json:"hits"`
	Error struct {
		RootCause []ErrorRootCause `json:"root_cause"`
	}
}

type ErrorRootCause struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

// GetErrors returns the error reasons when they are present in the response
func (r *SearchResponse) GetErrors() string {
	if len(r.Error.RootCause) == 0 {
		return ""
	}

	messages := make([]string, 0, len(r.Error.RootCause))

	for _, rc := range r.Error.RootCause {
		if rc.Reason != "" {
			// Deduplication
			if slices.Contains(messages, rc.Reason) {
				continue
			}

			messages = append(messages, rc.Reason)
		}
	}

	return strings.Join(messages, ", ")
}

type SearchHits struct {
	Total SearchTotal `json:"total"`
	Hits  []SearchHit `json:"hits"`
}

type SearchTotal struct {
	Value uint `json:"value"`
}

type SearchHit struct {
	Index  string         `json:"_index"`
	Type   string         `json:"_type"`
	Source map[string]any `json:"_source"`
	ID     string         `json:"_id"`
}

type SearchRequest struct {
	Query Query `json:"query"`
}

// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html
type Query struct {
	QueryString *QueryString `json:"query_string,omitempty"`
}

// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html
type QueryString struct {
	Query string `json:"query"`
}

type NodeInfo struct {
	IP     string     `json:"ip"`
	Ingest IngestInfo `json:"ingest"`
}

type IngestInfo struct {
	Total     IngestStats             `json:"total"`
	Pipelines map[string]PipelineInfo `json:"pipelines"`
}

type IngestStats struct {
	Count   float64 `json:"count"`
	Current float64 `json:"current"`
	Failed  float64 `json:"failed"`
}

type PipelineInfo struct {
	Count   float64 `json:"count"`
	Current float64 `json:"current"`
	Failed  float64 `json:"failed"`
}

type ClusterStats struct {
	Nodes       map[string]NodeInfo `json:"nodes"`
	ClusterName string              `json:"cluster_name"`
}

type Snapshot struct {
	Snapshot           string   `json:"snapshot"`
	UUID               string   `json:"uuid"`
	Repository         string   `json:"repository"`
	Indices            []string `json:"indices"`
	DataStreams        []string `json:"data_streams"`
	IncludeGlobalState bool     `json:"include_global_state"`
	State              string   `json:"state"`
	StartTimeInMillis  int      `json:"start_time_in_millis"`
	EndTimeInMillis    int      `json:"end_time_in_millis"`
	DurationInMillis   int      `json:"duration_in_millis"`
	Shards             struct {
		Total      int `json:"total"`
		Failed     int `json:"failed"`
		Successful int `json:"successful"`
	} `json:"shards"`
}

type SnapshotResponse struct {
	Snapshots []Snapshot `json:"snapshots"`
	Total     int        `json:"total"`
	Remaining int        `json:"remaining"`
}
