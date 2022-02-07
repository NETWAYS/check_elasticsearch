package cmd

import (
	"github.com/NETWAYS/go-check"
	"github.com/spf13/cobra"
	"os"
)

var (
	Timeout = 30
)

var rootCmd = &cobra.Command{
	Use:   "check_elasticsearch",
	Short: "Icinga check plugin to check Elasticsearch",
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		go check.HandleTimeout(Timeout)
	},
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
	pfs.StringVarP(&cliConfig.Hostname, "hostname", "H", "localhost",
		"Address of elasticsearch node")
	pfs.IntVarP(&cliConfig.Port, "port", "p", 9200,
		"Port of elasticsearch node")
	pfs.StringVarP(&cliConfig.Username, "username", "U", "",
		"Username if authentication is required")
	pfs.StringVarP(&cliConfig.Password, "password", "P", "",
		"Password if authentication is required")
	pfs.BoolVarP(&cliConfig.TLS, "tls", "S", false,
		"Use secure connection")
	pfs.BoolVar(&cliConfig.Insecure, "insecure", false,
		"Allow use of self signed certificates when using SSL")
	pfs.IntVarP(&Timeout, "timeout", "t", Timeout,
		"Timeout for the check")

	rootCmd.Flags().SortFlags = false
	pfs.SortFlags = false
}

func Help(cmd *cobra.Command, strings []string) {
	_ = cmd.Usage()

	os.Exit(3)
}
