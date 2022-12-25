package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

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

		// Creating connection params for sending messages.
		params := &lib.ConnectionParams{
			ClientID:     sendSenderID,
			BaseURL:      viper.GetString("backend.base_url"),
			IsTLSEnabled: viper.GetBool("backend.is_tls_enabled"),
		}

		// Creating a reader to read typed messages from stdin.
		reader := bufio.NewReader(os.Stdin)
		for {
			// Prompt.
			fmt.Printf(">> You: ") //nolint:forbidigo

			// Reading the input.
			messageBody, err := reader.ReadString('\n')
			if err != nil {
				color.Red("Error while reading message: %s\n", err.Error())
				return
			}

			// Remove trailing newline char.
			messageBody = strings.TrimSuffix(messageBody, "\n")
			// Forming the exact outgoing message.
			outgoingMessage := &lib.OutgoingMessageReq{
				RequestID:   uuid.NewString(),
				ReceiverIDs: []string{sendReceiverID},
				Message:     messageBody,
				SenderID:    params.ClientID,
			}

			// Sending the message whilst handling Cloud Run errors.
			if err := sendMessageWithColdStartHandling(outgoingMessage, params); err != nil {
				break
			}
		}
	},
}

// sendMessageWithColdStartHandling sends the given message using the given connection params.
// It also handles GCP Cloud Run's annoying 429 errors.
func sendMessageWithColdStartHandling(outMessage *lib.OutgoingMessageReq, params *lib.ConnectionParams) error {
	// Number of max retries.
	retryCount := viper.GetInt("general.cold_start_retry_count")

	// We only print cold-restart warning once, so a flag is required to keep track.
	var isWarningPrinted bool

	// Starting the retry loop to deal with GCP cold-start errors.
	for i := 0; i < retryCount; i++ {
		// Sending the message.
		_, err := lib.SendMessage(context.Background(), outMessage, params)
		if err == nil {
			return nil
		}

		// This will be logged upon every failure.
		color.Red("Error while sending message: %s\n", err.Error())

		// If the error is 429, we attempt a retry.
		if errors.Is(err, lib.ErrTooManyReq) {
			if !isWarningPrinted {
				color.Yellow("Looks like the server is under load. Retrying...")
			}

			isWarningPrinted = true
			time.Sleep(time.Second)
		}
	}

	// Retries didn't work.
	color.Red("The server is busy. Please try again in some time.")
	return errors.New("failed to send message")
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
