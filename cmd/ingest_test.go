package cmd

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"strings"
	"testing"
)

func TestIngest_ConnectionRefused(t *testing.T) {

	cmd := exec.Command("go", "run", "../main.go", "ingest", "--port", "9999")
	out, _ := cmd.CombinedOutput()

	actual := string(out)
	expected := "[UNKNOWN] - could not fetch cluster nodes statistics: Get \"http://localhost:9999/_nodes/stats\": dial"

	if !strings.Contains(actual, expected) {
		t.Error("\nActual: ", actual, "\nExpected: ", expected)
	}
}

type IngestTest struct {
	name     string
	server   *httptest.Server
	args     []string
	expected string
}

func TestIngestCmd(t *testing.T) {
	tests := []IngestTest{
		{
			name: "ingest-ok",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"_nodes":{"total":1,"successful":1,"failed":0},"cluster_name":"clustername","nodes":{"node1":{"ip":"127.0.0.1:9300","ingest":{"total":{"count":10,"time_in_millis":0,"current":3,"failed":5},"pipelines":{"mypipeline":{"count":10,"time_in_millis":0,"current":3,"failed":5,"processors":[{"set":{"type":"set","stats":{"count":0,"time_in_millis":0,"current":0,"failed":0}}}]}}}}}}`))
			})),
			args:     []string{"run", "../main.go", "ingest"},
			expected: "[OK] - Ingest operations alright \n \\_[OK] Number of failed ingest operations for mypipeline: 5; | pipelines.mypipeline.failed=5c pipelines.mypipeline.count=10c pipelines.mypipeline.current=3c\n",
		},
		{
			name: "ingest-ok-with-name",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"_nodes":{"total":1,"successful":1,"failed":0},"cluster_name":"clustername","nodes":{"mm3U-u0WTCeuZ325vZGY2w":{"ip":"127.0.0.1:9300","ingest":{"total":{"count":10,"time_in_millis":0,"current":3,"failed":5},"pipelines":{"foobar":{"count":10,"time_in_millis":0,"current":3,"failed":5,"processors":[{"set":{"type":"set","stats":{"count":0,"time_in_millis":0,"current":0,"failed":0}}}]},"mypipeline":{"count":10,"time_in_millis":0,"current":3,"failed":5,"processors":[{"set":{"type":"set","stats":{"count":0,"time_in_millis":0,"current":0,"failed":0}}}]}}}}}}`))
			})),
			args:     []string{"run", "../main.go", "ingest", "--pipeline", "foobar", "--pipeline", "notpresent"},
			expected: "[OK] - Ingest operations alright \n \\_[OK] Number of failed ingest operations for foobar: 5; | pipelines.foobar.failed=5c pipelines.foobar.count=10c pipelines.foobar.current=3c\n",
		},
		{
			name: "ingest-warn",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"_nodes":{"total":1,"successful":1,"failed":0},"cluster_name":"clustername","nodes":{"node1":{"ip":"127.0.0.1:9300","ingest":{"total":{"count":10,"time_in_millis":0,"current":3,"failed":5},"pipelines":{"mypipeline":{"count":10,"time_in_millis":0,"current":3,"failed":5,"processors":[{"set":{"type":"set","stats":{"count":0,"time_in_millis":0,"current":0,"failed":0}}}]}}}}}}`))
			})),
			args:     []string{"run", "../main.go", "ingest", "--failed-warning", "3"},
			expected: "[WARNING] - Ingest operations may not be alright \n \\_[WARNING] Number of failed ingest operations for mypipeline: 5; | pipelines.mypipeline.failed=5c pipelines.mypipeline.count=10c pipelines.mypipeline.current=3c\nexit status 1\n",
		},
		{
			name: "ingest-crit",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"_nodes":{"total":1,"successful":1,"failed":0},"cluster_name":"clustername","nodes":{"node1":{"ip":"127.0.0.1:9300","ingest":{"total":{"count":10,"time_in_millis":0,"current":3,"failed":5},"pipelines":{"mypipeline":{"count":10,"time_in_millis":0,"current":3,"failed":5,"processors":[{"set":{"type":"set","stats":{"count":0,"time_in_millis":0,"current":0,"failed":0}}}]}}}}}}`))
			})),
			args:     []string{"run", "../main.go", "ingest", "--failed-critical", "3"},
			expected: "[CRITICAL] - Ingest operations not alright \n \\_[CRITICAL] Number of failed ingest operations for mypipeline: 5; | pipelines.mypipeline.failed=5c pipelines.mypipeline.count=10c pipelines.mypipeline.current=3c\nexit status 2\n",
		},
		{
			name: "ingest-invalid",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{}`))
			})),
			args:     []string{"run", "../main.go", "ingest"},
			expected: "[UNKNOWN] - Ingest operations status unknown  | \nexit status 3\n",
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
