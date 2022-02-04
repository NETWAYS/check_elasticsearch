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

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		} else if total == 0 {
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
	fs.StringVarP(&cliQueryConfig.Query, "query", "q", "", "")
	fs.StringVarP(&cliQueryConfig.Index, "index", "I", "", "")
	fs.StringVarP(&cliQueryConfig.MessageKey, "msgkey", "k", "", "")
	fs.IntVarP(&cliQueryConfig.MessageLen, "msglen", "m", 80, "")
	fs.UintVarP(&cliQueryConfig.Warning, "warning", "w", 20, "")
	fs.UintVarP(&cliQueryConfig.Critical, "critical", "w", 50, "")
}
