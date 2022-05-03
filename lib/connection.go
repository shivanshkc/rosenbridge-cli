package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Connection represents a connection with Rosenbridge.
type Connection struct {
	// underlyingConn is the low-level connection object.
	underlyingConn *websocket.Conn
	// httpClient is the client for making HTTP requests to Rosenbridge.
	httpClient *http.Client
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
	endpoint := fmt.Sprintf("%s://%s/clients/%s/connection", wsProtocol, params.BaseURL, params.ClientID)

	// Establishing websocket connection.
	underlyingConn, response, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error in websocket.Dial: %w", err)
	}
	defer func() { _ = response.Body.Close() }()

	// Creating the HTTP client for synchronous requests.
	httpClient := &http.Client{}

	// Creating the connection abstraction.
	conn := &Connection{
		underlyingConn:                 underlyingConn,
		httpClient:                     httpClient,
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
// It uses an HTTP request to send the message.
func (c *Connection) SendMessage(ctx context.Context, request *OutgoingMessage) (*OutgoingMessageResponse, error) {
	// Setting the message type.
	request.Type = typeOutgoingMessage

	// Marshalling request to byte array.
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Converting the request byte array to io.Reader for the http client.
	bodyReader := bytes.NewReader(requestBytes)

	// Deciding on the protocol.
	httpProtocol := getHTTPProtocol(c.connectionParams)
	// Forming the endpoint.
	endpoint := fmt.Sprintf("%s://%s/clients/%s/messages", httpProtocol, c.connectionParams.BaseURL,
		c.connectionParams.ClientID)

	// Forming the HTTP request.
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to form the http request: %w", err)
	}

	// Executing the request.
	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http request: %w", err)
	}
	defer func() { _ = response.Body.Close() }()

	// Getting the response body.
	responseBody, err := unmarshalHTTPResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to get response body: %w", err)
	}

	// If the request failed completely, we create the error from the custom code of the response.
	if !isPositiveStatusCode(response.StatusCode) {
		return nil, fmt.Errorf("request failed: %s", responseBody.CustomCode)
	}

	// Converting the generic http response data into required struct.
	messageResponse, err := interface2OutgoingMessageResponse(responseBody.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to convert the http response into required struct: %w", err)
	}

	// If the http response data did not contain a request ID, we fetch it from the headers.
	if messageResponse.RequestID == "" {
		messageResponse.RequestID = response.Header.Get("x-request-id")
	}

	return messageResponse, nil
}

// SendMessageAsync sends a new message asynchronously.
// It uses the websocket connection for sending the message.
// The response of this request can be handled through the ResponseHandler function.
func (c *Connection) SendMessageAsync(ctx context.Context, request *OutgoingMessage) error {
	// Setting the message type.
	request.Type = typeOutgoingMessage

	// Marshalling request to byte array.
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Writing the message to the connection.
	if err := c.underlyingConn.WriteMessage(websocket.TextMessage, requestBytes); err != nil {
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
			messageType, err := unmarshalMessageType(message)
			if err != nil {
				// If the message type fails to be determined, we assume it to be an incoming message.
				conn.IncomingMessageHandler(ctx, nil, fmt.Errorf("failed to get message type: %w", err))
				continue
			}

			// Handling different message types.
			switch messageType {
			case typeIncomingMessage:
				inMessage, err := unmarshalIncomingMessage(message)
				if err != nil {
					conn.IncomingMessageHandler(ctx, nil,
						fmt.Errorf("failed to unmarshal message: %w", err))
					continue
				}
				conn.IncomingMessageHandler(ctx, inMessage, nil)
			case typeOutgoingMessageResponse:
				outMessageResp, err := unmarshalOutgoingMessageResponse(message)
				if err != nil {
					conn.OutgoingMessageResponseHandler(ctx, nil,
						fmt.Errorf("failed to unmarshal message: %w", err))
					continue
				}
				conn.OutgoingMessageResponseHandler(ctx, outMessageResp, nil)
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