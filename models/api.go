package models

// ReceivedMessage is the schema of a received message.
type ReceivedMessage struct {
	// MessageID is the identifier of the message.
	MessageID string `json:"message_id"`
	// MessageBody is the main content of the message.
	MessageBody string `json:"message_body"`
	// SenderID is the ID of the client who sent the message.
	SenderID string `json:"sender"`
	// SendAt is the time at which Rosenbridge received the message from the sender.
	SentAt int64 `json:"sent_at"`
}
