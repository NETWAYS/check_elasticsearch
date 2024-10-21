package cmd

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"strings"
	"testing"
)

func TestSnapshot_ConnectionRefused(t *testing.T) {

	cmd := exec.Command("go", "run", "../main.go", "snapshot", "--port", "9999")
	out, _ := cmd.CombinedOutput()

	actual := string(out)
	expected := "[UNKNOWN] - could not fetch snapshots: Get \"http://localhost:9999/_snapshot/*/*?order=desc\": dial"

	if !strings.Contains(actual, expected) {
		t.Error("\nActual: ", actual, "\nExpected: ", expected)
	}
}

func TestSnapshot_WithWrongFlags(t *testing.T) {

	cmd := exec.Command("go", "run", "../main.go", "snapshot", "--all", "--number", "9999")
	out, _ := cmd.CombinedOutput()

	actual := string(out)
	expected := "[UNKNOWN] - if any flags in the group"

	if !strings.Contains(actual, expected) {
		t.Error("\nActual: ", actual, "\nExpected: ", expected)
	}
}

type SnapshotTest struct {
	name     string
	server   *httptest.Server
	args     []string
	expected string
}

func TestSnapshotCmd(t *testing.T) {
	tests := []SnapshotTest{
		{
			name: "snapshot-invalid-response",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`Hey dude where my snapshot`))
			})),
			args:     []string{"run", "../main.go", "snapshot"},
			expected: "[UNKNOWN] - could not decode snapshot response",
		},
		{
			name: "snapshot-none-available-ok",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"snapshots":[],"total":0,"remaining":0}`))
			})),
			args:     []string{"run", "../main.go", "snapshot", "--no-snapshots-state", "OK"},
			expected: "[OK] - No snapshots found",
		},
		{
			name: "snapshot-none-available-warning",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"snapshots":[],"total":0,"remaining":0}`))
			})),
			args:     []string{"run", "../main.go", "snapshot", "--no-snapshots-state", "WARNING"},
			expected: "[WARNING] - No snapshots found",
		},
		{
			name: "snapshot-none-available-critical",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"snapshots":[],"total":0,"remaining":0}`))
			})),
			args:     []string{"run", "../main.go", "snapshot", "--no-snapshots-state", "CRITICAL"},
			expected: "[CRITICAL] - No snapshots found",
		},
		{
			name: "snapshot-none-available-unknown",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"snapshots":[],"total":0,"remaining":0}`))
			})),
			args:     []string{"run", "../main.go", "snapshot", "--no-snapshots-state", "UNKNOWN"},
			expected: "[UNKNOWN] - No snapshots found",
		},
		{
			name: "snapshot-none-available-default",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"snapshots":[],"total":0,"remaining":0}`))
			})),
			args:     []string{"run", "../main.go", "snapshot"},
			expected: "[UNKNOWN] - No snapshots found",
		},
		{
			name: "snapshot-ok",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"snapshots":[{"snapshot":"snapshot_1","uuid":"dKb54xw67gvdRctLCxSket","repository":"my_repository","version_id":1.1,"version":1,"indices":[],"data_streams":[],"feature_states":[],"include_global_state":true,"state":"SUCCESS","start_time":"2020-07-06T21:55:18.129Z","start_time_in_millis":1593093628850,"end_time":"2020-07-06T21:55:18.129Z","end_time_in_millis":1593094752018,"duration_in_millis":0,"failures":[{"fail": "didnotfindapidocsforfailures"}],"shards":{"total":0,"failed":0,"successful":0}},{"snapshot":"snapshot_2","uuid":"vdRctLCxSketdKb54xw67g","repository":"my_repository","version_id":2,"version":2,"indices":[],"data_streams":[],"feature_states":[],"include_global_state":true,"state":"SUCCESS","start_time":"2020-07-06T21:55:18.130Z","start_time_in_millis":1593093628851,"end_time":"2020-07-06T21:55:18.130Z","end_time_in_millis":1593094752019,"duration_in_millis":1,"failures":[],"shards":{"total":0,"failed":0,"successful":0}}],"next":"c25hcHNob3RfMixteV9yZXBvc2l0b3J5LHNuYXBzaG90XzI=","total":3,"remaining":1}
`))
			})),
			args:     []string{"run", "../main.go", "snapshot"},
			expected: "[OK] - All evaluated snapshots are in state SUCCESS",
		},
		{
			name: "snapshot-inprogress",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"snapshots":[{"snapshot":"snapshot_1","uuid":"dKb54xw67gvdRctLCxSket","repository":"my_repository","version_id":1,"version":1,"indices":[],"data_streams":[],"feature_states":[],"include_global_state":true,"state":"IN_PROGRESS","start_time":"2020-07-06T21:55:18.129Z","start_time_in_millis":1593093628850,"end_time":"2020-07-06T21:55:18.129Z","end_time_in_millis":1593094752018,"duration_in_millis":0,"failures":[],"shards":{"total":0,"failed":0,"successful":0}},{"snapshot":"snapshot_2","uuid":"vdRctLCxSketdKb54xw67g","repository":"my_repository","version_id":2,"version":2,"indices":[],"data_streams":[],"feature_states":[],"include_global_state":true,"state":"SUCCESS","start_time":"2020-07-06T21:55:18.130Z","start_time_in_millis":1593093628851,"end_time":"2020-07-06T21:55:18.130Z","end_time_in_millis":1593094752019,"duration_in_millis":1,"failures":[],"shards":{"total":0,"failed":0,"successful":0}}],"next":"c25hcHNob3RfMixteV9yZXBvc2l0b3J5LHNuYXBzaG90XzI=","total":3,"remaining":1}
