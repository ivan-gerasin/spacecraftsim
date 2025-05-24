package core

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Message represents a message to be sent to the server
type Message struct {
	ID     string        `json:"id"`
	Values []interface{} `json:"values"`
}

// SendMessage sends a message to the server
func SendMessage(w io.Writer, deviceID string, values []string) error {
	// Handle commands (messages starting with /)
	if strings.HasPrefix(deviceID, "/") {
		// Send command directly without wrapping in Message struct
		if _, err := fmt.Fprintln(w, deviceID); err != nil {
			return fmt.Errorf("failed to write command: %w", err)
		}
		return nil
	}

	// Convert string values to interface{} slice
	interfaceValues := make([]interface{}, len(values))
	for i, v := range values {
		interfaceValues[i] = v
	}

	msg := Message{
		ID:     deviceID,
		Values: interfaceValues,
	}

	// Wrap message in an array
	data, err := json.Marshal([]Message{msg})
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if _, err := fmt.Fprintln(w, string(data)); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

// ParseInput parses space-separated input into device ID and values
func ParseInput(input string) (deviceID string, values []string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", nil
	}
	return parts[0], parts[1:]
}
