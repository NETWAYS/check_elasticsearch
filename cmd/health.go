package cmd

import (
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/spf13/cobra"
)

// healthCmd represents the health command
var healthCmd = &cobra.Command{
	Use:   "health",
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

		health, err := client.Health()
		if err != nil {
			check.ExitError(err)
		}

		rc := 3
		output := "Cluster " + health.ClusterName + " is " + health.Status

		switch health.Status {
		case "green":
			rc = 0
		case "yellow":
			rc = 1
		default:
			rc = 2
		}

		// green = 0
		// yellow = 1
		// red = 2
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
