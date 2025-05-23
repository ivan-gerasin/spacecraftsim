package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"spacecraftsim/internal/commands"
	"spacecraftsim/internal/heartbeat"
	"spacecraftsim/internal/parser"
)

// Client represents a TCP client
type Client struct {
	conn     net.Conn
	parser   parser.MessageParser
	commands *commands.CommandRegistry
	monitor  *heartbeat.HeartbeatMonitor
	address  string
}

// New creates a new Client instance
func New(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	client := &Client{
		conn:     conn,
		parser:   &parser.JSONParser{},
		commands: commands.NewCommandRegistry(),
		address:  address,
	}

	// Initialize heartbeat monitor
	client.monitor = heartbeat.New(
		conn,
		func() { fmt.Println("\n[Server disconnected]") },
		func() { fmt.Println("\n[Server reconnected]") },
		func() { fmt.Println("\n[Failed to reconnect after 3 attempts. Only commands are available.]") },
		func(newConn net.Conn) { client.conn = newConn },
	)
	client.monitor.Start()

	return client, nil
}

// Close closes the client connection
func (c *Client) Close() error {
	if c.monitor != nil {
		c.monitor.Stop()
	}
	return c.conn.Close()
}

// Reconnect attempts to establish a new connection
func (c *Client) Reconnect() error {
	// Close existing connection and monitor
	if c.monitor != nil {
		c.monitor.Stop()
	}
	if c.conn != nil {
		c.conn.Close()
	}

	// Try to establish new connection
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return fmt.Errorf("failed to reconnect: %w", err)
	}

	// Update client state
	c.conn = conn
	c.monitor = heartbeat.New(
		conn,
		func() { fmt.Println("\n[Server disconnected]") },
		func() { fmt.Println("\n[Server reconnected]") },
		func() { fmt.Println("\n[Failed to reconnect after 3 attempts. Only commands are available.]") },
		func(newConn net.Conn) { c.conn = newConn },
	)
	c.monitor.Start()

	return nil
}

// SendBatch sends a batch of messages to the server
func (c *Client) SendBatch(messages []parser.Message) error {
	if !c.monitor.IsConnected() {
		return fmt.Errorf("server is not connected")
	}

	// Serialize the batch to JSON
	data, err := c.parser.SerializeBatch(messages)
	if err != nil {
		return fmt.Errorf("failed to serialize batch: %w", err)
	}

	// Send the JSON data with a newline
	data = append(data, '\n')
	_, err = c.conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send data: %w", err)
	}

	return nil
}

// RunCLI starts the command-line interface
func (c *Client) RunCLI() {
	scanner := bufio.NewScanner(os.Stdin)
	var batch []parser.Message

	fmt.Println("Enter messages (empty line to send) or commands (starting with /):")
	fmt.Print("> ")

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			if len(batch) > 0 {
				if err := c.SendBatch(batch); err != nil {
					fmt.Printf("Error: %v\n", err)
				} else {
					fmt.Println("[Sent successfully]")
					batch = nil // Clear the batch
				}
			}
		} else if commands.IsCommand(line) {
			// Handle command
			cmdName := commands.ParseCommand(line)
			if cmdName == "connect" && !c.monitor.IsConnected() {
				if err := c.Reconnect(); err != nil {
					fmt.Printf("Error: %v\n", err)
				} else {
					fmt.Println("[Reconnected successfully]")
				}
			} else if err := c.commands.Execute(cmdName, c.conn); err != nil {
				fmt.Printf("Command error: %v\n", err)
			}
		} else if c.monitor.IsConnected() {
			// Only process messages if connected
			msg, err := parser.ParseLine(line)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				batch = append(batch, msg)
			}
		} else {
			fmt.Println("Server is not connected. Only commands are available.")
		}

		fmt.Print("> ")
	}
}
