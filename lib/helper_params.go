package lib

// Types of data sent/received over the connection to/from Rosenbridge.
const (
	typeIncomingMessageReq string = "INCOMING_MESSAGE_REQ"
	typeOutgoingMessageReq string = "OUTGOING_MESSAGE_REQ"
	typeOutgoingMessageRes string = "OUTGOING_MESSAGE_RES"
	typeErrorRes           string = "ERROR_RES"
)

const (
	// codeOK is the success code for all scenarios.
	codeOK = "OK"
	// codeOffline indicates that the concerned client is offline.
	codeOffline = "OFFLINE" //nolint:unused
	// codeBridgeNotFound is sent when the required bridge does not exist.
	codeBridgeNotFound = "BRIDGE_NOT_FOUND" //nolint:unused
	// codeUnknown indicates that an unknown error occurred.
	codeUnknown = "UNKNOWN" //nolint:unused
)
