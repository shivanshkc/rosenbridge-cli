/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// messagesCmd represents the messages command.
var messagesCmd = &cobra.Command{
	Use:   "messages",
	Short: "Lists all messages that are persisted in the Rosenbridge database.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("messages called")
	},
}

func init() {
	rootCmd.AddCommand(messagesCmd)
}
