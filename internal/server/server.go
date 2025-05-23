package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"spacecraftsim/internal/commands"
	"spacecraftsim/internal/parser"
	"strings"
)

// Server represents a TCP server
type Server struct {
	address  string
	listener net.Listener
	parser   parser.MessageParser
	commands *commands.CommandRegistry
}

// New creates a new Server instance
func New(address string) *Server {
	return &Server{
		address:  address,
		parser:   &parser.JSONParser{},
		commands: commands.NewCommandRegistry(),
	}
}

// Start begins listening for connections
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	s.listener = listener

	log.Printf("Server listening on %s", s.address)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

// handleConnection processes a single client connection
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Printf("New connection from %s", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()

		// Check for special commands
		if line == "__kill__" {
			log.Printf("Received kill command from %s", conn.RemoteAddr())
			os.Exit(0)
		}

		if line == "__heartbeat__" {
			continue
		}

		// Process regular messages
		messages, err := s.parser.ParseBatch(strings.NewReader(line))
		if err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Process the messages
		for _, msg := range messages {
			// Format values with labels
			parts := []string{fmt.Sprintf("ID: %s", msg.ID)}
			for i, v := range msg.Values {
				parts = append(parts, fmt.Sprintf("VAL%d: %v", i+1, v))
			}
			log.Printf("%s", strings.Join(parts, " | "))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from connection: %v", err)
	}
}
