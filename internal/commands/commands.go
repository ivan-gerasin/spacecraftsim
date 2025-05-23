package commands

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// Command represents a command that can be executed
type Command interface {
	Execute(client net.Conn) error
}

// CommandRegistry holds all available commands
type CommandRegistry struct {
	commands map[string]Command
}

// NewCommandRegistry creates a new command registry
func NewCommandRegistry() *CommandRegistry {
	registry := &CommandRegistry{
		commands: make(map[string]Command),
	}

	// Register default commands
	registry.Register("exit", &ExitCommand{})
	registry.Register("kill", &KillCommand{})
	registry.Register("terminate", &TerminateCommand{})
	registry.Register("connect", &ConnectCommand{})

	return registry
}

// Register adds a new command to the registry
func (r *CommandRegistry) Register(name string, cmd Command) {
	r.commands[name] = cmd
}

// Execute runs a command by name
func (r *CommandRegistry) Execute(name string, client net.Conn) error {
	cmd, exists := r.commands[name]
	if !exists {
		return fmt.Errorf("unknown command: %s", name)
	}
	return cmd.Execute(client)
}

// ExitCommand terminates the client process
type ExitCommand struct{}

func (c *ExitCommand) Execute(client net.Conn) error {
	os.Exit(0)
	return nil
}

// KillCommand sends a signal to terminate the server
type KillCommand struct{}

func (c *KillCommand) Execute(client net.Conn) error {
	// Send special command to server
	_, err := client.Write([]byte("__kill__\n"))
	return err
}

// TerminateCommand terminates both client and server
type TerminateCommand struct{}

func (c *TerminateCommand) Execute(client net.Conn) error {
	// First kill the server
	if err := (&KillCommand{}).Execute(client); err != nil {
		return err
	}
	// Then exit the client
	return (&ExitCommand{}).Execute(client)
}

// ConnectCommand attempts to reconnect to the server
type ConnectCommand struct{}

func (c *ConnectCommand) Execute(client net.Conn) error {
	// Check if there's an active connection
	if client != nil {
		// Try to write a single byte to check connection
		one := []byte{0}
		_, err := client.Write(one)
		if err == nil {
			return fmt.Errorf("already connected to server")
		}
	}

	// Get the remote address before closing
	addr := client.RemoteAddr().String()

	// Close the current connection
	client.Close()

	// Try to establish new connection
	_, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to reconnect: %w", err)
	}

	fmt.Println("[Reconnected successfully]")
	return nil
}

// IsCommand checks if a line is a command
func IsCommand(line string) bool {
	return strings.HasPrefix(line, "/")
}

// ParseCommand extracts the command name from a line
func ParseCommand(line string) string {
	return strings.TrimPrefix(line, "/")
}
