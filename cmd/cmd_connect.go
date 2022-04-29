package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

// connectClientID binds with the client ID flag of the connect command.
var connectClientID string

// connectCmd represents the connect command.
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Establishes connection with Rosenbridge and starts streaming messages.",
	Long:  `Usage: rosen connect <client-id>`,
	Run: func(cmd *cobra.Command, args []string) {
		// Making sure the user provided a client ID.
		if connectClientID == "" {
			exitWithPrintf(1, cmd.Long)
		}

		// Validating the client ID.
		if err := checkClientID(connectClientID); err != nil {
			exitWithPrintf(1, err.Error())
		}

		// Establishing connection with Rosenbridge.
		conn, err := getWebsocketConn(connectClientID)
		if err != nil {
			exitWithPrintf(1, "Failed to connect: %s", err.Error())
		}
		color.Green("Connected with Rosenbridge.\n")
		// Connection will be closed upon function return.
		defer func() {
			_ = conn.Close()
			color.Green("Disconnected from Rosenbridge.\n")
		}()

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				// Not using exitWithPrintf here because of the defer-block above.
				color.Red("Error while reading messages: %s\n", err.Error())
				return
			}

			// Handling different message types.
			switch messageType {
			case websocket.CloseMessage:
				// Ending the loop upon connection closure.
				return
			case websocket.TextMessage:
				decodedMessage, _ := unmarshalReceivedMessage(message)
				printReceivedMessage(decodedMessage)
			case websocket.BinaryMessage:
			case websocket.PingMessage:
			case websocket.PongMessage:
			default:
			}
		}
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
