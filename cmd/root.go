package cmd

import (
	"os"

	"github.com/NETWAYS/go-check"
	"github.com/spf13/cobra"
)

var (
	timeout = 30
)

var rootCmd = &cobra.Command{
	Use:   "check_elasticsearch",
	Short: "Icinga check plugin to check Elasticsearch",
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		go check.HandleTimeout(timeout)
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
	pfs.StringArrayVarP(&cliConfig.Hostname, "hostname", "H", []string{"http://localhost:9200"},
		"URL of an Elasticsearch instance. Can be used multiple times.")
	pfs.StringVarP(&cliConfig.Username, "username", "U", "",
		"Username for HTTP Basic Authentication (CHECK_ELASTICSEARCH_USERNAME)")
	pfs.StringVarP(&cliConfig.Password, "password", "P", "",
		"Password for HTTP Basic Authentication (CHECK_ELASTICSEARCH_PASSWORD)")
	pfs.StringVarP(&cliConfig.Bearer, "bearer", "b", "",
		"Specify the Bearer Token for authentication (CHECK_ELASTICSEARCH_BEARER)")
	pfs.BoolVar(&cliConfig.Insecure, "insecure", false,
		"Skip the verification of the server's TLS certificate")
	pfs.StringVarP(&cliConfig.CAFile, "ca-file", "", "",
		"Specify the CA File for TLS authentication (CHECK_ELASTICSEARCH_CA_FILE)")
	pfs.StringVarP(&cliConfig.CertFile, "cert-file", "", "",
		"Specify the Certificate File for TLS authentication (CHECK_ELASTICSEARCH_CERT_FILE)")
	pfs.StringVarP(&cliConfig.KeyFile, "key-file", "", "",
		"Specify the Key File for TLS authentication (CHECK_ELASTICSEARCH_KEY_FILE)")
	pfs.IntVarP(&timeout, "timeout", "t", timeout,
		"Timeout in seconds for the plugin")

	rootCmd.Flags().SortFlags = false
	pfs.SortFlags = false

	loadFromEnv(&cliConfig)
}

func Help(cmd *cobra.Command, _ []string) {
	_ = cmd.Usage()

	os.Exit(3)
}
