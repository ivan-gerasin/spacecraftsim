package device

import (
	"log"
)

// Echo represents a device that prints received messages
type Echo struct {
	*BaseDevice
}

// NewEcho creates a new echo device
func NewEcho(id string) *Echo {
	e := &Echo{
		BaseDevice: NewBaseDevice(id, 0), // No periodic updates needed
	}
	return e
}

// HandleInput processes incoming messages and prints them
func (e *Echo) HandleInput(msg Message) error {
	log.Printf("Echo %s received: ID=%s, Values=%v, Source=%s",
		e.id, msg.ID, msg.Values, msg.Source)
	return nil
}

// Tick is not needed for echo device
func (e *Echo) Tick() error {
	return nil
}
