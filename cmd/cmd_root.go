package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// baseURL variable for binding with a persistent flag.
var baseURL string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "rosen",
	Short: "A distributed hub for real-time communication between servers.",
	Long:  "",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Declaring the optional "url" flag.
	rootCmd.PersistentFlags().StringVarP(&baseURL, "url", "u", "wss://rosenbridge.rudraksh.dev",
		"URL of any one node of a Rosenbridge cluster.")
}
