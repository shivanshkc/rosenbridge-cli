package cmd

import (
	"context"
	"fmt"

	"github.com/shivanshkc/rosenbridge-cli/lib"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// connectClientID binds with the client ID flag of the connect command.
var connectClientID string

// connectCmd represents the connect command.
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Establishes connection with Rosenbridge and starts streaming messages.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Validating the client ID.
		if err := checkClientID(connectClientID); err != nil {
			exitWithPrintf(1, err.Error())
		}

		// Getting a new connection to Rosenbridge.
		conn, err := lib.NewConnection(context.Background(), &lib.ConnectionParams{
			ClientID:     connectClientID,
			BaseURL:      viper.GetString("backend.base_url"),
			IsTLSEnabled: viper.GetBool("backend.is_tls_enabled"),
		})
		if err != nil {
			exitWithPrintf(1, "Failed to connect: %s", err.Error())
		}
		color.Green("Connected with Rosenbridge.\n")

		// Printing all incoming messages.
		conn.IncomingMessageHandler = printMessage

		// TODO: Just testing.
		if err := conn.SendMessageAsync(context.Background(), &lib.OutgoingMessage{
			RequestID:   uuid.NewString(),
			ReceiverIDs: []string{"sk"},
			Message:     "automated message",
			Persist:     "true",
		}); err != nil {
			panic("send message async error: " + err.Error())
		}

		// Blocking forever. TODO: Replace this with an interruption listener.
		select {}
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)

	// Setting up the --client-id or -c flag.
	connectCmd.Flags().StringVarP(&connectClientID, "client-id", "c", "",
		"ID of the client making the connection.")

	// The --client-id flag is required.
	if err := connectCmd.MarkFlagRequired("client-id"); err != nil {
		panic(fmt.Errorf("failed to mark client-id flag as required: %w", err))
	}
}