`))
			})),
			args:     []string{"run", "../main.go", "snapshot"},
			expected: "[UNKNOWN] - At least one evaluated snapshot is in state IN_PROGRESS",
		},
		{
			name: "snapshot-failed-with-all",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"snapshots":[{"snapshot":"snapshot_1","uuid":"dKb54xw67gvdRctLCxSket","repository":"my_repository","version_id":1,"version":1,"indices":[],"data_streams":[],"feature_states":[],"include_global_state":true,"state":"IN_PROGRESS","start_time":"2020-07-06T21:55:18.129Z","start_time_in_millis":1593093628850,"end_time":"2020-07-06T21:55:18.129Z","end_time_in_millis":1593094752018,"duration_in_millis":0,"failures":[],"shards":{"total":0,"failed":0,"successful":0}},{"snapshot":"snapshot_2","uuid":"vdRctLCxSketdKb54xw67g","repository":"my_repository","version_id":2,"version":2,"indices":[],"data_streams":[],"feature_states":[],"include_global_state":true,"state":"FAILED","start_time":"2020-07-06T21:55:18.130Z","start_time_in_millis":1593093628851,"end_time":"2020-07-06T21:55:18.130Z","end_time_in_millis":1593094752019,"duration_in_millis":1,"failures":[],"shards":{"total":0,"failed":0,"successful":0}}],"next":"c25hcHNob3RfMixteV9yZXBvc2l0b3J5LHNuYXBzaG90XzI=","total":3,"remaining":1}
`))
			})),
			args:     []string{"run", "../main.go", "snapshot", "--all"},
			expected: "[CRITICAL] - At least one evaluated snapshot is in state FAILED",
		},
		{
			name: "snapshot-partial-with-number",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"snapshots":[{"snapshot":"snapshot_1","uuid":"dKb54xw67gvdRctLCxSket","repository":"my_repository","version_id":1,"version":1,"indices":[],"data_streams":[],"feature_states":[],"include_global_state":true,"state":"SUCCESS","start_time":"2020-07-06T21:55:18.129Z","start_time_in_millis":1593093628850,"end_time":"2020-07-06T21:55:18.129Z","end_time_in_millis":1593094752018,"duration_in_millis":0,"failures":[],"shards":{"total":0,"failed":0,"successful":0}},{"snapshot":"snapshot_2","uuid":"vdRctLCxSketdKb54xw67g","repository":"my_repository","version_id":2,"version":2,"indices":[],"data_streams":[],"feature_states":[],"include_global_state":true,"state":"PARTIAL","start_time":"2020-07-06T21:55:18.130Z","start_time_in_millis":1593093628851,"end_time":"2020-07-06T21:55:18.130Z","end_time_in_millis":1593094752019,"duration_in_millis":1,"failures":[],"shards":{"total":0,"failed":0,"successful":0}}],"next":"c25hcHNob3RfMixteV9yZXBvc2l0b3J5LHNuYXBzaG90XzI=","total":3,"remaining":1}
`))
			})),
			args:     []string{"run", "../main.go", "snapshot", "--number", "4"},
			expected: "[WARNING] - At least one evaluated snapshot is in state PARTIAL",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.server.Close()

			// We need the random Port extracted
			u, _ := url.Parse(test.server.URL)
			cmd := exec.Command("go", append(test.args, "--port", u.Port())...)
			out, _ := cmd.CombinedOutput()

			actual := string(out)

			if !strings.Contains(actual, test.expected) {
				t.Error("\nActual: ", actual, "\nExpected: ", test.expected)
			}

		})
	}
}
