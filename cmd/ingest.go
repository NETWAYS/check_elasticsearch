package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/NETWAYS/go-check/result"
	"github.com/spf13/cobra"
)

// To store the CLI parameters.
type PipelineConfig struct {
	PipelineNames  []string
	FailedWarning  string
	FailedCritical string
}

const ingestOutput = "%s Number of failed ingest operations for %s: %g;"

var cliPipelineConfig PipelineConfig

var ingestCmd = &cobra.Command{
	Use:   "ingest",
	Short: "Checks the ingest statistics of Ingest Pipelines",
	Long:  `Checks the ingest statistics of Ingest Pipelines`,
	Run: func(_ *cobra.Command, _ []string) {
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

				pipelineMatched, regexErr := matches(pipelineName, cliPipelineConfig.PipelineNames)

				if regexErr != nil {
					check.ExitRaw(check.Unknown, "Invalid regular expression provided:", regexErr.Error())
				}

				if !pipelineMatched && len(cliPipelineConfig.PipelineNames) >= 1 {
					// If the pipeline doesn't matches a regex from the list we can skip it.
					continue
				}

				summary.WriteString("\n \\_")
				if failedCrit.DoesViolate(pp.Failed) {
					states = append(states, check.Critical)
					summary.WriteString(fmt.Sprintf(ingestOutput, "[CRITICAL]", pipelineName, pp.Failed))
				} else if failedWarn.DoesViolate(pp.Failed) {
					states = append(states, check.Warning)
					summary.WriteString(fmt.Sprintf(ingestOutput, "[WARNING]", pipelineName, pp.Failed))
				} else {
					states = append(states, check.OK)
					summary.WriteString(fmt.Sprintf(ingestOutput, "[OK]", pipelineName, pp.Failed))
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

	fs.StringArrayVar(&cliPipelineConfig.PipelineNames, "pipeline", []string{},
		"Name of the pipeline to check. Can be used multiple times and supports regex.")
	fs.StringVar(&cliPipelineConfig.FailedWarning, "failed-warning", "10",
		"Warning threshold for failed ingest operations. Use min:max for a range.")
	fs.StringVar(&cliPipelineConfig.FailedCritical, "failed-critical", "20",
		"Critical threshold for failed ingest operations. Use min:max for a range.")

	fs.SortFlags = false
}

// Matches a list of regular expressions against a string.
func matches(input string, regexToMatch []string) (bool, error) {
	for _, regex := range regexToMatch {
		re, err := regexp.Compile(regex)

		if err != nil {
			return false, err
		}

		if re.MatchString(input) {
			return true, nil
		}
	}

	return false, nil
}
