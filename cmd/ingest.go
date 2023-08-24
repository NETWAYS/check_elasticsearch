package cmd

import (
	"fmt"
	"strings"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/NETWAYS/go-check/result"
	"github.com/spf13/cobra"
)

// To store the CLI parameters.
type PipelineConfig struct {
	PipelineName   string
	FailedWarning  string
	FailedCritical string
}

var cliPipelineConfig PipelineConfig

var ingestCmd = &cobra.Command{
	Use:   "ingest",
	Short: "Checks the ingest statistics of Ingest Pipelines",
	Long:  `Checks the ingest statistics of Ingest Pipelines`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			rc       int
			output   string
			perfList perfdata.PerfdataList
		)

		failedCrit, err := check.ParseThreshold(cliPipelineConfig.FailedCritical)
		if err != nil {
			check.ExitError(err)
		}

		failedWarn, err := check.ParseThreshold(cliPipelineConfig.FailedWarning)
		if err != nil {
			check.ExitError(err)
		}

		client := cliConfig.NewClient()

		stats, err := client.NodeStats()
		if err != nil {
			check.ExitError(err)
		}

		// Calculate states capacity
		amountOfNodes := 0
		for _, node := range stats.Nodes {
			amountOfNodes += len(node.Ingest.Pipelines)
		}

		states := make([]int, 0, amountOfNodes)

		// Check status for each pipeline
		var summary strings.Builder

		for _, node := range stats.Nodes {
			for pipelineName, pp := range node.Ingest.Pipelines {

				// Skip if pipeline name doesn't match
				if cliPipelineConfig.PipelineName != "" {
					if cliPipelineConfig.PipelineName != pipelineName {
						continue
					}
				}

				summary.WriteString("\n \\_")
				if failedCrit.DoesViolate(pp.Failed) {
					states = append(states, check.Critical)
					summary.WriteString(fmt.Sprintf("[CRITICAL] Failed ingest operations for %s: %g;", pipelineName, pp.Failed))
				} else if failedWarn.DoesViolate(pp.Failed) {
					states = append(states, check.Warning)
					summary.WriteString(fmt.Sprintf("[WARNING] Failed ingest operations for %s: %g;", pipelineName, pp.Failed))
				} else {
					states = append(states, check.OK)
					summary.WriteString(fmt.Sprintf("[OK] Failed ingest operations for %s: %g;", pipelineName, pp.Failed))
				}

				perfList.Add(&perfdata.Perfdata{
					Label: fmt.Sprintf("pipelines.%s.failed", pipelineName),
					Uom:   "c",
					Value: pp.Failed})
				perfList.Add(&perfdata.Perfdata{
					Label: fmt.Sprintf("pipelines.%s.count", pipelineName),
					Uom:   "c",
					Value: pp.Count})
				perfList.Add(&perfdata.Perfdata{
					Label: fmt.Sprintf("pipelines.%s.current", pipelineName),
					Uom:   "c",
					Value: pp.Current})
			}
		}

		// Validate the various subchecks and use the worst state as return code
		switch result.WorstState(states...) {
		case 0:
			rc = check.OK
			output = "Ingest operations alright"
		case 1:
			rc = check.Warning
			output = "Ingest operations may not be alright"
		case 2:
			rc = check.Critical
			output = "Ingest operations not alright"
		default:
			rc = check.Unknown
			output = "Ingest operations status unknown"
		}

		check.ExitRaw(rc, output, summary.String(), "|", perfList.String())
	},
}

func init() {
	rootCmd.AddCommand(ingestCmd)

	fs := ingestCmd.Flags()

	fs.StringVar(&cliPipelineConfig.PipelineName, "pipeline", "",
		"Pipeline Name")

	fs.StringVar(&cliPipelineConfig.FailedWarning, "failed-warning", "10",
		"Warning threshold for failed ingest operations. Use min:max for a range.")
	fs.StringVar(&cliPipelineConfig.FailedCritical, "failed-critical", "20",
		"Critical threshold for failed ingest operations. Use min:max for a range.")

	fs.SortFlags = false
}
