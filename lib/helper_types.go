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

// BridgeMessage is the general schema of all messages that are sent over a bridge.
type BridgeMessage struct {
	// Type of the message. It can be used to differentiate and route various kinds of messages.
	Type string `json:"type"`
	// RequestID is used to correlate an outgoing-message-request with its corresponding response.
	RequestID string `json:"request_id"`
	// Body is the main content of this message.
	Body interface{} `json:"body"`
}

// IncomingMessageReq is the schema of an incoming message from Rosenbridge.
type IncomingMessageReq struct {
	// SenderID is the ID of the client who sent the message.
	SenderID string `json:"sender_id"`
	// Message is the main message content.
	Message string `json:"message"`
}

// OutgoingMessageReq is the schema of an outgoing message on Rosenbridge.
type OutgoingMessageReq struct {
	// SenderID is the ID of client who sent this message.
	SenderID string `json:"sender_id"`
	// Receivers is the list of client IDs that are intended to receive this message.
	ReceiverIDs []string `json:"receiver_ids"`
	// Message is the main message content that needs to be delivered.
	Message string `json:"message"`

	RequestID string `json:"-"`
}

// OutgoingMessageRes is the response of sending a message.
// It tells which of the clients received the message, and which ones did not, along with the reasons.
type OutgoingMessageRes struct {
	// Code is OK if the request is processable.
	// If it is negative, it means the message delivery was not even attempted.
	Code string `json:"code"`
	// Reason tells why the request is not processable (if it's not).
	Reason string `json:"reason"`
	// Report holds the message delivery status for each receiver.
	Report map[string][]*struct {
		// ClientID is the ID of the client to whom the bridge belongs.
		ClientID string `json:"client_id,omitempty"`
		// BridgeID is the unique ID of the bridge.
		BridgeID string `json:"bridge_id,omitempty"`
		// Code tells the final status of message delivery.
		Code string `json:"code"`
		// Reason tells why the delivery failed (if it failed).
		Reason string `json:"reason"`
	} `json:"report"`

	RequestID string `json:"-"`
}

// IncomingMessageHandlerFunc is the type of func that handles incoming messages.
// The error parameter notifies the caller of any errors that might occur while receiving/decoding the message.
//
// Note that if any error occurs before the type of the message itself could be determined, the message will be assumed
// as an IncomingMessageReq, and so the IncomingMessageHandlerFunc will be invoked.
type IncomingMessageHandlerFunc func(ctx context.Context, message *IncomingMessageReq, err error)

// OutgoingMessageResponseHandlerFunc is the type of func that handles outgoing message responses.
// The error parameter notifies the caller of any errors that might occur while receiving/decoding the message.
//
// Note that if any error occurs before the type of the message itself could be determined, the message will be assumed
// as an IncomingMessageReq, and so the IncomingMessageHandlerFunc will be invoked.
type OutgoingMessageResponseHandlerFunc func(ctx context.Context, response *OutgoingMessageRes, err error)

// ConnectionClosureHandlerFunc is the type of func that handles connection closures.
// The error parameter gives info on why the connection closed.
type ConnectionClosureHandlerFunc func(ctx context.Context, err interface{})
