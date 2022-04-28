/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// sendCmd represents the send command.
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Opens a console for writing messages to the intended client.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("send called")
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
