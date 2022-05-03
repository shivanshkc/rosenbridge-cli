package lib

// Types of data sent/received over the connection to/from Rosenbridge.
const (
	typeIncomingMessage         string = "INCOMING_MESSAGE"
	typeOutgoingMessage         string = "OUTGOING_MESSAGE"
	typeOutgoingMessageResponse string = "OUTGOING_MESSAGE_RESPONSE"
)

const (
	// PersistTrue persists the message without any dependency on the online/offline state of the receiver.
	PersistTrue PersistenceCriteria = "true"
	// PersistFalse never persists the message. If the receiver is offline, the message is lost forever.
	PersistFalse PersistenceCriteria = "false"
	// PersistIfOffline persists the message only if the receiver is offline.
	PersistIfOffline PersistenceCriteria = "if_offline"
)
