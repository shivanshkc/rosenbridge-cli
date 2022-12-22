package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/fatih/color"
)

// DefaultIncomingMessageHandler is the default handler for incoming messages.
func DefaultIncomingMessageHandler(ctx context.Context, message *IncomingMessageReq, err error) {}

// DefaultOutgoingMessageResponseHandler is the default handler for outgoing message responses.
func DefaultOutgoingMessageResponseHandler(ctx context.Context, response *OutgoingMessageRes, err error) {
}

// DefaultConnectionClosureHandler is the default handler for connection closures.
func DefaultConnectionClosureHandler(ctx context.Context, err interface{}) {
	if err != nil {
		color.Red("Connection closed with error: %v", err)
	}
}

// getWebsocketProtocol provides the correct websocket protocol based on the connection params.
func getWebsocketProtocol(params *ConnectionParams) string {
	if params.IsTLSEnabled {
		return "wss"
	}
	return "ws"
}

// getHTTPProtocol provides the correct http protocol based on the connection params.
func getHTTPProtocol(params *ConnectionParams) string {
	if params.IsTLSEnabled {
		return "https"
	}
	return "http"
}

// anyToBytes converts the provided input to a byte slice.
//
// If the conversion is not possible, it returns a non-nil error.
func anyToBytes(input interface{}) ([]byte, error) {
	switch asserted := input.(type) {
	case []byte:
		return asserted, nil
	case string:
		return []byte(asserted), nil
	case io.Reader:
		// Reading all the data.
		inputBytes, err := io.ReadAll(asserted)
		if err != nil {
			return nil, fmt.Errorf("error in io.ReadAll call: %w", err)
		}
		// Conversion successful.
		return inputBytes, nil
	default:
		// Marshalling to JSON. This works with all primitive data types and structs etc.
		inputBytes, err := json.Marshal(input)
		if err != nil {
			return nil, fmt.Errorf("error in json.Marshal call: %w", err)
		}
		// Conversion successful.
		return inputBytes, nil
	}
}

// anyToAny marshals the provided input and then un-marshals it into the provided output.
func anyToAny(input interface{}, targetOutput interface{}) error {
	// Marshalling the input.
	inputBytes, err := anyToBytes(input)
	if err != nil {
		return fmt.Errorf("error in anyToBytes call: %w", err)
	}

	// Unmarshalling into the target.
	if err := json.Unmarshal(inputBytes, targetOutput); err != nil {
		return fmt.Errorf("error in json.Unmarshal call: %w", err)
	}

	return nil
}
