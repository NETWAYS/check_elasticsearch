package cmd

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"strings"
	"testing"
)

func TestQuery_ConnectionRefused(t *testing.T) {

	cmd := exec.Command("go", "run", "../main.go", "query", "--port", "9999")
	out, _ := cmd.CombinedOutput()

	actual := string(out)
	expected := "UNKNOWN - could not fetch cluster info: dial"

	if !strings.Contains(actual, expected) {
		t.Error("\nActual: ", actual, "\nExpected: ", expected)
	}
}

type QueryTest struct {
	name     string
	server   *httptest.Server
	args     []string
	expected string
}

func TestQueryCmd(t *testing.T) {
	tests := []QueryTest{
		{
			name: "query-empty-return",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{}`))
			})),
			args:     []string{"run", "../main.go", "query"},
			expected: "OK - Total hits: 0 | total=0;20;50\n",
		},
		{
			name: "query-default",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"name":"test","cluster_name":"cluster","cluster_uuid":"6nZDLRvSQ1iDxZUTf0Hrmg","version":{"number":"7.17.7","build_flavor":"default","build_type":"docker","build_hash":"78dcaaa8cee33438b91eca7f5c7f56a70fec9e80","build_date":"2022-10-17T15:29:54.167373105Z","build_snapshot":false,"lucene_version":"8.11.1","minimum_wire_compatibility_version":"6.8.0","minimum_index_compatibility_version":"6.0.0-beta1"},"tagline":"YouKnow,forSearch"}`))
			})),
			args:     []string{"run", "../main.go", "query"},
			expected: "OK - Total hits: 0 | total=0;20;50\n",
		},
		{
			name: "query-ok",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"took":3,"timed_out":false,"_shards":{"total":1,"successful":1,"skipped":0,"failed":0},"hits":{"total":{"value":2,"relation":"eq"},"max_score":1.0,"hits":[{"_index":"my_index","_type":"_doc","_id":"yUi6voQB87C1kW3InC4l","_score":1.0,"_source":{"title":"One","tags":["ruby"]}},{"_index":"my_index","_type":"_doc","_id":"y0i9voQB87C1kW3I9y74","_score":1.0,"_source":{"title":"One","tags":["ruby"]}}]}}`))
			})),
			args: []string{"run", "../main.go", "query", "-I", "my_index", "-q", "*", "--msgkey", "title", "-w", "3"},
			expected: `OK - Total hits: 2
One
One
 | total=2;3;50
`,
		},
		{
			name: "query-critical",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"took":3,"timed_out":false,"_shards":{"total":1,"successful":1,"skipped":0,"failed":0},"hits":{"total":{"value":2,"relation":"eq"},"max_score":1.0,"hits":[{"_index":"my_index","_type":"_doc","_id":"yUi6voQB87C1kW3InC4l","_score":1.0,"_source":{"title":"One","tags":["ruby"]}},{"_index":"my_index","_type":"_doc","_id":"y0i9voQB87C1kW3I9y74","_score":1.0,"_source":{"title":"One","tags":["ruby"]}}]}}`))
			})),
			args: []string{"run", "../main.go", "query", "-I", "my_index", "-q", "*", "--msgkey", "title", "-c", "1"},
			expected: `CRITICAL - Total hits: 2
One
One
 | total=2;20;1
exit status 2
`,
		},
		{
			name: "query-warning",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"took":3,"timed_out":false,"_shards":{"total":1,"successful":1,"skipped":0,"failed":0},"hits":{"total":{"value":2,"relation":"eq"},"max_score":1.0,"hits":[{"_index":"my_index","_type":"_doc","_id":"yUi6voQB87C1kW3InC4l","_score":1.0,"_source":{"title":"One","tags":["ruby"]}},{"_index":"my_index","_type":"_doc","_id":"y0i9voQB87C1kW3I9y74","_score":1.0,"_source":{"title":"One","tags":["ruby"]}}]}}`))
			})),
			args: []string{"run", "../main.go", "query", "-I", "my_index", "-q", "*", "--msgkey", "title", "-w", "1"},
			expected: `WARNING - Total hits: 2
One
One
 | total=2;1;50
exit status 1
`,
		},
		{
			name: "query-invalid-json",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"took":3,"timed_out":false,"_shards":{"total":1,"successful":1,"skipped":0,"failed":0},"hits":{"total":{"value":"foo","relation":"eq"},"max_score":"bar","hits":[{"_index":"my_index","_type":"_doc","_id":"yUi6voQB87C1kW3InC4l","_score":"bla","_source":{"title":"One","tags":["ruby"]}}]}}`))
			})),
			args:     []string{"run", "../main.go", "query", "-I", "my_index", "-q", "*", "--msgkey", "title", "-w", "1"},
			expected: "UNKNOWN - error parsing the response body: json: cannot unmarshal string into Go struct field SearchTotal.hits.total.value of type uint (*fmt.wrapError)\nexit status 3\n",
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
