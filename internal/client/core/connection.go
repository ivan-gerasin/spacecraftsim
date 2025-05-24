package core

import (
	"fmt"
	"net"
	"time"
)

// Connection represents a TCP connection to the server
type Connection struct {
	conn net.Conn
}

// NewConnection creates a new connection to the server
func NewConnection(address string) (*Connection, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	return &Connection{conn: conn}, nil
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
