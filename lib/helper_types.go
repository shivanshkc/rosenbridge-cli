package lib

import (
	"context"
)

// ConnectionParams are the params required to create the connection.
type ConnectionParams struct {
	// ClientID is the ID to which the connection belongs.
	ClientID string
	// BaseURL is the URL of the Rosenbridge deployment without protocol.
	BaseURL string
	// IsTLSEnabled is a flag to tell if the deployment is TLS enabled.
	// If it is true, connection is attempted with "wss" protocol, otherwise "ws" is used.
	IsTLSEnabled bool
}

// PersistenceCriteria is an enum for the persistence criteria Rosenbridge supports.
type PersistenceCriteria string

// IncomingMessage is the schema of an incoming message from Rosenbridge.
type IncomingMessage struct {
	// RequestID is the unique identifier of the request that sent this message originally.
	RequestID string `json:"request_id"`
	// SenderID is the ID of the client who sent this message.
	SenderID string `json:"sender_id"`
	// Message is the main message content.
	Message string `json:"message"`
	// Persist is the persistence criteria of this message, set by the sender.
	Persist PersistenceCriteria `json:"persist"`

	// Type field is used internally and should not be populated by the caller.
	Type string `json:"type"`
}

// OutgoingMessage is the schema of an outgoing message on Rosenbridge.
type OutgoingMessage struct {
	// RequestID is the unique identifier of this request.
	RequestID string `json:"request_id"`
	// ReceiverIDs is the list of receiver client IDs that are intended to receive this message.
	ReceiverIDs []string `json:"receiver_ids"`
	// Message is the main message content.
	Message string `json:"message"`
	// Persist is the persistence criteria of the message.
	Persist PersistenceCriteria `json:"persist"`

	// Type field is used internally and should not be populated by the caller.
	Type string `json:"type"`
}

// OutgoingMessageResponse is the response of sending a message.
// It tells which of the clients received the message, and which ones did not, along with the reasons.
type OutgoingMessageResponse struct {
	// RequestID is the unique identifier of the request to which this response belongs.
	RequestID string `json:"request_id"`
	// Results holds the message delivery status for each receiver.
	Results []*struct {
		// ReceiverID is the ID of the receiver.
		ReceiverID string `json:"receiver_id"`
		// Code tells the final status of message delivery.
		Code string `json:"code"`
		// Reason tells why the delivery failed (if it failed).
		Reason string `json:"reason"`
	} `json:"results"`

	// Type field is used internally and should not be populated by the caller.
	Type string `json:"type"`
}

// IncomingMessageHandlerFunc is the type of func that handles incoming messages.
// The error parameter notifies the caller of any errors that might occur while receiving/decoding the message.
//
// Note that if any error occurs before the type of the message itself could be determined, the message will be assumed
// as an IncomingMessage, and so the IncomingMessageHandlerFunc will be invoked.
type IncomingMessageHandlerFunc func(ctx context.Context, message *IncomingMessage, err error)

// OutgoingMessageResponseHandlerFunc is the type of func that handles outgoing message responses.
// The error parameter notifies the caller of any errors that might occur while receiving/decoding the message.
//
// Note that if any error occurs before the type of the message itself could be determined, the message will be assumed
// as an IncomingMessage, and so the IncomingMessageHandlerFunc will be invoked.
type OutgoingMessageResponseHandlerFunc func(ctx context.Context, response *OutgoingMessageResponse, err error)

// ConnectionClosureHandlerFunc is the type of func that handles connection closures.
// The error parameter gives info on why the connection closed.
type ConnectionClosureHandlerFunc func(ctx context.Context, err interface{})

// httpResponseBody is the schema of a response body received from Rosenbridge.
type httpResponseBody struct {
	StatusCode int         `json:"status_code"`
	CustomCode string      `json:"custom_code"`
	Data       interface{} `json:"data"`
}
