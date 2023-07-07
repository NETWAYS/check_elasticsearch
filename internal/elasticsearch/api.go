package elasticsearch

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
	Hits SearchHits `json:"hits"`
}

type SearchHits struct {
	Total SearchTotal `json:"total"`
	Hits  []SearchHit `json:"hits"`
}

type SearchTotal struct {
	Value uint `json:"value"`
}

type SearchHit struct {
	Index  string                 `json:"_index"`
	Type   string                 `json:"_type"`
	Source map[string]interface{} `json:"_source"`
	Id     string                 `json:"_id"`
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
