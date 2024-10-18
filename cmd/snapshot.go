package cmd

import (
	"fmt"
	"strings"

	"errors"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
	"github.com/spf13/cobra"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Checks the status of Elasticsearch snapshots",
	Long: `Checks the status of Elasticsearch snapshots.
The plugin maps snapshot status to the following exit codes:

SUCCESS, Exit code 0
PARTIAL, Exit code 1
FAILED, Exit code 2
IN_PROGRESS, Exit code 3

If there are multiple snapshots the plugin uses the worst status.
`,
	Example: `
$ check_elasticsearch snapshot
[OK] - All evaluated snapshots are in state SUCCESS

$ check_elasticsearch snapshot --all
[CRITICAL] - At least one evaluated snapshot is in state FAILED

$ check_elasticsearch snapshot --number 5
[WARNING] - At least one evaluated snapshot is in state PARTIAL
`,
	Run: func(cmd *cobra.Command, _ []string) {
		repository, _ := cmd.Flags().GetString("repository")
		snapshot, _ := cmd.Flags().GetString("snapshot")
		numberOfSnapshots, _ := cmd.Flags().GetInt("number")
		evalAllSnapshots, _ := cmd.Flags().GetBool("all")
		noSnapshotsState, _ := cmd.Flags().GetString("no-snapshots-state")

		// Convert --no-snapshots-state to integer and validate input
		noSnapshotsStateInt, err := convertStateToInt(noSnapshotsState)
		if err != nil {
			check.ExitError(fmt.Errorf("invalid value for --no-snapshots-state: %s", noSnapshotsState))
		}

		var (
			rc     int
			output string
		)

		client := cliConfig.NewClient()

		snapResponse, err := client.Snapshot(repository, snapshot)

		if err != nil {
			check.ExitError(err)
		}

		// If all snapshots are to be evaluated
		if evalAllSnapshots {
			numberOfSnapshots = len(snapResponse.Snapshots)
		}

		// If more snapshots are requested than available
		if numberOfSnapshots > len(snapResponse.Snapshots) {
			numberOfSnapshots = len(snapResponse.Snapshots)
		}

		// Evaluate snapshots given their states
		sStates := make([]int, 0, len(snapResponse.Snapshots))

		// Check status for each snapshot
		var summary strings.Builder

		for _, snap := range snapResponse.Snapshots[:numberOfSnapshots] {

			summary.WriteString("\n \\_")

			switch snap.State {
			default:
				sStates = append(sStates, check.Unknown)
				summary.WriteString(fmt.Sprintf("[UNKNOWN] Snapshot: %s, State %s, Repository: %s", snap.Snapshot, snap.State, snap.Repository))
			case "SUCCESS":
				sStates = append(sStates, check.OK)
				summary.WriteString(fmt.Sprintf("[OK] Snapshot: %s, State %s, Repository: %s", snap.Snapshot, snap.State, snap.Repository))
			case "PARTIAL":
				sStates = append(sStates, check.Warning)
				summary.WriteString(fmt.Sprintf("[WARNING] Snapshot: %s, State %s, Repository: %s", snap.Snapshot, snap.State, snap.Repository))
			case "FAILED":
				sStates = append(sStates, check.Critical)
				summary.WriteString(fmt.Sprintf("[CRITICAL] Snapshot: %s, State %s, Repository: %s", snap.Snapshot, snap.State, snap.Repository))
			case "IN PROGRESS":
				sStates = append(sStates, check.Unknown)
				summary.WriteString(fmt.Sprintf("[UNKNOWN] Snapshot: %s, State %s, Repository: %s", snap.Snapshot, snap.State, snap.Repository))
			}
		}

		if len(snapResponse.Snapshots) == 0 {
			switch noSnapshotsStateInt {
			case 0:
				sStates = append(sStates, check.OK)
			case 1:
				sStates = append(sStates, check.Warning)
			case 2:
				sStates = append(sStates, check.Critical)
			case 3:
				sStates = append(sStates, check.Unknown)
			}
		}

		rc = result.WorstState(sStates...)

		if len(snapResponse.Snapshots) == 0 {
			output = "No snapshots found."
		} else {
			switch rc {
			case check.OK:
				output = "All evaluated snapshots are in state SUCCESS."
			case check.Warning:
				output = "At least one evaluated snapshot is in state PARTIAL."
			case check.Critical:
				output = "At least one evaluated snapshot is in state FAILED."
			case check.Unknown:
				output = "At least one evaluated snapshot is in state IN_PROGRESS."
			default:
				output = "Could not evaluate status of snapshots"
			}
		}

		check.ExitRaw(rc, output, "repository:", repository, "snapshot:", snapshot, summary.String())
	},
}

// Function to convert state to integer
func convertStateToInt(state string) (int, error) {
	state = strings.ToUpper(state)
	switch state {
	case "OK", "0":
		return 0, nil
	case "WARNING", "1":
		return 1, nil
	case "CRITICAL", "2":
		return 2, nil
	case "UNKNOWN", "3":
		return 3, nil
	default:
		return 0, errors.New("invalid state")
	}
}

func init() {
	rootCmd.AddCommand(snapshotCmd)

	fs := snapshotCmd.Flags()

	fs.StringP("snapshot", "s", "*",
		"Comma-separated list of snapshot names to retrieve. Wildcard (*) expressions are supported")
	fs.StringP("repository", "r", "*",
		"Comma-separated list of snapshot repository names used to limit the request")

	fs.IntP("number", "N", 1, "Check latest N number snapshots. If not set only the latest snapshot is checked")
	fs.BoolP("all", "a", false, "Check all retrieved snapshots. If not set only the latest snapshot is checked")

	fs.StringP("no-snapshots-state", "T", "UNKNOWN", "State to assign when no snapshots are found (0, 1, 2, 3, OK, WARNING, CRITICAL, UNKNOWN). If not set this defaults to UNKNOWN")

	snapshotCmd.MarkFlagsMutuallyExclusive("number", "all")
}
