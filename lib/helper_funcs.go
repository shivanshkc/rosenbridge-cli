package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

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

// createWebsocketURL converts "https" URL -> "wss" URL, and "http" URL -> "ws" URL.
func createWebsocketURL(baseURL string) string {
	// If the URL is already a websocket URL, there's nothing to do.
	if strings.HasPrefix(baseURL, "wss://") || strings.HasPrefix(baseURL, "ws://") {
		return baseURL
	}

	// Handling https URLs.
	if strings.HasPrefix(baseURL, "https://") {
		return "wss://" + strings.TrimPrefix(baseURL, "https://")
	}

	// Handling http URLs.
	if strings.HasPrefix(baseURL, "http://") {
		return "ws://" + strings.TrimPrefix(baseURL, "http://")
	}

	// Should never be the case because of the validations.
	return baseURL
}

// createHTTPURL converts "wss" URL -> "https" URL, and "ws" URL -> "http" URL.
func createHTTPURL(baseURL string) string {
	// If the URL is already an HTTP URL, there's nothing to do.
	if strings.HasPrefix(baseURL, "https://") || strings.HasPrefix(baseURL, "http://") {
		return baseURL
	}

	// Handling wss URLs.
	if strings.HasPrefix(baseURL, "wss://") {
		return "https://" + strings.TrimPrefix(baseURL, "wss://")
	}

	// Handling http URLs.
	if strings.HasPrefix(baseURL, "ws://") {
		return "http://" + strings.TrimPrefix(baseURL, "ws://")
	}

	// Should never be the case because of the validations.
	return baseURL
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

// isCode2xx tells if the given status code is 2xx.
func isCode2xx(statusCode int) bool {
	return statusCode/100 == 2
}
