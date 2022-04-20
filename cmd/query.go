package cmd

import (
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/spf13/cobra"
)

type QueryConfig struct {
	Index      string
	Query      string
	MessageKey string
	MessageLen int
	Critical   uint
	Warning    uint
}

var cliQueryConfig QueryConfig

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Checks the total hits/results of an Elasticsearch query",
	Long: `Checks the total hits/results of an Elasticsearch query.
The plugin is currently capable to return the total hits of documents based on a provided query string.

For more information to the syntax, please visit:
https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html`,
	Example: "check_elasticsearch query -q \"event.dataset:sample_web_logs and @timestamp:[now-5m TO now]\" " +
		"-I \"kibana_sample_data_logs\" -k \"message\"",
	Run: func(cmd *cobra.Command, args []string) {
		client := cliConfig.Client()
		err := client.Connect()
		if err != nil {
			check.ExitError(err)
		}

		total, messages, err := client.SearchMessages(
			cliQueryConfig.Index,
			cliQueryConfig.Query,
			cliQueryConfig.MessageKey)

		if err != nil {
			check.ExitError(err)
		}

		rc := 3
		output := fmt.Sprintf("Total hits: %d", total)

		if len(messages) > 0 {
			output += "\n"
			for _, msg := range messages {
				if len(msg) > cliQueryConfig.MessageLen {
					msg = msg[0:cliQueryConfig.MessageLen]
				}
				output += msg + "\n"
			}
		}

		if total >= cliQueryConfig.Critical {
			rc = check.Critical
		} else if total >= cliQueryConfig.Warning {
			rc = check.Warning
		} else if total >= 0 {
			rc = check.OK
		}

		p := perfdata.PerfdataList{
			{Label: "total", Value: total,
				Warn: &check.Threshold{Upper: float64(cliQueryConfig.Warning)},
				Crit: &check.Threshold{Upper: float64(cliQueryConfig.Critical)},
			},
		}

		check.ExitRaw(rc, output, "|", p.String())
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)

	fs := queryCmd.Flags()
	fs.StringVarP(&cliQueryConfig.Query, "query", "q", "",
		"Elasticsearch query")
	fs.StringVarP(&cliQueryConfig.Index, "index", "I", "_all",
		"The index which will be used ")
	fs.StringVarP(&cliQueryConfig.MessageKey, "msgkey", "k", "",
		"Message of messagekey to display")
	fs.IntVarP(&cliQueryConfig.MessageLen, "msglen", "m", 80,
		"Number of characters to display in latest message")
	fs.UintVarP(&cliQueryConfig.Warning, "warning", "w", 20,
		"Warning threshold for total hits")
	fs.UintVarP(&cliQueryConfig.Critical, "critical", "c", 50,
		"Critical threshold for total hits")

	fs.SortFlags = false
}
