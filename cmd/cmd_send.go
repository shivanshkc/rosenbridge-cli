package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/shivanshkc/rosenbridge-cli/lib"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// These variables bind with the flags of the send command.
var sendSenderID, sendReceiverID string

// sendCmd represents the send command.
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Opens a console for writing messages to the intended client.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Validating the IDs.
		if err := checkClientID(sendSenderID); err != nil {
			exitWithPrintf(1, err.Error())
		}
		if err := checkClientID(sendReceiverID); err != nil {
			exitWithPrintf(1, err.Error())
		}

		// Getting a new connection to Rosenbridge.
		conn, err := lib.NewConnection(context.Background(), &lib.ConnectionParams{
			ClientID:     sendSenderID,
			BaseURL:      viper.GetString("backend.base_url"),
			IsTLSEnabled: viper.GetBool("backend.is_tls_enabled"),
		})
		if err != nil {
			exitWithPrintf(1, "Failed to connect: %s", err.Error())
		}
		color.Green("Connected with Rosenbridge.\n")
		// Connection will be closed upon function return.
		defer func() {
			_ = conn.Close()
			color.Green("Disconnected from Rosenbridge.\n")
		}()

		// Creating a reader to read typed messages from stdin.
		reader := bufio.NewReader(os.Stdin)
		for {
			// Prompt.
			fmt.Printf(">> You: ")

			// Reading the input.
			messageBody, err := reader.ReadString('\n')
			if err != nil {
				color.Red("Error while reading message: %s\n", err.Error())
				return
			}

			// Forming the exact outgoing message.
			outgoingMessage := &lib.OutgoingMessage{
				RequestID:   uuid.NewString(),
				ReceiverIDs: []string{sendReceiverID},
				Message:     messageBody,
				Persist:     lib.PersistTrue,
			}

			// Sending the message.
			if err := conn.SendMessageAsync(context.Background(), outgoingMessage); err != nil {
				color.Red("Error while sending message: %s\n", err.Error())
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	// Setting up the --sender or -s flag.
	sendCmd.Flags().StringVarP(&sendSenderID, "sender", "s", "",
		"ID of the client sending the message(s).")

	// The --sender flag is required.
	if err := sendCmd.MarkFlagRequired("sender"); err != nil {
		panic(fmt.Errorf("failed to mark sender flag as required: %w", err))
	}

	// Setting up the --receiver or -r flag.
	sendCmd.Flags().StringVarP(&sendReceiverID, "receiver", "r", "",
		"ID of the client receiving the message(s).")

	// The --receiver flag is required.
	if err := sendCmd.MarkFlagRequired("receiver"); err != nil {
		panic(fmt.Errorf("failed to mark receiver flag as required: %w", err))
	}
}
