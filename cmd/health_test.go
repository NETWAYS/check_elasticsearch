package cmd

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"strings"
	"testing"
)

func TestHealth_ConnectionRefused(t *testing.T) {

	cmd := exec.Command("go", "run", "../main.go", "health", "--port", "9999")
	out, _ := cmd.CombinedOutput()

	actual := string(out)
	expected := "[UNKNOWN] - could not fetch cluster health: Get \"http://localhost:9999/_cluster/health\": dial"

	if !strings.Contains(actual, expected) {
		t.Error("\nActual: ", actual, "\nExpected: ", expected)
	}
}

type HealthTest struct {
	name     string
	server   *httptest.Server
	args     []string
	expected string
}

func TestHealthCmd(t *testing.T) {
	tests := []HealthTest{
		{
			name: "health-basic-auth-ok",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				user, pass, ok := r.BasicAuth()
				if ok {
					// Just for testing, this is now how to handle BasicAuth properly
					if user == "username" && pass == "password" {
						w.Header().Set("X-Elastic-Product", "Elasticsearch")
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`{"cluster_name":"test","status":"green","timed_out":false,"number_of_nodes":1,"number_of_data_nodes":1,"active_primary_shards":3,"active_shards":3,"relocating_shards":0,"initializing_shards":0,"unassigned_shards":0,"delayed_unassigned_shards":0,"number_of_pending_tasks":0,"number_of_in_flight_fetch":0,"task_max_waiting_in_queue_millis":0,"active_shards_percent_as_number":100.0}`))
						return
					}
				}
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`The Authorization header wasn't set`))
			})),
			args:     []string{"run", "../main.go", "health", "--username", "username", "--password", "password"},
			expected: "[OK] - Cluster test is green | status=0 nodes=1 data_nodes=1 active_primary_shards=3 active_shards=3\n",
		},
		{
			name: "health-invalid",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{}`))
			})),
			args:     []string{"run", "../main.go", "health"},
			expected: "[UNKNOWN] - Cluster status unknown | status=3 nodes=0 data_nodes=0 active_primary_shards=0 active_shards=0\nexit status 3\n",
		},
		{
			name: "health-404",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{}`))
			})),
			args:     []string{"run", "../main.go", "health"},
			expected: "[UNKNOWN] - request failed for cluster health: 401 Unauthorized (*errors.errorString)\nexit status 3\n",
		},
		{
			name: "health-unknown",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"cluster_name":"test","status":"foobar","timed_out":false,"number_of_nodes":1,"number_of_data_nodes":1,"active_primary_shards":3}`))
			})),
			args:     []string{"run", "../main.go", "health"},
			expected: "[UNKNOWN] - Cluster test is foobar | status=3 nodes=1 data_nodes=1 active_primary_shards=3 active_shards=0\nexit status 3\n",
		},
		{
			name: "health-ok",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"cluster_name":"test","status":"green","timed_out":false,"number_of_nodes":1,"number_of_data_nodes":1,"active_primary_shards":3,"active_shards":3,"relocating_shards":0,"initializing_shards":0,"unassigned_shards":0,"delayed_unassigned_shards":0,"number_of_pending_tasks":0,"number_of_in_flight_fetch":0,"task_max_waiting_in_queue_millis":0,"active_shards_percent_as_number":100.0}`))
			})),
			args:     []string{"run", "../main.go", "health"},
			expected: "[OK] - Cluster test is green | status=0 nodes=1 data_nodes=1 active_primary_shards=3 active_shards=3\n",
		},
		{
			name: "health-yellow",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"cluster_name":"test","status":"yellow","timed_out":false,"number_of_nodes":1,"number_of_data_nodes":1,"active_primary_shards":3,"active_shards":3,"relocating_shards":0,"initializing_shards":0,"unassigned_shards":0,"delayed_unassigned_shards":0,"number_of_pending_tasks":0,"number_of_in_flight_fetch":0,"task_max_waiting_in_queue_millis":0,"active_shards_percent_as_number":100.0}`))
			})),
			args:     []string{"run", "../main.go", "health"},
			expected: "[WARNING] - Cluster test is yellow | status=1 nodes=1 data_nodes=1 active_primary_shards=3 active_shards=3\nexit status 1\n",
		},
		{
			name: "health-red",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"cluster_name":"test","status":"red","timed_out":false,"number_of_nodes":1,"number_of_data_nodes":1,"active_primary_shards":3,"active_shards":3,"relocating_shards":0,"initializing_shards":0,"unassigned_shards":0,"delayed_unassigned_shards":0,"number_of_pending_tasks":0,"number_of_in_flight_fetch":0,"task_max_waiting_in_queue_millis":0,"active_shards_percent_as_number":100.0}`))
			})),
			args:     []string{"run", "../main.go", "health"},
			expected: "[CRITICAL] - Cluster test is red | status=2 nodes=1 data_nodes=1 active_primary_shards=3 active_shards=3\nexit status 2\n",
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

			if actual != test.expected {
				t.Error("\nActual: ", actual, "\nExpected: ", test.expected)
			}

		})
	}
}
