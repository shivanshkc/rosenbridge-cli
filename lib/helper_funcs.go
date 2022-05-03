package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// DefaultIncomingMessageHandler is the default handler for incoming messages.
func DefaultIncomingMessageHandler(ctx context.Context, message *IncomingMessage, err error) {}

// DefaultOutgoingMessageResponseHandler is the default handler for outgoing message responses.
func DefaultOutgoingMessageResponseHandler(ctx context.Context, response *OutgoingMessageResponse, err error) {
}

// DefaultConnectionClosureHandler is the default handler for connection closures.
func DefaultConnectionClosureHandler(ctx context.Context, err interface{}) {}

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

// unmarshalMessageType provides the type of the message.
func unmarshalMessageType(message []byte) (string, error) {
	// Decoding into a simple map.
	decoded := map[string]interface{}{}
	if err := json.Unmarshal(message, &decoded); err != nil {
		return "", fmt.Errorf("failed to unmarshal message: %w", err)
	}
	// Checking if there's a type key.
	mType, exists := decoded["type"]
	if !exists {
		return "", fmt.Errorf("no message type")
	}
	// Checking if the message type is string.
	mTypeString, asserted := mType.(string)
	if !asserted {
		return "", fmt.Errorf("invalid message type: %v", mType)
	}
	return mTypeString, nil
}

// unmarshalIncomingMessage decodes the provided byte slice into an IncomingMessage.
func unmarshalIncomingMessage(message []byte) (*IncomingMessage, error) {
	inMessage := &IncomingMessage{}
	if err := json.Unmarshal(message, inMessage); err != nil {
		return nil, fmt.Errorf("error in json.Unmarshal call: %w", err)
	}
	return inMessage, nil
}

// unmarshalOutgoingMessageResponse decodes the provided byte slice into an OutgoingMessageResponse.
func unmarshalOutgoingMessageResponse(message []byte) (*OutgoingMessageResponse, error) {
	outMessageResp := &OutgoingMessageResponse{}
	if err := json.Unmarshal(message, outMessageResp); err != nil {
		return nil, fmt.Errorf("error in json.Unmarshal call: %w", err)
	}
	return outMessageResp, nil
}

// unmarshalHTTPResponse decodes a http response body into its struct.
func unmarshalHTTPResponse(response *http.Response) (*httpResponseBody, error) {
	// Reading into a byte slice.
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshalling into the struct.
	responseBody := &httpResponseBody{}
	if err := json.Unmarshal(bodyBytes, responseBody); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return responseBody, nil
}

// interface2OutgoingMessageResponse converts the provided interface into an *OutgoingMessageResponse.
// Any conversion errors are reported by the error parameter.
func interface2OutgoingMessageResponse(input interface{}) (*OutgoingMessageResponse, error) {
	// Marshalling into json for later unmarshalling into struct.
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	// Unmarshalling into struct.
	messageResponse := &OutgoingMessageResponse{}
	if err := json.Unmarshal(inputBytes, messageResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal input into struct: %w", err)
	}

	return messageResponse, nil
}

// isPositiveStatusCode tells if the provided http status code implies successful operation.
func isPositiveStatusCode(code int) bool {
	return code < 300
}
