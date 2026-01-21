package cmd

import (
	"fmt"
	"strings"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/spf13/cobra"
)

type QueryConfig struct {
	Index      string
	Query      string
	MessageKey string
	MessageLen int
	Critical   string
	Warning    string
}

var (
	cliQueryConfig QueryConfig
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Checks the total hits/results of an Elasticsearch query (using a query_string query type)",
	Long: `Checks the total hits/results of an Elasticsearch query (using a query_string query type).
The plugin is currently capable to return the total hits of documents based on a provided query string.

For more information to the syntax, please visit:
https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html`,
	Example: "check_elasticsearch query -q \"event.dataset:sample_web_logs and @timestamp:[now-5m TO now]\" " +
		"-I \"kibana_sample_data_logs\" -k \"message\"",
	Run: func(_ *cobra.Command, _ []string) {
		var (
			rc     int
			output strings.Builder
		)

		client := cliConfig.NewClient()

		total, messages, err := client.SearchMessages(
			cliQueryConfig.Index,
			cliQueryConfig.Query,
			cliQueryConfig.MessageKey)

		if err != nil {
			check.ExitError(err)
		}

		output.WriteString(fmt.Sprintf("Search query hits: %d", total))

		if len(messages) > 0 {
			output.WriteString("\n")
			for _, msg := range messages {
				if len(msg) > cliQueryConfig.MessageLen {
					msg = msg[0:cliQueryConfig.MessageLen]
				}
				output.WriteString(msg + "\n")
			}
		}

		crit, err := check.ParseThreshold(cliQueryConfig.Critical)
		if err != nil {
			check.ExitError(err)
		}

		warn, err := check.ParseThreshold(cliQueryConfig.Warning)
		if err != nil {
			check.ExitError(err)
		}

		if crit.DoesViolate(float64(total)) {
			rc = check.Critical
		} else if warn.DoesViolate(float64(total)) {
			rc = check.Warning
		} else {
			rc = check.OK
		}

		p := perfdata.PerfdataList{
			{Label: "query_hits", Value: total, Warn: warn, Crit: crit, Uom: "c"},
		}

		check.ExitRaw(rc, output.String(), "|", p.String())
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)

	fs := queryCmd.Flags()
	fs.StringVarP(&cliQueryConfig.Query, "query", "q", "",
		"The Elasticsearch query to run (query_string type syntax)")
	fs.StringVarP(&cliQueryConfig.Index, "index", "I", "_all",
		"Name of the Index which will be used")
	fs.StringVarP(&cliQueryConfig.MessageKey, "msgkey", "k", "",
		"Name of a field to display in the output (e.g. a message body)")
	fs.IntVarP(&cliQueryConfig.MessageLen, "msglen", "m", 80,
		"Maximum number of characters to display from the requested field (default 80)")
	fs.StringVarP(&cliQueryConfig.Warning, "warning", "w", "20",
		"Warning threshold for total hits")
	fs.StringVarP(&cliQueryConfig.Critical, "critical", "c", "50",
		"Critical threshold for total hits")

	fs.SortFlags = false
}
