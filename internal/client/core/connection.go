package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// Connection represents a TCP connection to the server
type Connection struct {
	conn           net.Conn
	messageHandler func(Message)
}

// NewConnection creates a new connection to the server
func NewConnection(address string) (*Connection, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	c := &Connection{conn: conn}
	go c.readMessages()
	return c, nil
}

// SetMessageHandler sets the handler for incoming messages
func (c *Connection) SetMessageHandler(handler func(Message)) {
	c.messageHandler = handler
}

// SendMessage sends a message to the server
func (c *Connection) SendMessage(msg Message) error {
	data, err := json.Marshal([]Message{msg})
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	if _, err := fmt.Fprintln(c.conn, string(data)); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	return nil
}

// readMessages reads messages from the connection
func (c *Connection) readMessages() {
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		var messages []Message
		if err := json.Unmarshal(scanner.Bytes(), &messages); err != nil {
			continue // Skip invalid messages
		}
		for _, msg := range messages {
			if c.messageHandler != nil {
				c.messageHandler(msg)
			}
		}
	}
}

// Write writes data to the connection
func (c *Connection) Write(p []byte) (n int, err error) {
	return c.conn.Write(p)
}

// Close closes the connection
func (c *Connection) Close() error {
	return c.conn.Close()
}

// StartHeartbeat starts sending heartbeat messages
func (c *Connection) StartHeartbeat(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			if _, err := fmt.Fprintln(c.conn, "__heartbeat__"); err != nil {
				return
			}
		}
	}()
}
