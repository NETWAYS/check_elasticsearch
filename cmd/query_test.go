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

	expected := "[UNKNOWN] - could not execute search request: Get \"http://localhost:9999/_all/_search?size=1&track_total_hits=true\": dial"

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
			expected: "[OK] - Search query hits: 0 | query_hits=0c;20;50\n",
		},
		{
			name: "query-no-such-index",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error":{"root_cause":[{"type":"index_not_found_exception","reason":"no such index [example]","resource.type":"index_or_alias","resource.id":"kibaa_sample_data_logs","index_uuid":"_na_","index":"example"}],"type":"index_not_found_exception","reason":"no such index [example]","resource.type":"index_or_alias","resource.id":"example","index_uuid":"_na_","index":"example"},"status":404}`))
			})),
			args:     []string{"run", "../main.go", "query", "-q", "foo", "-I", "foo"},
			expected: "[UNKNOWN] - failed to run query: no such index [example] (*errors.errorString)\nexit status 3\n",
		},
		{
			name: "query-example",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"took":4,"timed_out":false,"_shards":{"total":1,"successful":1,"skipped":0,"failed":0},"hits":{"total":{"value":14074,"relation":"eq"},"max_score":8.386352,"hits":[{"_index":"kibana_sample_data_logs","_type":"_doc","_id":"dGX9CYgBFkvhWgFatiP9","_score":8.386352,"_source":{"agent":"Mozilla/5.0(X11;Linuxi686)AppleWebKit/534.24(KHTML,likeGecko)Chrome/11.0.696.50Safari/534.24","bytes":1831,"clientip":"30.156.16.164","extension":"","geo":{"srcdest":"US:IN","src":"US","dest":"IN","coordinates":{"lat":55.53741389,"lon":-132.3975144}},"host":"elastic-elastic-elastic.org","index":"kibana_sample_data_logs","ip":"30.156.16.163","machine":{"ram":9663676416,"os":"winxp"},"memory":73240,"message":"30.156.16.163--[2018-09-01T12:44:53.756Z]\"GET/wp-content/HTTP/1.1\"4041831\"-\"\"Mozilla/5.0(X11;Linuxi686)AppleWebKit/534.24(KHTML,likeGecko)Chrome/11.0.696.50Safari/534.24\"","phpmemory":73240,"referer":"http://www.elastic-elastic-elastic.com/success/timothy-l-kopra","request":"/wp-content/","response":404,"tags":["success","info"],"timestamp":"2023-06-10T12:44:53.756Z","url":"https://elastic-elastic-elastic.org/wp-content//","utc_time":"2023-06-10T12:44:53.756Z","event":{"dataset":"sample_web_logs"}}}]}}`))
			})),
			args: []string{"run", "../main.go", "query", "-q", "vent.dataset:sample_web_logs and @timestamp:[now-5m TO now]", "-I", "kibana_sample_data_logs", "-k", "message"},
			expected: `[CRITICAL] - Search query hits: 14074
30.156.16.163--[2018-09-01T12:44:53.756Z]"GET/wp-content/HTTP/1.1"4041831"-""Moz
 | query_hits=14074c;20;50
exit status 2
`,
		},
		{
			name: "query-ok",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"took":3,"timed_out":false,"_shards":{"total":1,"successful":1,"skipped":0,"failed":0},"hits":{"total":{"value":2,"relation":"eq"},"max_score":1.0,"hits":[{"_index":"my_index","_type":"_doc","_id":"yUi6voQB87C1kW3InC4l","_score":1.0,"_source":{"title":"One","tags":["ruby"]}},{"_index":"my_index","_type":"_doc","_id":"y0i9voQB87C1kW3I9y74","_score":1.0,"_source":{"title":"One","tags":["ruby"]}}]}}`))
			})),
			args: []string{"run", "../main.go", "query", "-I", "my_index", "-q", "*", "--msgkey", "title", "-w", "3"},
			expected: `[OK] - Search query hits: 2
One
One
 | query_hits=2c;3;50
`,
		},
		{
			name: "query-ok-no-msgkey",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Elastic-Product", "Elasticsearch")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"took":3,"timed_out":false,"_shards":{"total":1,"successful":1,"skipped":0,"failed":0},"hits":{"total":{"value":2,"relation":"eq"},"max_score":1.0,"hits":[{"_index":"my_index","_type":"_doc","_id":"yUi6voQB87C1kW3InC4l","_score":1.0,"_source":{"title":"One","tags":["ruby"]}},{"_index":"my_index","_type":"_doc","_id":"y0i9voQB87C1kW3I9y74","_score":1.0,"_source":{"title":"One","tags":["ruby"]}}]}}`))
			})),
			args: []string{"run", "../main.go", "query", "-I", "my_index", "-q", "*", "-w", "3"},
			expected: `[OK] - Search query hits: 2 | query_hits=2c;3;50
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
			expected: `[CRITICAL] - Search query hits: 2
One
One
 | query_hits=2c;20;1
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
			expected: `[WARNING] - Search query hits: 2
One
One
 | query_hits=2c;1;50
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
			expected: "[UNKNOWN] - error parsing the response body: json: cannot unmarshal string into Go struct field SearchTotal.hits.total.value of type uint (*fmt.wrapError)\nexit status 3\n",
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
