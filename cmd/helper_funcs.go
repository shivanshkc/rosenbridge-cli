package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/shivanshkc/rosenbridge-cli/models"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

// getWebsocketConn establishes the websocket connection as per the configs and returns it.
func getWebsocketConn(clientID string) (*websocket.Conn, error) {
	// Forming the API endpoint URL.
	endpoint := fmt.Sprintf("ws://%s/api/clients/%s/connection", viper.GetString("backend.addr"), clientID)
	// Establishing websocket connection.
	conn, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error in websocket.Dial: %w", err)
	}
	return conn, nil
}

// exitWithPrintf prints the provided message in Printf style and then calls os.Exit with provided code.
func exitWithPrintf(code int, format string, a ...interface{}) {
	switch code {
	case 0:
		color.Green(format+"\n", a...)
	default:
		color.Red(format+"\n", a...)
	}
	os.Exit(code)
}

// unmarshalReceivedMessage decodes the received message bytes into the intended struct.
func unmarshalReceivedMessage(message []byte) (*models.ReceivedMessage, error) {
	messageStruct := &models.ReceivedMessage{}
	if err := json.Unmarshal(message, messageStruct); err != nil {
		return nil, fmt.Errorf("failed to unmarshal received message: %w", err)
	}
	return messageStruct, nil
}

// printReceivedMessage prints the provided message in appropriate format.
// If the provided message is nil, it prints that the message is ill-formatted.
func printReceivedMessage(message *models.ReceivedMessage) {
	if message == nil {
		color.Yellow(">> [%s] Message is of unrecognized format.\n", time.Now().Format(time.Kitchen))
		return
	}
	color.Yellow(">> [%s] %s: %s\n", time.Now().Format(time.Kitchen), message.SenderID, message.MessageBody)
}
