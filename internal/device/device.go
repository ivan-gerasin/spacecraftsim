package device

import (
	"fmt"
	"time"
)

// Message represents a message that can be sent between devices
type Message struct {
	ID     string
	Values []interface{}
	Time   time.Time
	Source string
}

// Device represents a ship module with input/output capabilities
type Device interface {
	// ID returns the unique identifier of the device
	ID() string

	// HandleInput processes an incoming message
	HandleInput(msg Message) error

	// Tick is called periodically to update device state
	Tick() error

	// Subscribe registers the device to receive messages on specific topics
	Subscribe(bus Bus) error

	// GetTickRate returns the device's tick rate
	GetTickRate() time.Duration
}

// Bus defines the interface for inter-device communication
type Bus interface {
	// Publish sends a message to all subscribers of a topic
	Publish(topic string, msg Message) error

	// Subscribe registers a device to receive messages on a topic
	Subscribe(topic string, device Device) error

	// Unsubscribe removes a device's subscription to a topic
	Unsubscribe(topic string, device Device) error
}

// BaseDevice provides common functionality for devices
type BaseDevice struct {
	id       string
	bus      Bus
	topics   []string
	lastTick time.Time
	tickRate time.Duration
}

// NewBaseDevice creates a new base device
func NewBaseDevice(id string, tickRate time.Duration) *BaseDevice {
	return &BaseDevice{
		id:       id,
		tickRate: tickRate,
		lastTick: time.Now(),
	}
}

// ID returns the device's identifier
func (d *BaseDevice) ID() string {
	return d.id
}

// Subscribe registers the device to receive messages on specific topics
func (d *BaseDevice) Subscribe(bus Bus) error {
	d.bus = bus
	for _, topic := range d.topics {
		if err := bus.Subscribe(topic, d); err != nil {
			return fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
		}
	}
	return nil
}

// AddTopic adds a topic to the device's subscription list
func (d *BaseDevice) AddTopic(topic string) {
	d.topics = append(d.topics, topic)
}

// HandleInput provides a default implementation for BaseDevice
func (d *BaseDevice) HandleInput(msg Message) error {
	return nil
}

// Tick provides a default implementation for BaseDevice
func (d *BaseDevice) Tick() error {
	return nil
}

// GetTickRate returns the device's tick rate
func (d *BaseDevice) GetTickRate() time.Duration {
	return d.tickRate
}
