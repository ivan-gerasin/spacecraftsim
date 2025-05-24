package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"spacecraftsim/internal/device"
	"spacecraftsim/internal/parser"
	"spacecraftsim/internal/ship"
	"strings"
	"time"
)

// Server represents a TCP server
type Server struct {
	address  string
	listener net.Listener
	parser   parser.MessageParser
	ship     *ship.Ship
}

// New creates a new Server instance
func New(address string) *Server {
	s := &Server{
		address: address,
		parser:  &parser.JSONParser{},
		ship:    ship.New(),
	}

	// Register some example devices
	s.registerDevices()

	return s
}

// registerDevices adds some example devices to the ship
func (s *Server) registerDevices() {
	// Create an example device
	logger := device.NewLogger("logger1")
	echo1 := device.NewEcho("echo1")

	// Register the devices
	if err := s.ship.RegisterDevice(logger); err != nil {
		log.Printf("Error registering logger: %v", err)
	}
	if err := s.ship.RegisterDevice(echo1); err != nil {
		log.Printf("Error registering echo1: %v", err)
	}

	// Start the ship
	s.ship.Start()
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
			s.ship.Stop()
			os.Exit(0)
		}

		if line == "__heartbeat__" {
			continue
		}

		// Process regular messages
		messages, err := s.parser.ParseBatch(strings.NewReader(line))
		if err != nil {
			log.Printf("Error parsing message: %v", err)
			resp := parser.ResponseMessage{
				Type:  "error",
				Error: fmt.Sprintf("Failed to parse message: %v", err),
			}
			if err := json.NewEncoder(conn).Encode(resp); err != nil {
				log.Printf("Error sending error response: %v", err)
			}
			continue
		}

		// Process the messages
		for _, msg := range messages {
			// Convert parser message to device message
			devMsg := device.Message{
				ID:     msg.ID,
				Values: msg.Values,
				Time:   time.Now(),
				Source: conn.RemoteAddr().String(),
			}

			// Route message to appropriate device
			if err := s.ship.HandleMessage(devMsg); err != nil {
				log.Printf("Error handling message: %v", err)
				resp := parser.ResponseMessage{
					Type:  "error",
					ID:    msg.ID,
					Error: fmt.Sprintf("Failed to handle message: %v", err),
				}
				if err := json.NewEncoder(conn).Encode(resp); err != nil {
					log.Printf("Error sending error response: %v", err)
				}
			} else {
				// Send success response
				resp := parser.ResponseMessage{
					Type:   "success",
					ID:     msg.ID,
					Values: msg.Values,
				}
				if err := json.NewEncoder(conn).Encode(resp); err != nil {
					log.Printf("Error sending success response: %v", err)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from connection: %v", err)
	}
}
