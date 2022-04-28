/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// clientsCmd represents the clients command.
var clientsCmd = &cobra.Command{
	Use:   "clients",
	Short: "Shows information about the clients currently connected to Rosenbridge.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("clients called")
	},
}

func init() {
	rootCmd.AddCommand(clientsCmd)
}
