package cmd

import (
	"github.com/NETWAYS/go-check"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "check_elasticsearch",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: Help,
}

func Execute(version string) {
	defer check.CatchPanic()

	rootCmd.Version = version
	rootCmd.VersionTemplate()

	if err := rootCmd.Execute(); err != nil {
		check.ExitError(err)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.DisableAutoGenTag = true

	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	pfs := rootCmd.PersistentFlags()
	pfs.StringVarP(&cliConfig.Hostname, "hostname", "H", "localhost", "")
	pfs.StringVarP(&cliConfig.Username, "username", "U", "", "")
	pfs.StringVarP(&cliConfig.Password, "password", "P", "", "")
	pfs.IntVarP(&cliConfig.Port, "port", "p", 9200, "")
	pfs.BoolVarP(&cliConfig.TLS, "tls", "S", false, "")
	pfs.BoolVar(&cliConfig.Insecure, "insecure", false, "")

	pfs.SortFlags = false
}

func Help(cmd *cobra.Command, strings []string) {
	_ = cmd.Usage()

	os.Exit(3)
}
