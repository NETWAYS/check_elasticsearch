package cmd

import (
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Checks the health status of an Elasticsearch cluster",
	Long: `Checks the health status of an Elasticsearch cluster

The cluster health status is:
	green = OK
	yellow = WARNING
	red = CRITICAL`,
	Example: "  check_elasticsearch health --hostname \"127.0.0.1\" --port 9200 --username \"exampleUser\"  " +
		"--password \"examplePass\" --tls --insecure",
	Run: func(cmd *cobra.Command, args []string) {
		client := cliConfig.Client()
		err := client.Connect()
		if err != nil {
			check.ExitError(err)
		}

		health, err := client.Health()
		if err != nil {
			check.ExitError(err)
		}

		var rc int
		switch health.Status {
		case "green":
			rc = check.OK
		case "yellow":
			rc = check.Warning
		case "red":
			rc = check.Critical
		default:
			rc = check.Unknown
		}

		var output = "Cluster status unknown"
		if health.Status != "" {
			output = "Cluster " + health.ClusterName + " is " + health.Status
		}

		// green = 0
		// yellow = 1
		// red = 2
		// unknown = 3
		p := perfdata.PerfdataList{
			{Label: "status", Value: rc},
			{Label: "nodes", Value: health.NumberOfNodes},
			{Label: "data_nodes", Value: health.NumberOfDataNodes},
			{Label: "active_primary_shards", Value: health.ActivePrimaryShards},
			{Label: "active_shards", Value: health.ActiveShards},
		}

		check.ExitRaw(rc, output, "|", p.String())
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)
	healthCmd.DisableFlagsInUseLine = true
}
