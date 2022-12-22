package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Connection represents a connection with Rosenbridge.
type Connection struct {
	// underlyingConn is the low-level connection object.
	underlyingConn *websocket.Conn
	// connectionParams are the parameters required to create the connection.
	connectionParams *ConnectionParams

	// IncomingMessageHandler handles incoming message.
	IncomingMessageHandler IncomingMessageHandlerFunc
	// OutgoingMessageResponseHandler handles outgoing message responses.
	OutgoingMessageResponseHandler OutgoingMessageResponseHandlerFunc
	// ConnectionClosureHandler handles connection closures.
	ConnectionClosureHandler ConnectionClosureHandlerFunc
}

// NewConnection creates and returns a new connection.
func NewConnection(ctx context.Context, params *ConnectionParams) (*Connection, error) {
	// Deciding on the protocol.
	wsProtocol := getWebsocketProtocol(params)
	// Forming the API endpoint URL.
	endpoint := fmt.Sprintf("%s://%s/api/bridge", wsProtocol, params.BaseURL)

	// Request headers.
	headers := &http.Header{}
	headers.Set("x-client-id", params.ClientID)

	// Establishing websocket connection.
	underlyingConn, response, err := websocket.DefaultDialer.Dial(endpoint, *headers)
	if err != nil {
		return nil, fmt.Errorf("error in websocket.Dial: %w", err)
	}
	defer func() { _ = response.Body.Close() }()

	// Creating the connection abstraction.
	conn := &Connection{
		underlyingConn:                 underlyingConn,
		connectionParams:               params,
		IncomingMessageHandler:         DefaultIncomingMessageHandler,
		OutgoingMessageResponseHandler: DefaultOutgoingMessageResponseHandler,
		ConnectionClosureHandler:       DefaultConnectionClosureHandler,
	}

	// Starting a separate goroutine to listen to websocket messages.
	go websocketMessageReader(ctx, conn)
	return conn, nil
}

// SendMessage sends a new message synchronously.
// It is a stateless way to send a message and hence does not need to be associated to a connection.
func SendMessage(ctx context.Context, request *OutgoingMessageReq, params *ConnectionParams) (
	*OutgoingMessageRes, error,
) {
	request.SenderID = params.ClientID
	// Marshalling the request to byte array.
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	// Converting the request byte array to io.Reader for the http client.
	bodyReader := bytes.NewReader(requestBytes)

	// Deciding on the protocol.
	httpProtocol := getHTTPProtocol(params)
	// Forming the endpoint.
	endpoint := fmt.Sprintf("%s://%s/api/message", httpProtocol, params.BaseURL)

	// Forming the HTTP request.
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to form the http request: %w", err)
	}
	httpRequest.Header.Set("x-request-id", request.RequestID)

	// Executing the request.
	response, err := (&http.Client{}).Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http request: %w", err)
	}
	defer func() { _ = response.Body.Close() }()

	// Decoding the response body.
	outMessageRes := &OutgoingMessageRes{}
	if err := anyToAny(response.Body, outMessageRes); err != nil {
		return nil, fmt.Errorf("failed to get response body: %w", err)
	}

	// If the request failed completely, we create the error from the custom code of the response.
	if outMessageRes.Code != codeOK {
		return nil, fmt.Errorf("request failed: %s", outMessageRes.Reason)
	}

	outMessageRes.RequestID = response.Header.Get("x-request-id")
	return outMessageRes, nil
}

// SendMessageAsync sends a new message asynchronously.
// It uses the websocket connection for sending the message.
// The response of this request can be handled through the ResponseHandler function.
func (c *Connection) SendMessageAsync(ctx context.Context, request *OutgoingMessageReq) error {
	message := &BridgeMessage{
		Type:      typeOutgoingMessageReq,
		RequestID: request.RequestID,
		Body:      request,
	}

	// Marshalling the message to byte array.
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Writing the message to the connection.
	if err := c.underlyingConn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	return nil
}

// Close closes the underlying connection.
func (c *Connection) Close() error {
	if err := c.underlyingConn.Close(); err != nil {
		return fmt.Errorf("failed to close underlying conn: %w", err)
	}
	return nil
}

// websocketMessageReader manages the websocket connection and messages by calling appropriate handlers.
//
//nolint:cyclop
func websocketMessageReader(ctx context.Context, conn *Connection) {
	// This routine returns when the connection closes.
	defer conn.ConnectionClosureHandler(ctx, recover())

	// Starting an infinite loop to process all websocket communication.
	// This loop panics when the connection is closed.
	for {
		wsMessageType, message, err := conn.underlyingConn.ReadMessage()
		if err != nil {
			// This invokes the ClosureHandler with the given error.
			panic(fmt.Errorf("error in ReadMessage: %w", err))
		}

		// Handling different websocket message types.
		switch wsMessageType {
		case websocket.CloseMessage:
			// This closes the connection with nil error.
			panic(nil)
		case websocket.TextMessage:
			bridgeMessage := &BridgeMessage{}
			if err := anyToAny(message, bridgeMessage); err != nil {
				// If the message type fails to be determined, we assume it to be an incoming message.
				conn.IncomingMessageHandler(ctx, nil, fmt.Errorf("failed to decode message: %w", err))
				continue
			}

			// Handling different message types.
			switch bridgeMessage.Type {
			case typeIncomingMessageReq:
				inMessageReq := &IncomingMessageReq{}
				if err := anyToAny(bridgeMessage.Body, inMessageReq); err != nil {
					conn.IncomingMessageHandler(ctx, nil,
						fmt.Errorf("failed to unmarshal message: %w", err))
					continue
				}
				conn.IncomingMessageHandler(ctx, inMessageReq, nil)
			case typeOutgoingMessageRes:
				outMessageRes := &OutgoingMessageRes{}
				if err := anyToAny(bridgeMessage.Body, outMessageRes); err != nil {
					conn.OutgoingMessageResponseHandler(ctx, nil,
						fmt.Errorf("failed to unmarshal message: %w", err))
					continue
				}
				conn.OutgoingMessageResponseHandler(ctx, outMessageRes, nil)
			case typeErrorRes:
				// If the response type is error, we assume it to be an incoming message.
				conn.IncomingMessageHandler(ctx, nil, errors.New("unknown error"))
			default:
				// Unknown message types are simply ignored.
			}
		case websocket.BinaryMessage:
		case websocket.PingMessage:
		case websocket.PongMessage:
		default:
		}
	}
}
