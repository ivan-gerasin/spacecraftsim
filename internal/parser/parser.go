package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Message represents a single sensor message
// Values can be float64 or string
type Message struct {
	ID     string        `json:"id"`
	Values []interface{} `json:"values"`
}

// MessageParser defines the interface for parsing messages
type MessageParser interface {
	// ParseBatch reads and parses a batch of messages from the reader
	ParseBatch(reader io.Reader) ([]Message, error)
	// SerializeBatch converts a batch of messages to JSON
	SerializeBatch(messages []Message) ([]byte, error)
}

// JSONParser implements MessageParser for JSON messages
type JSONParser struct{}

// ParseBatch implements the MessageParser interface for JSON messages
func (p *JSONParser) ParseBatch(reader io.Reader) ([]Message, error) {
	var messages []Message
	if err := json.NewDecoder(reader).Decode(&messages); err != nil {
		return nil, err
	}
	return messages, nil
}

// SerializeBatch converts messages to JSON
func (p *JSONParser) SerializeBatch(messages []Message) ([]byte, error) {
	return json.Marshal(messages)
}

// ParseLine converts a single line of input into a Message
// Values are float64 if possible, otherwise string
func ParseLine(line string) (Message, error) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return Message{}, fmt.Errorf("invalid message format: need ID and at least one value")
	}

	msg := Message{
		ID:     parts[0],
		Values: make([]interface{}, len(parts)-1),
	}

	for i, val := range parts[1:] {
		if num, err := strconv.ParseFloat(val, 64); err == nil {
			msg.Values[i] = num
		} else {
			msg.Values[i] = val
		}
	}

	return msg, nil
}
