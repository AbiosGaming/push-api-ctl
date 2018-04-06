package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var clientIDFlag string
var clientSecretFlag string
var apiBaseURLFlag string
var wsBaseURLFlag string

var rootCmd = &cobra.Command{
	SilenceUsage:  true, // Don't print usage on errors that propagated to the top
	SilenceErrors: true, // Don't let cobra print the error message, we'll print it ourselves in Execute
	Use:           "push-api-ctl",
	Short:         "Abios Push Service Management",
	Long:          `Command line tool for managing the push api`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stderr, "%s\n", cmd.UsageString())
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiBaseURLFlag, "api-url", "https://api.abiosgaming.com/v2", "Root URL for access token generation")
	rootCmd.PersistentFlags().StringVar(&wsBaseURLFlag, "push-service-url", "https://ws.abiosgaming.com/v0", "Root URL for the push service")
	rootCmd.PersistentFlags().StringVar(&clientIDFlag, "client-id", "", "Client id (required)")
	rootCmd.PersistentFlags().StringVar(&clientSecretFlag, "client-secret", "", "Client secret (required)")

	rootCmd.MarkPersistentFlagRequired("client-id")
	rootCmd.MarkPersistentFlagRequired("client-secret")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
